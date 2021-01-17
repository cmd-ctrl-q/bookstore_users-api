package mysql_utils

import (
	"strings"

	"github.com/cmd-ctrl-q/bookstore_users-api/utils/errors"
	"github.com/go-sql-driver/mysql"
)

const (
	ErrorNoRows = "no rows in result set"
)

// ParseError returns a rest error.
// This function should be used for all errors coming from a mysql database.
func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), ErrorNoRows) {
			return errors.NewNotFoundError("no record matching given id")
		}
		return errors.NewInternalServerError("errors parsing database response")
	}

	switch sqlErr.Number {
	case 1062:
		return errors.NewBadRequestError("invalid data")
	}
	return errors.NewInternalServerError("error processing request")
}
