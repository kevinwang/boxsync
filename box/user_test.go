package box

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentUser(t *testing.T) {
	server, client := newTestServerClient("/users/me", `{
		"id": "1234",
		"name": "Example User",
		"login": "user@example.com"
	}`)
	defer server.Close()

	user, err := client.GetCurrentUser()
	assert.NoError(t, err, "Function should not return error")
	assert.Equal(t, "1234", user.ID, "ID should be \"1234\"")
	assert.Equal(t, "Example User", user.Name, "Name should be \"Example User\"")
	assert.Equal(t, "user@example.com", user.Login, "Login should be \"user@example.com\"")
}
