package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cmd-ctrl-q/bookstore_oauth-go/oauth/errors"
	"github.com/cmd-ctrl-q/golang-restclient/rest"
)

const (
	headerXPublic   = "X-Public"
	headerXClientID = "X-Client-Id"
	headerXCallerID = "X-Caller-Id"

	paramAccessToken = "access_token"
)

var (
	oauthRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8080", // oauth api should be running on 8080
		Timeout: 200 * time.Millisecond,
	}
)

type accessToken struct {
	ID       string `json:"is"`
	UserID   int64  `json:"user_id"`
	ClientID int64  `json:"client_id"`
}

// IsPublic validates requests
func IsPublic(request *http.Request) bool {
	// if nil then is public
	if request == nil {
		return true
	}
	return request.Header.Get(headerXPublic) == "true"
}

func GetCallerID(request *http.Request) int64 {
	if request == nil {
		return 0
	}
	callerID, err := strconv.ParseInt(request.Header.Get(headerXCallerID), 10, 64)
	if err != nil {
		return 0
	}
	return callerID
}

func GetClientID(request *http.Request) int64 {
	if request == nil {
		return 0
	}
	clientID, err := strconv.ParseInt(request.Header.Get(headerXClientID), 10, 64)
	if err != nil {
		return 0
	}
	return clientID
}

func AuthenticateRequest(request *http.Request) *errors.RestErr {
	if request == nil {
		return nil
	}

	// clean the header
	cleanRequest(request)

	// fill request with new values
	// e.g. http://api.bookstore.com/resource?access_token=abc123
	accessTokenID := strings.TrimSpace(request.URL.Query().Get(paramAccessToken))
	// dont want to process an access token that is empty
	if accessTokenID == "" {
		return nil
	}

	fmt.Println("accessTokenID: ", accessTokenID)
	at, err := getAccessToken(accessTokenID)
	if err != nil {
		// status code doesn't exists in oauth api
		if err.Status == http.StatusNotFound {
			// return nil because dont want to throw non-oauth api errors.
			return nil
		}
		return err
	}

	// add the headers so we can use them
	request.Header.Add(headerXClientID, fmt.Sprintf("%v", at.ClientID))
	request.Header.Add(headerXCallerID, fmt.Sprintf("%v", at.UserID))

	// request is authenticated
	return nil
}

func cleanRequest(request *http.Request) {
	// anytime you have a parameter as a pointer, always check if its nil
	if request == nil {
		return
	}

	request.Header.Del(headerXClientID)
	request.Header.Del(headerXCallerID)
}

func getAccessToken(accessTokenID string) (*accessToken, *errors.RestErr) {
	// call/get oauth api to get access token. see bookstore_oauth-api/app/application.go
	response := oauthRestClient.Get(fmt.Sprintf("/oauth/access_token/%s", accessTokenID))
	fmt.Println("response: ", response)

	// rest client timeout
	if response == nil || response.Response == nil {
		return nil, errors.NewInternalServerError("invalid restclient response when tyring to get access token")
	}

	// any other error.
	// invalid error whose struct signature doesnt match our restErr fields
	if response.StatusCode > 299 {
		var restErr errors.RestErr
		if err := json.Unmarshal(response.Bytes(), &restErr); err != nil {
			return nil, errors.NewInternalServerError("invalid error interface when trying to get access token")
		}

		return nil, &restErr
	}

	var at accessToken
	if err := json.Unmarshal(response.Bytes(), &at); err != nil {
		return nil, errors.NewInternalServerError("error when trying to unmarshal access token reponse")
	}
	return &at, nil
}
