package utils

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/anand-kay/linkedin-clone/server"

	"github.com/dgrijalva/jwt-go"
)

type ContextKey string

// Request context keys for authorization
var ContextUserIDKey ContextKey = "UserID"
var ContextEmailKey ContextKey = "Email"
var ContextFirstNameKey ContextKey = "FirstName"
var ContextLastNameKey ContextKey = "LastName"

// ReadEnvFile - Reads the .env file and returns the config as a string
func ReadEnvFile(envPath string) (string, error) {
	envBytes, err := ioutil.ReadFile(envPath)
	if err != nil {
		return "", err
	}

	return string(envBytes), nil
}

// GetServerFromReqContext - Returns the Server instance attached in the request context
func GetServerFromReqContext(req *http.Request) *server.Server {
	return req.Context().Value("server").(*server.Server)
}

// GenerateJWT - Generates a new JWT
func GenerateJWT(userID int64, email string, firstName string, lastName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserID":    userID,
		"Email":     email,
		"FirstName": firstName,
		"LastName":  lastName,
		"Access":    "all",
		"ExpiresAt": time.Now().Unix(),
	})

	signedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateID - Validates user id and post id
func ValidateID(id string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil || idInt < 1 {
		return errors.New("Invalid ID")
	}

	return nil
}
