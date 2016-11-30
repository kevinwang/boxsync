package box

import (
	"encoding/json"
)

func (c *client) GetCurrentUser() (*User, error) {
	body, err := c.Get("/users/me")
	if err != nil {
		return nil, err
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
