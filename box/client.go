package box

import (
	"errors"
	"io/ioutil"
	"net/http"
)

const (
	defaultAPIBaseURL = "https://api.box.com/2.0"
)

type Client interface {
	Get(endpoint string) ([]byte, error)

	GetCurrentUser() (*User, error)
	GetFolderContents(id string) (*FolderContents, error)
	GetFolder(id string) (*Folder, error)
	GetFile(id string) (*File, error)

	DownloadFile(id, destPath string) error
}

type client struct {
	client     *http.Client
	apiBaseURL string
}

func NewClient(httpClient *http.Client) Client {
	return &client{
		client:     httpClient,
		apiBaseURL: defaultAPIBaseURL,
	}
}

func (c *client) Get(endpointPath string) ([]byte, error) {
	r, err := c.client.Get(c.endpointURL(endpointPath))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if r.StatusCode >= 400 {
		return nil, errors.New(r.Status)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *client) endpointURL(path string) string {
	return c.apiBaseURL + path
}
