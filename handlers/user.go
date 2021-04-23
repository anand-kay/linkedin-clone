package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/anand-kay/linkedin-clone/libs"
	"github.com/anand-kay/linkedin-clone/models"
	"github.com/anand-kay/linkedin-clone/utils"
)

type response struct {
	User  models.User
	Level uint8 `json:"level"`
}

var otherUser models.User

// UserInfo - Returns the user info of a particular id
func UserInfo(w http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value(utils.ContextUserIDKey).(string)
	otherUserID := req.URL.Query().Get("id")

	err := utils.ValidateID(otherUserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid user id"))

		return
	}

	otherUser.ID = otherUserID

	err = otherUser.GetUserInfo(utils.GetServerFromReqContext(req).DB)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))

		return
	}

	var res response

	res.User = otherUser

	if otherUserID != userID {
		res.Level = libs.CheckLevel(req.Context(), utils.GetServerFromReqContext(req).RedisDB, userID, otherUserID)
	}

	resJSON, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating JSON response"))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}
