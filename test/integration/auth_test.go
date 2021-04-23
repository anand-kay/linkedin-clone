package integration

import (
	"database/sql/driver"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/anand-kay/linkedin-clone/libs"
)

type AnyPassword struct{}

func (ap AnyPassword) Match(v driver.Value) bool {
	_, ok := v.(string)

	return ok
}

func TestSignup(t *testing.T) {
	sqlMock := setupServer()

	resCode, _ := runTestSignup(sqlMock, "1", "test1@test.com", "fgh54cv34bvb7", "Don", "Smith")
	if resCode != http.StatusOK {
		t.Error("TestSignup - Failed")
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	resCode, _ = runTestSignup(sqlMock, "2", "test2@test.com", "$@df", "Jake", "Sparrow")
	if resCode != http.StatusBadRequest {
		t.Error("TestSignup - Failed")
	}

	teardownServer()
}

func TestLogin(t *testing.T) {
	sqlMock := setupServer()

	resCode, _ := runTestLogin(sqlMock, "testlogin@test.com", "ewdf64vc6vc")
	if resCode != http.StatusOK {
		t.Error("TestLogin - Failed")
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	resCode, _ = runTestLogin(sqlMock, "telogtest.com", "ewdf64vc6vc")
	if resCode != http.StatusBadRequest {
		t.Error("TestLogin - Failed")
	}

	teardownServer()
}

func TestLogout(t *testing.T) {
	setupServer()

	req, err := http.NewRequest("POST", "http://localhost:3000/logout", nil)
	if err != nil {
		log.Fatalln("Error while trying to create new logout request", err)
	}

	req.Header.Set("Auth-Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2Nlc3MiOiJhbGwiLCJFbWFpbCI6InRlc3QxQHRlc3QuY29tIiwiRXhwaXJlc0F0IjoxNjE3NDUxODEzLCJGaXJzdE5hbWUiOiJEb24iLCJMYXN0TmFtZSI6IlNtaXRoIiwiVXNlcklEIjoxfQ.0fFudOT4sVcai6TNnseSj0_zwnpCy1WYvsNjnf4ercU")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("Error while trying to logout", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Error("TestLogout - Failed")
	}

	teardownServer()
}

func runTestSignup(sqlMock sqlmock.Sqlmock, id string, email string, password string, firstName string, lastName string) (int, string) {
	columns := []string{"id"}
	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users(email, password, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING id;`)).WithArgs(email, AnyPassword{}, firstName, lastName).WillReturnRows(sqlmock.NewRows(columns).FromCSVString(id))
	sqlMock.ExpectCommit()

	bodyReader := strings.NewReader(`{"email": "` + email + `", "password": "` + password + `", "first_name": "` + firstName + `", "last_name": "` + lastName + `"}`)
	res, err := http.Post("http://localhost:3000/signup", "application/json", bodyReader)
	if err != nil {
		log.Fatalln(err)
	}
	body, _ := ioutil.ReadAll(res.Body)

	return res.StatusCode, string(body)
}

func runTestLogin(sqlMock sqlmock.Sqlmock, email string, password string) (int, string) {
	hashedPwd, err := libs.HashPassword(password)
	if err != nil {
		log.Fatalln("Error while hashing password")
	}

	columns := []string{"id", "password"}
	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password FROM users WHERE email=$1;`)).WithArgs(email).WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1," + hashedPwd))

	bodyReader := strings.NewReader(`{"email": "` + email + `", "password": "` + password + `"}`)
	res, err := http.Post("http://localhost:3000/login", "application/json", bodyReader)
	if err != nil {
		log.Fatalln(err)
	}
	body, _ := ioutil.ReadAll(res.Body)

	return res.StatusCode, string(body)
}
