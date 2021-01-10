package users

import (
	"net/http"
	"strconv"

	"github.com/cmd-ctrl-q/bookstore_users-api/domain/users"
	"github.com/cmd-ctrl-q/bookstore_users-api/services"
	"github.com/cmd-ctrl-q/bookstore_users-api/utils/errors"
	"github.com/gin-gonic/gin"
)

// CreateUser handles CreateUser requests
func CreateUser(c *gin.Context) {
	var user users.User
	restErr := errors.NewBadRequestError("invalid json body")
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(restErr.Status, restErr)
		return
	}

	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

// GetUser handles GetUser requests
func GetUser(c *gin.Context) {
	userID, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64) // param, base, bitSize
	if userErr != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("user id should be a number"))
		return
	}

	user, getErr := services.GetUser(userID)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	c.JSON(http.StatusOK, user)
}
