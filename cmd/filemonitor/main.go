package main

import (
	"flag"
	"fmt"
	"os"

	"gitlab-beta.engr.illinois.edu/sp-box/boxsync/filemonitor"
)

func printFileEvent(watchEvent *filemonitor.FileWatchEvent) {
	if watchEvent.Type == filemonitor.EvTypeCreate {
		fileInfo, _ := os.Stat(watchEvent.FilePath)
		if fileInfo.IsDir() {
			fmt.Println("created dir: ", watchEvent.FilePath)
		} else {
			fmt.Println("created file: ", watchEvent.FilePath)
		}

	} else if watchEvent.Type == filemonitor.EvTypeWrite {
		fmt.Println("write to file: ", watchEvent.FilePath)
	} else if watchEvent.Type == filemonitor.EvTypeRemove {
		fmt.Println("remove file: ", watchEvent.FilePath)
	} else if watchEvent.Type == filemonitor.EvTypeRename {
		fmt.Println("rename file: ", watchEvent.FilePath)
	} else if watchEvent.Type == filemonitor.EvTypeChmod {
		fmt.Println("chmod file: ", watchEvent.FilePath)
	}
}

func main() {
	pathPtr := flag.String("path", "./", "file path to watch")
	flag.Parse()

	//simple main to test if filemonitor package works....
	watcher := filemonitor.NewWatcher(printFileEvent)

	watcher.AddAll(*pathPtr)

	for {
		select {
		case fileEvent := <-watcher.FileEventC:
			printFileEvent(&fileEvent)
		case <-watcher.CommonSignalC:
			//for now just handle kill signals
			fmt.Fprintln(os.Stderr, "\n Kill signal triggered, quit...")
			watcher.Close()
			return
		}
	}
}
