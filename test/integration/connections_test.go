package integration

import (
	"database/sql"
	"log"
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestFetchConnections(t *testing.T) {
	sqlMock := setupServer()

	userID := "1"

	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT user_2 FROM connections WHERE user_1=$1;`)).WithArgs(userID).WillReturnRows(sqlmock.NewRows([]string{"user_2"}).FromCSVString("2").FromCSVString("3").FromCSVString("4"))
	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT user_1 FROM connections WHERE user_2=$1;`)).WithArgs(userID).WillReturnRows(sqlmock.NewRows([]string{"user_1"}).FromCSVString("5").FromCSVString("6").FromCSVString("7"))
	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT sender_id, receiver_id FROM pending_requests WHERE sender_id=$1 OR receiver_id=$1;`)).WithArgs(userID).WillReturnRows(sqlmock.NewRows([]string{"sender_id", "receiver_id"}).FromCSVString("1,8").FromCSVString("9,1").FromCSVString("10,1").FromCSVString("1,11").FromCSVString("12,1"))

	req, err := http.NewRequest("GET", "http://localhost:3000/connections", nil)
	if err != nil {
		log.Fatalln("Error while trying to create new FetchConnections request", err)
	}

	req.Header.Set("Auth-Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2Nlc3MiOiJhbGwiLCJFbWFpbCI6InRlc3QxQHRlc3QuY29tIiwiRXhwaXJlc0F0IjoxNjE3NDUxODEzLCJGaXJzdE5hbWUiOiJEb24iLCJMYXN0TmFtZSI6IlNtaXRoIiwiVXNlcklEIjoxfQ.0fFudOT4sVcai6TNnseSj0_zwnpCy1WYvsNjnf4ercU")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("Error while trying to send a FetchConnections request", err)
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Error("FetchConnections - Failed")
	}

	teardownServer()
}

func TestSendReq(t *testing.T) {
	sqlMock := setupServer()

	userID := "1"
	otherUserID := "2"

	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE id=$1;`)).WithArgs(otherUserID).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString(otherUserID))
	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM pending_requests WHERE (sender_id=$1 AND receiver_id=$2) OR (sender_id=$2 AND receiver_id=$1);`)).WithArgs(userID, otherUserID).WillReturnError(sql.ErrNoRows)
	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM connections WHERE (user_1=$1 AND user_2=$2) OR (user_1=$2 AND user_2=$1);`)).WithArgs(userID, otherUserID).WillReturnError(sql.ErrNoRows)
	sqlMock.ExpectExec(regexp.QuoteMeta(`INSERT INTO pending_requests (sender_id, receiver_id) VALUES ($1, $2);`)).WithArgs(userID, otherUserID).WillReturnResult(sqlmock.NewResult(1, 1))

	req, err := http.NewRequest("POST", "http://localhost:3000/connections/sendreq?id="+otherUserID, nil)
	if err != nil {
		log.Fatalln("Error while trying to create new SendReq request", err)
	}

	req.Header.Set("Auth-Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2Nlc3MiOiJhbGwiLCJFbWFpbCI6InRlc3QxQHRlc3QuY29tIiwiRXhwaXJlc0F0IjoxNjE3NDUxODEzLCJGaXJzdE5hbWUiOiJEb24iLCJMYXN0TmFtZSI6IlNtaXRoIiwiVXNlcklEIjoxfQ.0fFudOT4sVcai6TNnseSj0_zwnpCy1WYvsNjnf4ercU")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("Error while trying to send a SendReq request", err)
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Error("TestSendReq - Failed")
	}

	teardownServer()
}

func TestAcceptReq(t *testing.T) {
	sqlMock := setupServer()

	userID := "1"
	otherUserID := "2"

	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE id=$1;`)).WithArgs(otherUserID).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString(otherUserID))
	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM pending_requests WHERE sender_id=$1 AND receiver_id=$2;`)).WithArgs(otherUserID, userID).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
	sqlMock.ExpectBegin()
	sqlMock.ExpectExec(regexp.QuoteMeta(`INSERT INTO connections(user_1, user_2) VALUES ($1, $2);`)).WithArgs(otherUserID, userID).WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta(`DELETE FROM pending_requests WHERE sender_id=$1 AND receiver_id=$2;`)).WithArgs(otherUserID, userID).WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectCommit()

	req, err := http.NewRequest("POST", "http://localhost:3000/connections/acceptreq?id="+otherUserID, nil)
	if err != nil {
		log.Fatalln("Error while trying to create new AcceptReq request", err)
	}

	req.Header.Set("Auth-Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2Nlc3MiOiJhbGwiLCJFbWFpbCI6InRlc3QxQHRlc3QuY29tIiwiRXhwaXJlc0F0IjoxNjE3NDUxODEzLCJGaXJzdE5hbWUiOiJEb24iLCJMYXN0TmFtZSI6IlNtaXRoIiwiVXNlcklEIjoxfQ.0fFudOT4sVcai6TNnseSj0_zwnpCy1WYvsNjnf4ercU")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("Error while trying to send a AcceptReq request", err)
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Error("TestAcceptReq - Failed")
	}

	teardownServer()
}

func TestRevokeReq(t *testing.T) {
	sqlMock := setupServer()

	userID := "1"
	otherUserID := "2"

	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE id=$1;`)).WithArgs(otherUserID).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString(otherUserID))
	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM pending_requests WHERE sender_id=$1 AND receiver_id=$2;`)).WithArgs(otherUserID, userID).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
	sqlMock.ExpectExec(regexp.QuoteMeta(`DELETE FROM pending_requests WHERE sender_id=$1 AND receiver_id=$2;`)).WithArgs(otherUserID, userID).WillReturnResult(sqlmock.NewResult(1, 1))

	req, err := http.NewRequest("DELETE", "http://localhost:3000/connections/revokereq?id="+otherUserID, nil)
	if err != nil {
		log.Fatalln("Error while trying to create new RevokeReq request", err)
	}

	req.Header.Set("Auth-Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2Nlc3MiOiJhbGwiLCJFbWFpbCI6InRlc3QxQHRlc3QuY29tIiwiRXhwaXJlc0F0IjoxNjE3NDUxODEzLCJGaXJzdE5hbWUiOiJEb24iLCJMYXN0TmFtZSI6IlNtaXRoIiwiVXNlcklEIjoxfQ.0fFudOT4sVcai6TNnseSj0_zwnpCy1WYvsNjnf4ercU")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("Error while trying to send a RevokeReq request", err)
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Error("TestRevokeReq - Failed")
	}

	teardownServer()
}
