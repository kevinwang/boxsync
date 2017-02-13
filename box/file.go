package box

import (
	"encoding/json"
)

func (c *client) GetFile(fileID string) (*File, error) {
	body, err := c.Get("/files/" + fileID)
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
