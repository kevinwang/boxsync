package box

import (
	"errors"
	"io/ioutil"

	"gitlab-beta.engr.illinois.edu/sp-box/boxsync/httpclient"
)

const (
	defaultApiBaseUrl = "https://api.box.com/2.0"
)

type Client interface {
	Get(endpoint string) ([]byte, error)

	GetCurrentUser() (*User, error)
}

type client struct {
	client     httpclient.Client
	apiBaseUrl string
}

func NewClient(httpClient httpclient.Client) Client {
	return &client{
		client:     httpClient,
		apiBaseUrl: defaultApiBaseUrl,
	}
}

func (c *client) Get(endpointPath string) ([]byte, error) {
	r, err := c.client.Get(c.endpointUrl(endpointPath))
	if err != nil {
		return nil, err
	} else if r.StatusCode >= 400 {
		return nil, errors.New(r.Status)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *client) endpointUrl(path string) string {
	return c.apiBaseUrl + path
}
