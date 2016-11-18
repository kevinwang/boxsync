package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gitlab-beta.engr.illinois.edu/sp-box/boxsync/auth"
)

func main() {
	client := auth.Login()

	r, err := client.Get("https://uofi.app.box.com/api/2.0/folders/0")
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
}
