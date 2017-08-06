package tickspot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type User struct {
	ID        int    `json:"id,omitempty" yaml:"user"`
	FirstName string `json:"first_name,omitempty" yaml:"first_name"`
	LastName  string `json:"last_name,omitempty" yaml:"last_name"`
	Email     string `json:"email,omitempty" yaml:"email"`
	TimeZone  string `json:"timezone,omitempty" yaml:"time_zone"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func (u User) GetID() int {
	return u.ID
}

// GET /users.json
func (t *Tick) GetUsers() ([]*User, error) {
	path := fmt.Sprintf("/users.json")

	var b *bytes.Buffer
	resp, err := t.SendRequest(MethodGET, path, b)
	if err != nil {
		return nil, err
	}

	users := []*User{}
	err = json.Unmarshal(resp.Body, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (t *Tick) GetUserByEmail(email string) (*User, error) {
	users, err := t.GetUsers()
	if err != nil {
		return nil, err
	}

	email = strings.ToLower(email)
	for _, u := range users {
		if email == strings.ToLower(u.Email) {
			return u, nil
		}
	}

	return nil, fmt.Errorf("User with email %s not found", email)
}

//GET /users/deleted.json
func (t *Tick) GetDeletedUsers() ([]*User, error) {
	return []*User{}, nil
}
