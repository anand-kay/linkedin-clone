package models

import "database/sql"

func (user *User) GetUserInfo(db *sql.DB) error {
	uRow := db.QueryRow("SELECT email, first_name, last_name FROM users WHERE id=$1;", user.ID)

	return uRow.Scan(&user.Email, &user.FirstName, &user.LastName)
}
