package models

import (
	"context"
	"database/sql"
)

func GetConnectionIDs(db *sql.DB, userID string) ([]string, error) {
	var connectionIDs []string

	cRows, err := db.Query("SELECT user_2 FROM connections WHERE user_1=$1;", userID)
	if err != nil {
		return nil, err
	}
	defer cRows.Close()

	rightIDs, err := extractConnectionIDs(cRows)
	if err != nil {
		return nil, err
	}
	connectionIDs = append(connectionIDs, rightIDs...)

	cRows, err = db.Query("SELECT user_1 FROM connections WHERE user_2=$1;", userID)
	if err != nil {
		return nil, err
	}
	defer cRows.Close()

	leftIDs, err := extractConnectionIDs(cRows)
	if err != nil {
		return nil, err
	}
	connectionIDs = append(connectionIDs, leftIDs...)

	return connectionIDs, nil
}

func GetPendingRequests(db *sql.DB, userID string) ([]Connection, error) {
	var pendingRequests []Connection

	prRows, err := db.Query("SELECT sender_id, receiver_id FROM pending_requests WHERE sender_id=$1 OR receiver_id=$1;", userID)
	if err != nil {
		return nil, err
	}
	defer prRows.Close()

	for prRows.Next() {
		var pendingRequest Connection

		if err = prRows.Scan(&pendingRequest.User1, &pendingRequest.User2); err != nil {
			return nil, err
		}

		pendingRequests = append(pendingRequests, pendingRequest)
	}

	return pendingRequests, nil
}

func CheckUserExists(db *sql.DB, userID string) error {
	var id string

	uRow := db.QueryRow("SELECT id FROM users WHERE id=$1;", userID)

	return uRow.Scan(&id)
}

func (connection *Connection) CheckPendingReqs(db *sql.DB) error {
	var id string

	prRow := db.QueryRow("SELECT id FROM pending_requests WHERE (sender_id=$1 AND receiver_id=$2) OR (sender_id=$2 AND receiver_id=$1);", connection.User1, connection.User2)

	return prRow.Scan(&id)
}

func (connection *Connection) CheckConnections(db *sql.DB) error {
	var id string

	cRow := db.QueryRow("SELECT id FROM connections WHERE (user_1=$1 AND user_2=$2) OR (user_1=$2 AND user_2=$1);", connection.User1, connection.User2)

	return cRow.Scan(&id)
}

func (connection *Connection) SendReq(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO pending_requests (sender_id, receiver_id) VALUES ($1, $2);", connection.User1, connection.User2)

	return err
}

func (connection *Connection) CheckReqExists(db *sql.DB) error {
	var id string

	prRow := db.QueryRow("SELECT id FROM pending_requests WHERE sender_id=$1 AND receiver_id=$2;", connection.User2, connection.User1)

	return prRow.Scan(&id)
}

func (connection *Connection) AcceptReq(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "INSERT INTO connections(user_1, user_2) VALUES ($1, $2);", connection.User2, connection.User1)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM pending_requests WHERE sender_id=$1 AND receiver_id=$2;", connection.User2, connection.User1)
	if err != nil {
		return err
	}

	return nil
}

func (connection *Connection) RevokeReq(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM pending_requests WHERE sender_id=$1 AND receiver_id=$2;", connection.User2, connection.User1)
	if err != nil {
		return err
	}

	return nil
}

func extractConnectionIDs(cRows *sql.Rows) ([]string, error) {
	var connectionIDs []string

	for cRows.Next() {
		var connectionID string

		if err := cRows.Scan(&connectionID); err != nil {
			return nil, err
		}

		connectionIDs = append(connectionIDs, connectionID)
	}

	return connectionIDs, nil
}
