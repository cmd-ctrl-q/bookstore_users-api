package app

import (
	"github.com/cmd-ctrl-q/bookstore_utils-go/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default() // point router to gin engine
)

// StartApplication is caleld in main.go to start app
func StartApplication() {

	mapUrls()

	logger.Info("about to launch application...")
	router.Run(":8081") // run router
}
