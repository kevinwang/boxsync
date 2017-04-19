package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"

	"gitlab.engr.illinois.edu/sp-box/boxsync/auth"
	"gitlab.engr.illinois.edu/sp-box/boxsync/box"
	"gitlab.engr.illinois.edu/sp-box/boxsync/cache"
	"gitlab.engr.illinois.edu/sp-box/boxsync/filemonitor"
	//"gitlab.engr.illinois.edu/sp-box/boxsync/sync"
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

	//syncRoot, err := sync.GetSyncRootFolder(client)
	//if err != nil {
	//log.Fatal(err)
	//}
	//fmt.Printf("Sync root: %s %q\n", syncRoot.ID, syncRoot.Name)

	//err = sync.DownloadAll(client, syncRoot.ID, sync.LocalSyncRoot)
	//if err != nil {
	//log.Fatal(err)
	//}

	//events, err := client.GetEvents(box.StreamPositionNow)
	//events, err := client.GetEvents("5912462320132937")
	//if err != nil {
	//log.Fatal(err)
	//}
	//fmt.Println(events)

	/*
		if len(os.Args) < 2 {
			log.Fatal("Specify filename for upload")
		}
		file, err := client.UploadFile(os.Args[1], "0")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Upload successful")
		fmt.Println(file)
	*/

	boxRoot := path.Join(os.Getenv("HOME"), "Box Sync")
	if _, err := os.Stat(boxRoot); os.IsNotExist(err) {
		fmt.Printf("Creating directory %s\n", boxRoot)
		err := os.MkdirAll(boxRoot, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	cache, err := cache.NewCache(client)
	if err != nil {
		log.Fatal(err)
	}

	killSignalC := make(chan os.Signal, 1)
	signal.Notify(killSignalC, os.Interrupt, os.Kill)

	watcher := filemonitor.NewWatcher(func(*filemonitor.FileWatchEvent) {
		_ = cache.RescanLocalTree()
	})
	watcher.AddAll(boxRoot)

	for {
		select {
		case <-watcher.FileEventC:
		case <-killSignalC:
			//for now just handle kill signals
			log.Print("Kill signal triggered, quit...")
			watcher.Close()
			return
		}
	}
}
