package users

import (
	"net/http"
	"strconv"

	"github.com/cmd-ctrl-q/bookstore_users-api/domain/users"
	"github.com/cmd-ctrl-q/bookstore_users-api/services"
	"github.com/cmd-ctrl-q/bookstore_users-api/utils/errors"
	"github.com/gin-gonic/gin"
)

type TestService struct{}

// TestServiceInterface is a mock for the Service interface to test functions in the users_controller.go
func TestServiceInterface() {
	// services.UsersService.
	// mockService
}

func getUserID(userIDParam string) (int64, *errors.RestErr) {
	userID, userErr := strconv.ParseInt(userIDParam, 10, 64) // param, base, bitSize
	if userErr != nil {
		return 0, errors.NewBadRequestError("user id should be a number")
	}
	return userID, nil
}

// Create handles POST requests and creates a new user based on the data in the request
func Create(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	result, saveErr := services.UsersService.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}

// Get handles GET requests and returns a user based on the user id
func Get(c *gin.Context) {
	userID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		c.JSON(http.StatusBadRequest, idErr)
		return
	}

	user, getErr := services.UsersService.GetUser(userID)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public") == "true"))
}

// Update handles PUT and PATCH requests and updates a user based on the user id
func Update(c *gin.Context) {
	// get users id from url
	userID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		c.JSON(http.StatusBadRequest, idErr)
		return
	}

	// the updated user's data coming in through the request
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	// user id is valid. store it into the user object
	user.ID = userID

	// if method = patch
	isPartial := c.Request.Method == http.MethodPatch

	result, err := services.UsersService.UpdateUser(isPartial, user)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

// Delete handles a DELETE request and deletes a user based on user id
func Delete(c *gin.Context) {
	// get users id from url
	userID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		c.JSON(http.StatusBadRequest, idErr)
		return
	}

	if err := services.UsersService.DeleteUser(userID); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

// Search searches for the query parameter 'status'
func Search(c *gin.Context) {
	status := c.Query("status")

	users, err := services.UsersService.SearchUser(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, users.Marshall(c.GetHeader("X-Public") == "true"))
}
