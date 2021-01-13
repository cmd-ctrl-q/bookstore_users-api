package users

import (
	"strings"
	"unicode"

	"github.com/cmd-ctrl-q/bookstore_users-api/utils/errors"
)

const (
	// StatusActive ...
	StatusActive = "active"
)

// User is a struct for a user
type User struct {
	ID          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:"-"` // hide field
}

// Users is a slice of User
type Users []User

// Validate cleans up user data before allowing other services to process it.
func (user *User) Validate() *errors.RestErr {
	user.FirstName = strings.TrimSpace(strings.ToLower(user.FirstName))
	user.LastName = strings.TrimSpace(strings.ToLower(user.LastName))

	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return errors.NewBadRequestError("invalid email address")
	}

	user.Password = strings.TrimSpace(user.Password)
	// valid := isValid(user.Password)
	// if !valid {
	// 	return errors.NewBadRequestError("invalid password")
	// }
	return nil
}

// TODO: create method that validates a password,
// including security measures to define on password.
// 33 > password > 7 chars
// password has atleast 1 uppercase
// password has atleast 1 lowercase
// password has atleast 1 number
// password has atleast 1 special character

func isValid(password string) bool {

	var (
		upper   = false
		lower   = false
		number  = false
		special = false
	)

	if len(password) < 8 || len(password) > 32 {
		return false
	}

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsLower(c):
			lower = true
		case unicode.IsNumber(c):
			number = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case c == ' ':
			return false
		}
	}
	return upper && lower && number && special
}
