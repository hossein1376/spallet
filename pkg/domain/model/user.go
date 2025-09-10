package model

type User struct {
	ID       UserID `json:"id"`
	Username string `json:"username"`
}

type UserID int64
