package box

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	defaultAPIBaseURL   = "https://api.box.com/2.0"
	defaultAPIUploadURL = "https://upload.box.com/api/2.0"
)

type Client interface {
	Get(endpointPath string) ([]byte, error)
	GetByURL(url string) ([]byte, error)
	Post(endpointPath, bodyType string, body io.Reader, upload bool) ([]byte, error)
	Delete(endpointPath string) ([]byte, error)
	Options(endpointPath string) ([]byte, error)

	GetCurrentUser() (*User, error)

	CreateFolder(name, parentID string) (*Folder, error)
	GetFolder(id string) (*Folder, error)
	GetFolderContents(id string) (*FolderContents, error)
	DeleteFolder(id string, recursive bool) error

	GetFile(id string) (*File, error)
	DownloadFile(id, destPath string) error
	UploadFile(srcPath, parentID string) (*File, error)
	UploadFileVersion(fileID, srcPath string) (*File, error)
	DeleteFile(id string) error

	GetEvents(streamPosition string) (*EventCollection, error)
	GetLongPollURL() (string, error)
	GetEventStream(longPollURL, streamPosition string, quit <-chan struct{}) (<-chan Event, <-chan error, error)
}

type client struct {
	client           *http.Client
	apiBaseURL       string
	apiUploadBaseURL string
}

func NewClient(httpClient *http.Client) Client {
	return &client{
		client:           httpClient,
		apiBaseURL:       defaultAPIBaseURL,
		apiUploadBaseURL: defaultAPIUploadURL,
	}
}

func (c *client) Get(endpointPath string) ([]byte, error) {
	return c.GetByURL(c.endpointURL(endpointPath))
}

func (c *client) GetByURL(url string) ([]byte, error) {
	r, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return handleResponse(r)
}

func (c *client) Post(endpointPath, bodyType string, body io.Reader, upload bool) ([]byte, error) {
	url := c.endpointURL(endpointPath)
	if upload {
		url = c.uploadEndpointURL(endpointPath)
	}
	r, err := c.client.Post(url, bodyType, body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return handleResponse(r)
}

func (c *client) Delete(endpointPath string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", c.endpointURL(endpointPath), nil)
	if err != nil {
		return nil, err
	}

	r, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return handleResponse(r)
}

func (c *client) Options(endpointPath string) ([]byte, error) {
	req, err := http.NewRequest("OPTIONS", c.endpointURL(endpointPath), nil)
	if err != nil {
		return nil, err
	}

	r, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return handleResponse(r)
}

func (c *client) endpointURL(path string) string {
	return c.apiBaseURL + path
}

func (c *client) uploadEndpointURL(path string) string {
	return c.apiUploadBaseURL + path
}

func handleResponse(r *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if r.StatusCode >= 400 {
		return nil, errors.New(r.Status + " -- " + string(body))
	}

	return body, nil
}
