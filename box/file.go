package box

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
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

func (c *client) UploadFile(srcPath, parentID string) (*File, error) {
	filename := path.Base(srcPath)
	file, err := os.Open(srcPath)
	if err != nil {
		return nil, err
	}

	fileBody := &bytes.Buffer{}
	writer := multipart.NewWriter(fileBody)

	attr, err := attributesJSON(filename, parentID)
	if err != nil {
		return nil, err
	}

	if err := writer.WriteField("attributes", string(attr)); err != nil {
		return nil, err
	}

	filePart, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(filePart, file); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	respBody, err := c.PostMultipart("/files/content", writer, fileBody)
	if err != nil {
		return nil, err
	}

	return handleUploadResponse(respBody)
}

func (c *client) UploadFileVersion(fileID, srcPath string) (*File, error) {
	filename := path.Base(srcPath)
	file, err := os.Open(srcPath)
	if err != nil {
		return nil, err
	}

	fileBody := &bytes.Buffer{}
	writer := multipart.NewWriter(fileBody)

	filePart, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(filePart, file); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	respBody, err := c.PostMultipart("/files/"+fileID+"/content", writer, fileBody)
	if err != nil {
		return nil, err
	}

	return handleUploadResponse(respBody)
}

func attributesJSON(filename, parentID string) ([]byte, error) {
	attributes := UploadAttributes{
		Name:   filename,
		Parent: UploadParent{ID: parentID},
	}
	attributesJSON, err := json.Marshal(attributes)
	if err != nil {
		return nil, err
	}
	return attributesJSON, nil
}

func handleUploadResponse(body []byte) (*File, error) {
	var collection Collection
	err := json.Unmarshal(body, &collection)
	if err != nil {
		return nil, err
	}
	if collection.Count != 1 {
		return nil, errors.New("Collection count is not 1")
	}
	var fileObj File
	err = json.Unmarshal(collection.Entries[0], &fileObj)
	if err != nil {
		return nil, err
	}

	return &fileObj, nil
}
