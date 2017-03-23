package box

import (
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

const (
	defaultAPIBaseURL   = "https://api.box.com/2.0"
	defaultAPIUploadURL = "https://upload.box.com/api/2.0"
)

type Client interface {
	Get(endpoint string) ([]byte, error)
	PostMultipart(endpointPath string, writer *multipart.Writer, body io.Reader) ([]byte, error)

	GetCurrentUser() (*User, error)
	GetFolderContents(id string) (*FolderContents, error)
	GetFolder(id string) (*Folder, error)
	GetFile(id string) (*File, error)
	GetEvents(streamPosition string) (*EventCollection, error)

	DownloadFile(id, destPath string) error
	UploadFile(srcPath, parentID string) (*File, error)
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
	r, err := c.client.Get(c.endpointURL(endpointPath))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return handleResponse(r)
}

func (c *client) PostMultipart(endpointPath string, writer *multipart.Writer, body io.Reader) ([]byte, error) {
	r, err := c.client.Post(c.uploadEndpointURL(endpointPath), writer.FormDataContentType(), body)
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
