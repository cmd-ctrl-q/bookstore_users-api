package services

import (
	"github.com/cmd-ctrl-q/bookstore_users-api/domain/users"
	"github.com/cmd-ctrl-q/bookstore_users-api/utils/crypto_utils"
	"github.com/cmd-ctrl-q/bookstore_users-api/utils/date_utils"
	"github.com/cmd-ctrl-q/bookstore_users-api/utils/errors"
)

var (
	// UsersService should be used in the controller in order to use the userServiceInterface methods
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct{}

// UserServiceInterface is an interface for user CRUD methods
type usersServiceInterface interface {
	CreateUser(users.User) (*users.User, *errors.RestErr)
	GetUser(int64) (*users.User, *errors.RestErr)
	UpdateUser(bool, users.User) (*users.User, *errors.RestErr)
	DeleteUser(int64) *errors.RestErr
	SearchUser(string) (users.Users, *errors.RestErr)
}

// CreateUser creates the user data that is received from the CreateUser controller
func (s *usersService) CreateUser(user users.User) (*users.User, *errors.RestErr) {

	// validate user data
	if err := user.Validate(); err != nil {
		return nil, err
	}

	user.Status = users.StatusActive
	user.DateCreated = date_utils.GetNowDBFormat()
	user.Password = crypto_utils.GetMD5(user.Password)

	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUser gets and returns a user by their user id
func (s *usersService) GetUser(userID int64) (*users.User, *errors.RestErr) {

	// create new instance of user and give it userID
	result := &users.User{ID: userID}

	// check if user id exists in db
	if err := result.Get(); err != nil {
		return nil, err
	}

	// user exists
	return result, nil
}

// UpdateUser updates current user and returns an updated user
func (s *usersService) UpdateUser(isPartial bool, user users.User) (*users.User, *errors.RestErr) {

	// check and return current user in db
	current, err := UsersService.GetUser(user.ID)
	if err != nil {
		return nil, err
	}

	// validate the updated user fields
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// method = patch, else method = put
	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}
		if user.LastName != "" {
			current.LastName = user.LastName
		}
		if user.Email != "" {
			current.Email = user.Email
		}
	} else {
		// update current user with the new user data
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
	}

	// update the db with the current user data
	if err := current.Update(); err != nil {
		return nil, err
	}

	return current, nil
}

// DeleteUser attempts to delete a user from the database
func (s *usersService) DeleteUser(userID int64) *errors.RestErr {
	user := &users.User{ID: userID}
	return user.Delete()
}

// SearchUser gets all of the users with a particular status and returns the list of users with that status
func (s *usersService) SearchUser(status string) (users.Users, *errors.RestErr) {
	dao := &users.User{}
	return dao.FindByStatus(status)
}
