package libs

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis"
	"github.com/anand-kay/linkedin-clone/libs"
	"github.com/go-redis/redis/v8"
)

func TestGraph(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub redis db connection", err)
	}
	defer s.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	testPopulateGraph(t, rdb)

	testCheckLevel(t, rdb)
}

func testPopulateGraph(t *testing.T, rdb *redis.Client) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{"user_1", "user_2"}

	mock.ExpectQuery("SELECT user_1, user_2 FROM connections;").WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,4").FromCSVString("6,4").FromCSVString("9,4").FromCSVString("7,6").FromCSVString("8,7"))

	err = libs.PopulateGraph(context.Background(), db, rdb)
	if err != nil {
		t.Errorf("an error '%s' was not expected when populating graph in redis db", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if ok, err := rdb.SIsMember(context.Background(), "users", "1").Result(); !ok || err != nil {
		if err != nil {
			t.Errorf("an error '%s' was not expected when checking for users in redis db", err)
		} else {
			t.Error("User should've existed")
		}
	}

	if ok, err := rdb.SIsMember(context.Background(), "users", "6").Result(); !ok || err != nil {
		if err != nil {
			t.Errorf("an error '%s' was not expected when checking for users in redis db", err)
		} else {
			t.Error("User should've existed")
		}
	}

	if ok, err := rdb.SIsMember(context.Background(), "users", "2").Result(); ok || err != nil {
		if err != nil {
			t.Errorf("an error '%s' was not expected when checking for users in redis db", err)
		} else {
			t.Error("User shouldn't have existed")
		}
	}

	if ok, err := rdb.SIsMember(context.Background(), "user:7", "6").Result(); !ok || err != nil {
		if err != nil {
			t.Errorf("an error '%s' was not expected when checking for connections in redis db", err)
		} else {
			t.Error("Connection should've existed")
		}
	}

	if ok, err := rdb.SIsMember(context.Background(), "user:8", "7").Result(); !ok || err != nil {
		if err != nil {
			t.Errorf("an error '%s' was not expected when checking for connections in redis db", err)
		} else {
			t.Error("Connection should've existed")
		}
	}

	if ok, err := rdb.SIsMember(context.Background(), "user:1", "8").Result(); ok || err != nil {
		if err != nil {
			t.Errorf("an error '%s' was not expected when checking for connections in redis db", err)
		} else {
			t.Error("Connection shouldn't have existed")
		}
	}

	if ok, err := rdb.SIsMember(context.Background(), "user:7", "9").Result(); ok || err != nil {
		if err != nil {
			t.Errorf("an error '%s' was not expected when checking for connections in redis db", err)
		} else {
			t.Error("Connection shouldn't have existed")
		}
	}

	if ok, err := rdb.SIsMember(context.Background(), "user:7", "10").Result(); ok || err != nil {
		if err != nil {
			t.Errorf("an error '%s' was not expected when checking for connections in redis db", err)
		} else {
			t.Error("Connection shouldn't have existed")
		}
	}
}

func testCheckLevel(t *testing.T, rdb *redis.Client) {
	if libs.CheckLevel(context.Background(), rdb, "9", "8") != 99 {
		t.Error("Incorrect connection level")
	}

	if libs.CheckLevel(context.Background(), rdb, "1", "4") != 1 {
		t.Error("Incorrect connection level")
	}

	if libs.CheckLevel(context.Background(), rdb, "4", "7") != 2 {
		t.Error("Incorrect connection level")
	}

	if libs.CheckLevel(context.Background(), rdb, "7", "9") != 3 {
		t.Error("Incorrect connection level")
	}

	if libs.CheckLevel(context.Background(), rdb, "8", "8") != 0 {
		t.Error("Incorrect connection level")
	}
}
