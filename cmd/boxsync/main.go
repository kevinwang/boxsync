package main

import (
	"fmt"
	"log"

	"gitlab-beta.engr.illinois.edu/sp-box/boxsync/auth"
	"gitlab-beta.engr.illinois.edu/sp-box/boxsync/box"
	"gitlab-beta.engr.illinois.edu/sp-box/boxsync/sync"
)

func main() {
	httpClient, err := auth.Login()
	if err != nil {
		log.Fatal(err)
	}

	client := box.NewClient(httpClient)

	user, err := client.GetCurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user.ID)

	syncRoot, err := sync.GetSyncRootFolder(client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sync root: %s %q\n", syncRoot.ID, syncRoot.Name)

	err = sync.DownloadAll(client, syncRoot.ID, sync.LocalSyncRoot)
	if err != nil {
		log.Fatal(err)
	}
}
