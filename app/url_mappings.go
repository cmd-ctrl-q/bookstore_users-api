package app

import (
	"github.com/cmd-ctrl-q/bookstore_users-api/controllers/ping"
	"github.com/cmd-ctrl-q/bookstore_users-api/controllers/users"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)

	// external
	router.POST("/users", users.Create) // internal

	router.GET("/users/:user_id", users.Get)
	router.PUT("/users/:user_id", users.Update)
	router.PATCH("/users/:user_id", users.Update)
	router.DELETE("/users/:user_id", users.Delete)
	router.GET("/internal/users/search", users.Search) // internal
	router.POST("/users/login", users.Login)           // internal

}
