package users

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cmd-ctrl-q/bookstore_users-api/datasources/mysql/users_db"
	"github.com/cmd-ctrl-q/bookstore_users-api/utils/mysql_utils"
	"github.com/cmd-ctrl-q/bookstore_utils-go/logger"
	"github.com/cmd-ctrl-q/bookstore_utils-go/rest_errors"
)

const (
	queryInsertUser             = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser                = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=?;"
	queryUpdateUser             = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser             = "DELETE FROM users WHERE id=?;"
	queryFindByStatus           = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email=? AND password=? AND status=?;"
)

var (
	usersDB = make(map[int64]*User)
)

// Get attempts to get the users data from the db
func (user *User) Get() *rest_errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to get user statement", err)
		return rest_errors.NewInternalServerError("error when trying to get user", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.ID)
	// populate user fields with the incoming data from the row
	if getErr := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		// user not found
		logger.Error("error when trying to get user by id", getErr)
		return rest_errors.NewInternalServerError("error when trying to get user", errors.New("database error")) // new
		// return mysql_utils.ParseError(getErr) // old
	}
	return nil
}

// Save attempts to save a user in the database
func (user *User) Save() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("error when trying to get user", errors.New("database error")) // new
		// return errors.NewInternalServerError(err.Error()) // old
	}
	defer stmt.Close()

	// add user to db
	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		logger.Error("error when trying to save user", saveErr)
		return rest_errors.NewInternalServerError("error when trying to get user", errors.New("database error")) // new
		// return mysql_utils.ParseError(saveErr) // old
	}

	// get the last row (ie. userID) the user was inserted
	userID, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error")) // new
		// return mysql_utils.ParseError(err) // old
	}
	user.ID = userID
	return nil
}

// Update updates an existing user's fields in the db
func (user *User) Update() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return rest_errors.NewInternalServerError("error when trying to update user", errors.New("database error")) // new
		// return errors.NewInternalServerError(err.Error()) // old
	}
	defer stmt.Close()

	// attempt to update user
	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		logger.Error("error when trying to update user", err)
		return rest_errors.NewInternalServerError("error when trying to update user", errors.New("database error")) // new
		// return mysql_utils.ParseError(err) // old
	}
	return nil
}

// Delete attempts to delete an existing user from the db
func (user *User) Delete() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", errors.New("database error")) // new
		// return errors.NewInternalServerError((err.Error()))
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.ID); err != nil {
		logger.Error("error when trying to delete user", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", errors.New("database error")) // new
		// return mysql_utils.ParseError(err)
	}
	return nil
}

// FindByStatus finds users in the database based on an input status and returns the list of users.
func (user *User) FindByStatus(status string) (Users, *rest_errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find users by status statement", err)
		return nil, rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error"))
	}
	defer stmt.Close()

	// get rows from db
	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find users by status", err)
		return nil, rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error")) // new
		// return nil, mysql_utils.ParseError(err) // old
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when scanning user row into user struct", err)
			return nil, rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error")) // new
			// return nil, mysql_utils.ParseError(err) // old
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}

// FindByEmailAndPassword finds user by their email and password
func (user *User) FindByEmailAndPassword() *rest_errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error"))
	}
	defer stmt.Close()

	// find user by email and password and set their status to active
	result := stmt.QueryRow(user.Email, user.Password, StatusActive)
	// populate user fields with the incoming data from the row
	if getErr := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		// user not found
		if strings.Contains(getErr.Error(), mysql_utils.ErrorNoRows) {
			return rest_errors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when trying to get user by email and password", getErr)
		return rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error")) // new
		// return mysql_utils.ParseError(getErr) // old
	}
	return nil
}
