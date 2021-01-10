package app

import "github.com/gin-gonic/gin"

var (
	router = gin.Default() // point router to gin engine
)

// StartApplication is caleld in main.go to start app
func StartApplication() {

	mapUrls()           // define maps
	router.Run(":8080") // run router

}
