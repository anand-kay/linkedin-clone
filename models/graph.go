package models

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
)

// GetConnections - Fetches all the existing connections
func GetConnections(db *sql.DB) ([]*Connection, error) {
	var connections []*Connection

	cRows, err := db.Query("SELECT user_1, user_2 FROM connections;")
	if err != nil {
		return nil, err
	}
	defer cRows.Close()

	for cRows.Next() {
		connection := new(Connection)

		if err := cRows.Scan(&connection.User1, &connection.User2); err != nil {
			return nil, err
		}

		connections = append(connections, connection)
	}

	return connections, nil
}

// AddToGraph - Adds a connection to the graph
func (connection *Connection) AddToGraph(ctx context.Context, rdb *redis.Client) error {
	_, err := rdb.SAdd(ctx, "users", connection.User1, connection.User2).Result()
	if err != nil {
		return err
	}

	_, err = rdb.SAdd(ctx, "user:"+connection.User1, connection.User2).Result()
	if err != nil {
		return err
	}

	_, err = rdb.SAdd(ctx, "user:"+connection.User2, connection.User1).Result()
	if err != nil {
		return err
	}

	return nil
}
