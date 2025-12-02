package entity

import (
	"time"
)

type User struct {
	ID         int
	CreatedAt  time.Time
	UpdatedAt  time.Time
 	Username 	string
	Password	string
	Email 		string
	Role		string
}