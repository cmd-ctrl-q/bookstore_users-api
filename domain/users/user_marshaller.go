package users

import "encoding/json"

type Users []User

// PublicUser struct is for public requests
type PublicUser struct {
	ID          int64  `json:"id"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
}

// PrivateUser struct is for private/internal requests
type PrivateUser struct {
	ID          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
}

// Marshall marshals user data into json before sending it to the caller
func (users Users) Marshall(isPublic bool) []interface{} {
	result := make([]interface{}, len(users))
	for i, user := range users {
		result[i] = user.Marshall(isPublic)
	}
	return result
}

// Marshall marshals public or private data depending on isPublic value
func (user *User) Marshall(isPublic bool) interface{} {
	if isPublic {
		return &PublicUser{
			ID:          user.ID,
			DateCreated: user.DateCreated,
			Status:      user.Status,
		}
	}

	userJSON, _ := json.Marshal(user)
	var privateUser PrivateUser
	// unmarshal the json into the PrivateUser object
	json.Unmarshal(userJSON, &privateUser)

	return privateUser
}
