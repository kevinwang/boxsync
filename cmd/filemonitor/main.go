package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

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

	killSignalC := make(chan os.Signal, 1)
	signal.Notify(killSignalC, os.Interrupt, os.Kill)

	//simple main to test if filemonitor package works....
	watcher := filemonitor.NewWatcher(printFileEvent)

	watcher.AddAll(*pathPtr)

	toExcludePatterns := [2]string{"**/*.ppp", "**/*.ttt"}
	toExcludeFiles := [2]string{"/home/ani91/Desktop/tmpDir/eee.eee", "/home/ani91/Desktop/tmpDir/kkk.kkk"}
	toExcludeFolders := [2]string{"/home/ani91/Desktop/tmpDir/tmpDir/subTmpdir", "home/ani91/Desktop/tmpDir/excludeDir"}
	watcher.AddExcludePatterns(toExcludePatterns)
	watcher.AddExcludeFiles(toExcludeFiles)
	watcher.AddExcludeFiles(toExcludeFolders)

	for {
		select {
		case fileEvent := <-watcher.FileEventC:
			printFileEvent(&fileEvent)
		case <-killSignalC:
			//for now just handle kill signals
			fmt.Fprintln(os.Stderr, "\n Kill signal triggered, quit...")
			watcher.Close()
			return
		}
	}
}
