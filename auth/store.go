package auth

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"golang.org/x/oauth2"
)

const (
	storeFileName = ".boxsync_session.json"
)

func saveToken(token *oauth2.Token) error {
	tokJSON, err := json.Marshal(token)
	if err != nil {
		return err
	}

	ioutil.WriteFile(getStoreFilePath(), tokJSON, 0600)

	return nil
}

func loadToken() (*oauth2.Token, error) {
	tokJSON, err := ioutil.ReadFile(getStoreFilePath())
	if err != nil {
		return nil, err
	}

	var token oauth2.Token
	err = json.Unmarshal(tokJSON, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func clearStore() {
	os.Remove(getStoreFilePath())
}

func getStoreFilePath() string {
	return path.Join(os.Getenv("HOME"), storeFileName)
}
