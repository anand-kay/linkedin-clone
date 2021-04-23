package handlers

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/anand-kay/linkedin-clone/libs"
	"github.com/anand-kay/linkedin-clone/models"
	"github.com/anand-kay/linkedin-clone/utils"
)

// Signup - Signs up new user with email, password, first name, last name
func Signup(w http.ResponseWriter, req *http.Request) {
	var user models.User

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	json.Unmarshal(reqBody, &user)

	err = libs.ValidateSignupForm(user.Email, user.Password, user.FirstName, user.LastName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	hashedPwd, err := libs.HashPassword(user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	user.Password = hashedPwd

	token, err := user.CreateUser(req.Context(), utils.GetServerFromReqContext(req).DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	http.SetCookie(w, &http.Cookie{
		Name:    "x-auth",
		Value:   token,
		Expires: time.Now().Add(2 * time.Hour),
	})
	w.Write([]byte(token))
}

// Login - Logs in user with email, password
func Login(w http.ResponseWriter, req *http.Request) {
	var user models.User

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	json.Unmarshal(reqBody, &user)

	err = libs.ValidateLoginForm(user.Email, user.Password, user.FirstName, user.LastName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	userID, storedPwd, err := user.RetreiveIdAndHashedPwd(utils.GetServerFromReqContext(req).DB)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("User not found"))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	isPwdMatch := libs.CheckHash(user.Password, storedPwd)
	if !isPwdMatch {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Incorrect password"))
		return
	}

	token, err := utils.GenerateJWT(userID, user.Email, user.FirstName, user.LastName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	http.SetCookie(w, &http.Cookie{
		Name:    "x-auth",
		Value:   token,
		Expires: time.Now().Add(2 * time.Hour),
	})
	w.Write([]byte(token))
}

// Logout - Logs out a user
func Logout(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}
