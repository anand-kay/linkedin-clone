package integration

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreatePost(t *testing.T) {
	sqlMock := setupServer()

	postText := "from a test"

	sqlMock.ExpectExec(regexp.QuoteMeta(`INSERT INTO posts(user_id, text) VALUES ($1, $2);`)).WithArgs("1", postText).WillReturnResult(sqlmock.NewResult(1, 1))

	bodyReader := strings.NewReader(`{"text": "` + postText + `"}`)
	req, err := http.NewRequest("POST", "http://localhost:3000/post/create", bodyReader)
	if err != nil {
		log.Fatalln("Error while trying to create new create post request", err)
	}

	req.Header.Set("Auth-Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2Nlc3MiOiJhbGwiLCJFbWFpbCI6InRlc3QxQHRlc3QuY29tIiwiRXhwaXJlc0F0IjoxNjE3NDUxODEzLCJGaXJzdE5hbWUiOiJEb24iLCJMYXN0TmFtZSI6IlNtaXRoIiwiVXNlcklEIjoxfQ.0fFudOT4sVcai6TNnseSj0_zwnpCy1WYvsNjnf4ercU")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("Error while trying to send a create post request", err)
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Error("TestCreatePost - Failed")
	}

	teardownServer()
}

func TestFetchPostByID(t *testing.T) {
	sqlMock := setupServer()

	postID := "1"

	columns := []string{"user_id", "text"}
	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT user_id, text FROM posts WHERE id=$1;`)).WithArgs(postID).WillReturnRows(sqlmock.NewRows(columns).FromCSVString("4, This is a post text"))

	req, err := http.NewRequest("GET", "http://localhost:3000/post/"+postID, nil)
	if err != nil {
		log.Fatalln("Error while trying to create new FetchPostByID request", err)
	}

	req.Header.Set("Auth-Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2Nlc3MiOiJhbGwiLCJFbWFpbCI6InRlc3QxQHRlc3QuY29tIiwiRXhwaXJlc0F0IjoxNjE3NDUxODEzLCJGaXJzdE5hbWUiOiJEb24iLCJMYXN0TmFtZSI6IlNtaXRoIiwiVXNlcklEIjoxfQ.0fFudOT4sVcai6TNnseSj0_zwnpCy1WYvsNjnf4ercU")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("Error while trying to send a FetchPostByID request", err)
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Error("FetchPostByID - Failed")
	}

	teardownServer()
}

func TestFetchAllPosts(t *testing.T) {
	sqlMock := setupServer()

	otherUserID := "3"
	page := 1
	limit := 10

	columns := []string{"id", "text"}
	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT id, text FROM posts WHERE user_id=$1 LIMIT $2 OFFSET $3;`)).WithArgs(otherUserID, limit, (page * limit)).WillReturnRows(sqlmock.NewRows(columns).FromCSVString("11, Eleventh post").FromCSVString("12, Twelfth post").FromCSVString("13, Thirteenth post").FromCSVString("14, Fourteenth post"))

	req, err := http.NewRequest("GET", "http://localhost:3000/post/posts?id="+otherUserID+"&page="+strconv.Itoa(page)+"&limit="+strconv.Itoa(limit), nil)
	if err != nil {
		log.Fatalln("Error while trying to create new FetchAllPosts request", err)
	}

	req.Header.Set("Auth-Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2Nlc3MiOiJhbGwiLCJFbWFpbCI6InRlc3QxQHRlc3QuY29tIiwiRXhwaXJlc0F0IjoxNjE3NDUxODEzLCJGaXJzdE5hbWUiOiJEb24iLCJMYXN0TmFtZSI6IlNtaXRoIiwiVXNlcklEIjoxfQ.0fFudOT4sVcai6TNnseSj0_zwnpCy1WYvsNjnf4ercU")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("Error while trying to send a FetchAllPosts request", err)
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Error("FetchAllPosts - Failed")
	}

	teardownServer()
}
