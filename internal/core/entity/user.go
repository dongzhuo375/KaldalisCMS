package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Password  string // This will store the hashed password
	Email     string
	Role      string
}

// SetPassword hashes the given plaintext password and stores the hash in the User struct.
func (u *User) SetPassword(plainTextPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword compares the given plaintext password with the stored hash.
// It returns true if the password is correct, and false otherwise.
func (u *User) CheckPassword(plainTextPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainTextPassword))
	return err == nil
}