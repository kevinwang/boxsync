package box

import (
	"io/ioutil"
	"net/http"
)

type Client interface {
	Get(endpoint string) ([]byte, error)

	GetCurrentUser() (*User, error)
}

type client struct {
	client *http.Client
}

func NewClient(httpClient *http.Client) Client {
	return &client{
		client: httpClient,
	}
}

func (c *client) Get(endpoint string) ([]byte, error) {
	r, err := c.client.Get("https://api.box.com/2.0" + endpoint)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
