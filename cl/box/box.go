package main

import (
	"fmt"
	"github.com/urfave/cli"
	"gitlab.engr.illinois.edu/sp-box/boxsync/auth"
	"gitlab.engr.illinois.edu/sp-box/boxsync/box"
	"log"
	"os"
)

func main() {

	httpClient, err := auth.Login()
	if err != nil {
		log.Fatal(err)
	}

	client := box.NewClient(httpClient)

	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "getCurrUser",
			Aliases: []string{"u"},
			Usage:   "Print current user ID",
			Action: func(c *cli.Context) error {
				user, err := client.GetCurrentUser()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Current User Id: " + user.ID)
				return nil
			},
		},
	}

	app.Run(os.Args)
}
