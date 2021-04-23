package integration

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserInfo(t *testing.T) {
	sqlMock := setupServer()

	resCode, _ := runTestUserInfo(sqlMock, "2", "usrinf@test.com", "Dwight", "Blake")
	if resCode != http.StatusOK {
		t.Error("TestUserInfo - Failed")
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	resCode, _ = runTestUserInfo(sqlMock, "r", "usrinf@test.com", "Dwight", "Blake")
	if resCode != http.StatusBadRequest {
		t.Error("TestUserInfo - Failed")
	}

	teardownServer()
}

func runTestUserInfo(sqlMock sqlmock.Sqlmock, id string, email string, firstName string, lastName string) (int, string) {
	columns := []string{"email", "first_name", "last_name"}
	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT email, first_name, last_name FROM users WHERE id=$1;`)).WithArgs(id).WillReturnRows(sqlmock.NewRows(columns).AddRow(email, firstName, lastName))

	req, err := http.NewRequest("GET", "http://localhost:3000/user/info?id="+id, nil)
	if err != nil {
		log.Fatalln("Error while trying to create new user info request", err)
	}

	req.Header.Set("Auth-Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2Nlc3MiOiJhbGwiLCJFbWFpbCI6InRlc3QxQHRlc3QuY29tIiwiRXhwaXJlc0F0IjoxNjE3NDUxODEzLCJGaXJzdE5hbWUiOiJEb24iLCJMYXN0TmFtZSI6IlNtaXRoIiwiVXNlcklEIjoxfQ.0fFudOT4sVcai6TNnseSj0_zwnpCy1WYvsNjnf4ercU")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	resBody, _ := ioutil.ReadAll(res.Body)

	return res.StatusCode, string(resBody)
}
