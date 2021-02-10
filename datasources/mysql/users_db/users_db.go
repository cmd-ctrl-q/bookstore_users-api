package users_db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// Driver
	"github.com/go-sql-driver/mysql"

	"github.com/cmd-ctrl-q/bookstore_utils-go/logger"
)

const (
	mysqlUsersUsername = "mysql_users_username"
	mysqlUsersPassword = "mysql_users_password"
	mysqlUsersHost     = "mysql_users_host"
	mysqlUsersScheme   = "mysql_users_schema"
)

var (
	// Client is the users database (schema: users_db)
	Client *sql.DB

	// os.Setenv("", mysqlUsersUsername)
	// os.Setenv("", mysqlUsersPassword)
	// os.Setenv("", mysqlUsersHost)
	// os.Setenv("", mysqlUsersScheme)

	username = os.Getenv(mysqlUsersUsername)
	password = os.Getenv(mysqlUsersPassword)
	host     = os.Getenv(mysqlUsersHost)
	scheme   = os.Getenv(mysqlUsersScheme)
)

func init() {
	// define datasource name. // user:password@tcp(host)/schema?charset=utf8
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username,
		password,
		host,
		scheme,
	)

	var err error
	Client, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	// ping test database
	if err = Client.Ping(); err != nil {
		panic(err)
	}

	mysql.SetLogger(logger.GetLogger())
	log.Println("database successfully configured")
}
