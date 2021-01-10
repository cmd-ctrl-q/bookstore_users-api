package services

import (
	"github.com/cmd-ctrl-q/bookstore_users-api/domain/users"
	"github.com/cmd-ctrl-q/bookstore_users-api/utils/errors"
)

// CreateUser creates the user data that is received from the CreateUser controller
func CreateUser(user users.User) (*users.User, *errors.RestErr) {

	// prepare user data
	if err := user.Validate(); err != nil {
		return nil, err
	}
	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUser gets and returns a user by their user id
func GetUser(userID int64) (*users.User, *errors.RestErr) {

	// create new instance of user and give it userID
	result := &users.User{ID: userID}

	// check if user id exists in db
	if err := result.Get(); err != nil {
		return nil, err
	}

	// user exists
	return result, nil
}
