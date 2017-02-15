package main

import (
	"fmt"
	"log"

	"gitlab-beta.engr.illinois.edu/sp-box/boxsync/auth"
	"gitlab-beta.engr.illinois.edu/sp-box/boxsync/box"
)

func main() {
	httpClient, err := auth.Login()
	if err != nil {
		log.Fatal(err)
	}

	client := box.NewClient(httpClient)
	user, _ := client.GetCurrentUser()
	fmt.Println(user.ID)

	//r, err := httpClient.Get("https://api.box.com/2.0/folders/0")
	/*
		if err != nil {
			log.Fatal(err)
		}
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Println(string(body))
	*/
	//fmt.Println(string("...heheh"))
	r, _ := client.GetFolderEntity("0")
	//r, _ := client.GetFile("10257272849")
	fmt.Println(r.FolderId)

}
