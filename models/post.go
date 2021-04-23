package models

import (
	"database/sql"
)

// CreatePost - Creates a new post entry in the DB
func (post *Post) CreatePost(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO posts(user_id, text) VALUES ($1, $2);", post.UserID, post.Text)

	return err
}

func (post *Post) GetPostByID(db *sql.DB) error {
	pRow := db.QueryRow("SELECT user_id, text FROM posts WHERE id=$1;", post.ID)

	return pRow.Scan(&post.UserID, &post.Text)

	// switch err := pRow.Scan(&post.UserID, &post.Text); err {
	// case sql.ErrNoRows:
	// 	return err
	// case nil:
	// 	return nil
	// default:
	// 	return err
	// }
}

func GetPostsByUserID(db *sql.DB, userID string, page int, limit int) ([]Post, error) {
	pRows, err := db.Query("SELECT id, text FROM posts WHERE user_id=$1 LIMIT $2 OFFSET $3;", userID, limit, (page * limit))
	if err != nil {
		return nil, err
	}
	defer pRows.Close()

	var posts []Post

	for pRows.Next() {
		var post Post

		if err := pRows.Scan(&post.ID, &post.Text); err != nil {
			return nil, err
		}

		post.UserID = userID

		posts = append(posts, post)
	}

	return posts, nil
}
