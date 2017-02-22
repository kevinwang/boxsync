package box

import (
	"encoding/json"
	"io"
	"os"
	"path"
)

func (c *client) GetFile(id string) (*File, error) {
	body, err := c.Get("/files/" + id)
	if err != nil {
		return nil, err
	}
	var file File
	err = json.Unmarshal(body, &file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (c *client) DownloadFile(id, destPath string) error {
	dir := path.Dir(destPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	r, err := c.client.Get(c.endpointURL("/files/" + id + "/content"))
	if err != nil {
		return err
	}
	defer r.Body.Close()

	_, err = io.Copy(out, r.Body)

	return nil
}
