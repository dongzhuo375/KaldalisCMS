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


//dzCake,这里写下基本的业务逻辑