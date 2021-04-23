package models

import (
	"context"
	"database/sql"

	"github.com/anand-kay/linkedin-clone/utils"
)

// CreateUser - Creates a new user entry in the DB
func (user *User) CreateUser(ctx context.Context, db *sql.DB) (string, error) {
	txOptions := &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	}

	tx, err := db.BeginTx(ctx, txOptions)
	if err != nil {
		return "", err
	}

	userID, err := user.insertUserToDB(tx)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	token, err := utils.GenerateJWT(userID, user.Email, user.FirstName, user.LastName)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return token, nil
}

// RetreiveHashedPwd - Fetches the id and hashed password of the user
func (user *User) RetreiveIdAndHashedPwd(db *sql.DB) (int64, string, error) {
	var id int64
	var password string

	userRow := db.QueryRow("SELECT id, password FROM users WHERE email=$1;", user.Email)

	switch err := userRow.Scan(&id, &password); err {
	case sql.ErrNoRows:
		return -1, "", sql.ErrNoRows
	case nil:
		return id, password, nil
	default:
		return -1, "", err
	}
}

func (user *User) insertUserToDB(tx *sql.Tx) (int64, error) {
	var userID int64

	// Using QueryRow() instead of Exec() to retreive the userID
	err := tx.QueryRow(
		"INSERT INTO users(email, password, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING id;",
		user.Email, user.Password, user.FirstName, user.LastName).Scan(&userID)
	if err != nil {
		return -1, err
	}

	return userID, nil
}
