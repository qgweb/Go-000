package model

type User struct {
	ID     int64  `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Passwd string `json:"passwd" db:"passwd"`
}
