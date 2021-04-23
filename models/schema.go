package models

// User - Blueprint of user
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Post - Blueprint of post
type Post struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}

// Connection - Blueprint of connection
type Connection struct {
	User1 string
	User2 string
}
