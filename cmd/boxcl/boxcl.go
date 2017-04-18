package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"

	"gitlab.engr.illinois.edu/sp-box/boxsync/auth"
	"gitlab.engr.illinois.edu/sp-box/boxsync/box"
	"gitlab.engr.illinois.edu/sp-box/boxsync/sync"
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
		{
			Name:    "checkSyncFolder",
			Aliases: []string{"s"},
			Usage:   "Check for Box Sync Folder",
			Action: func(c *cli.Context) error {
				syncRoot, err := sync.GetSyncRootFolder(client)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Sync root: %s %q\n", syncRoot.ID, syncRoot.Name)
				return nil
			},
		},
		{
			Name:    "downloadAll",
			Aliases: []string{"dA"},
			Usage:   "Download All files",
			Action: func(c *cli.Context) error {
				syncRoot, err := sync.GetSyncRootFolder(client)
				if err != nil {
					log.Fatal(err)
				}
				err = sync.DownloadAll(client, syncRoot.ID, sync.LocalSyncRoot)
				if err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
		{
			Name:    "upload",
			Aliases: []string{"up"},
			Usage:   "Upload file",
			Action: func(c *cli.Context) error {
				if len(os.Args) < 3 {
					log.Fatal("Specify filename for upload")
				}
				parentId := "0"
				if len(os.Args) > 3 {
					parentId = os.Args[3]
				}
				file, err := client.UploadFile(os.Args[2], parentId)

				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Upload successful")
				fmt.Println(file)
				return nil
			},
		},
		{
			Name:    "uploadNewVersion",
			Aliases: []string{"upN"},
			Usage:   "Upload existing file with new version",
			Action: func(c *cli.Context) error {
				if len(os.Args) < 4 {
					log.Fatal("Specify fileId for upload & source path")
				}

				file, err := client.UploadFileVersion(os.Args[2], os.Args[3])

				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Upload new version successful")
				fmt.Println(file)
				return nil
			},
		},
		{
			Name:    "watchEvents",
			Aliases: []string{"wE"},
			Usage:   "Output event stream in real time",
			Action: func(c *cli.Context) error {
				url, err := client.GetLongPollURL()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Long poll URL: " + url)

				quit := make(chan struct{})
				events, errs, err := client.GetEventStream(url, box.StreamPositionNow, quit)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Event watching started, CTRL-C to quit")

				for {
					select {
					case event := <-events:
						fmt.Println(event)
					case err := <-errs:
						log.Fatal(err)
					}
				}
				return nil
			},
		},
		{
			Name:  "mkdir",
			Usage: "Create folder",
			Action: func(c *cli.Context) error {
				if len(os.Args) < 3 {
					log.Fatal("Specify folder name")
				}
				parentId := "0"
				if len(os.Args) > 3 {
					parentId = os.Args[3]
				}
				folder, err := client.CreateFolder(os.Args[2], parentId)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Folder created")
				fmt.Println(folder)
				return nil
			},
		},
		{
			Name:    "rmFile",
			Aliases: []string{"rm"},
			Usage:   "Delete file",
			Action: func(c *cli.Context) error {
				if len(os.Args) < 3 {
					log.Fatal("Specify file id")
				}
				err := client.DeleteFile(os.Args[2])
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("File deleted")
				return nil
			},
		},
		{
			Name:    "rmFolder",
			Aliases: []string{"rmdir"},
			Usage:   "Delete folder recursively",
			Action: func(c *cli.Context) error {
				if len(os.Args) < 3 {
					log.Fatal("Specify folder id")
				}
				err := client.DeleteFolder(os.Args[2], true)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Folder deleted")
				return nil
			},
		},
	}

	app.Run(os.Args)
}
