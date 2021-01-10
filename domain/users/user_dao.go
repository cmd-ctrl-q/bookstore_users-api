package users

import (
	"fmt"

	"github.com/cmd-ctrl-q/bookstore_users-api/utils/errors"
)

var (
	usersDB = make(map[int64]*User)
)

// Get gets the user's primary key / id
func (user *User) Get() *errors.RestErr {
	result := usersDB[user.ID]
	if result == nil {
		return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.ID))
	}

	// fill in user fields with data from db
	user.ID = result.ID
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.Email = result.Email
	user.DateCreated = result.DateCreated

	return nil
}

// Save attempts to save a user in the database
func (user *User) Save() *errors.RestErr {
	current := usersDB[user.ID]
	if current != nil {
		// user already exists in db
		if current.Email == user.Email {
			return errors.NewBadRequestError(fmt.Sprintf("email %s already registered", user.Email))
		}
		return errors.NewBadRequestError(fmt.Sprintf("user %d already exists", user.ID))
	}
	// save new user
	usersDB[user.ID] = user
	return nil
}
