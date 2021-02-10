package users

import (
	"net/http"
	"strconv"

	"github.com/cmd-ctrl-q/bookstore_oauth-go/oauth"
	"github.com/cmd-ctrl-q/bookstore_oauth-go/oauth/errors"
	"github.com/cmd-ctrl-q/bookstore_users-api/domain/users"
	"github.com/cmd-ctrl-q/bookstore_users-api/services"
	"github.com/cmd-ctrl-q/bookstore_utils-go/rest_errors"
	"github.com/gin-gonic/gin"
)

// TestService is a struct for mocking
type TestService struct{}

// TestServiceInterface is a mock for the Service interface to test functions in the users_controller.go
func TestServiceInterface() {
	// services.UsersService.
	// mockService
}

func getUserID(userIDParam string) (int64, *rest_errors.RestErr) {
	userID, userErr := strconv.ParseInt(userIDParam, 10, 64) // param, base, bitSize
	if userErr != nil {
		return 0, rest_errors.NewBadRequestError("user id should be a number")
	}
	return userID, nil
}

// Create handles POST requests and creates a new user based on the data in the request
func Create(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
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
	// call oauth library
	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.Status, err)
		return
	}

	// if caller is still 0, it means the access token does not exist.
	// ie the caller is not authorized to access this scope.
	// if callerID := oauth.GetCallerID(c.Request); callerID == 0 {
	// 	err := errors.RestErr{
	// 		Status:  http.StatusUnauthorized,
	// 		Message: "resource not available",
	// 	}
	// 	c.JSON(err.Status, err)
	// 	return
	// }

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

	// if caller id (owner of the at) is the user id, then display the full/private user
	if oauth.GetCallerID(c.Request) == user.ID {
		c.JSON(http.StatusOK, user.Marshall(false))
		return
	}

	// else validate if request is public or private
	// if X-Public = true, display public user
	// if X-Public = false, display private user
	// c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public") == "true"))
	c.JSON(http.StatusOK, user.Marshall(oauth.IsPublic(c.Request)))
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

// Login logs in user
func Login(c *gin.Context) {
	// get user data from incoming request
	var request users.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user, err := services.UsersService.LoginUser(request)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public") == "true"))
	// c.JSON(http.StatusOK, user) // returns also password
}
