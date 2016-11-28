package filemonitor

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type onFileEventCallback func(*FileWatchEvent)

type EventType int

const (
	EvTypeCreate EventType = iota
	EvTypeWrite
	EvTypeRemove
	EvTypeRename
	EvTypeChmod
)

type FileWatcher struct {
	watcher         *fsnotify.Watcher
	triggerInstsMap map[string]*TriggerInst
	mutexLock       sync.Mutex
	callback        onFileEventCallback
}

type FileWatchEvent struct {
	FilePath string
	Type     EventType
}

//--------------------------------------
//Public methods start here
//--------------------------------------

func NewWatcher(callback onFileEventCallback) *FileWatcher {
	fileWatcher := new(FileWatcher)

	var err error
	fileWatcher.watcher, err = fsnotify.NewWatcher()

	//fsnotify fails...
	if err != nil {
		fmt.Fprintln(os.Stderr, "calling fsnotify.NewWatcher fails:", err)
		fileWatcher = nil
		return nil
	}

	fileWatcher.triggerInstsMap = map[string]*TriggerInst{}
	fileWatcher.callback = callback
	//start a thread to watch
	go fileWatcher.startRunning()

	return fileWatcher
}

func (fileWatcher *FileWatcher) AddAll(filePath string) {
	dirs, err := fileWatcher.getSubFolders(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "directory travering err:", err)
		return
	}

	for _, dir := range dirs {
		fileWatcher.Add(dir)
	}
}

func (fileWatcher *FileWatcher) Add(filePath string) {
	//Has not been initlized, do nothing now.
	if fileWatcher.watcher == nil {
		return
	}

	fileWatcher.mutexLock.Lock()
	defer fileWatcher.mutexLock.Unlock()

	filePath = filepath.Clean(filePath)

	err := fileWatcher.watcher.Add(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tried to watch:", filePath, ", but err calling fsnotify add:", err)
	}
}

func (fileWatcher *FileWatcher) Remove(filePath string) {
	//Has not been initilized, do nothing now.
	if fileWatcher.watcher == nil {
		return
	}

	fileWatcher.mutexLock.Lock()
	defer fileWatcher.mutexLock.Unlock()

	filePath = filepath.Clean(filePath)

	err := fileWatcher.watcher.Remove(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tried to remove:", filePath, ", but err calling fsnotify remove:", err)
	}
}

func (fileWatcher *FileWatcher) WaitForKill() {
	onSignalKill := make(chan os.Signal, 1)
	signal.Notify(onSignalKill, os.Interrupt, os.Kill)
	<-onSignalKill
	fmt.Fprintln(os.Stderr, "\nKill signal triggered, quit...")
}

func (fileWatcher *FileWatcher) Close() {
	fileWatcher.watcher.Close()
	fileWatcher.watcher = nil
}

//------------------------------------------
//Private methods start here
//------------------------------------------

func (fileWatcher *FileWatcher) getSubFolders(filePath string) (dirs []string, err error) {
	err = filepath.Walk(filePath, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		dirs = append(dirs, newPath)
		return nil

	})

	return dirs, err
}

func (fileWatcher *FileWatcher) startRunning() {
	for {
		select {
		case fileEvent := <-fileWatcher.watcher.Events:
			fileWatcher.triggerEvent(&fileEvent)
		case errorEvent := <-fileWatcher.watcher.Errors:
			fmt.Fprintln(os.Stderr, errorEvent.Error())
		}
	}
}

func (fileWatcher *FileWatcher) triggerEvent(fileEvent *fsnotify.Event) {
	fileWatcher.mutexLock.Lock()
	defer fileWatcher.mutexLock.Unlock()

	var triggerInst *TriggerInst
	var ok bool
	triggerInst, ok = fileWatcher.triggerInstsMap[fileEvent.Name]
	//first time event firing
	if !ok {
		triggerInst = &TriggerInst{filePath: fileEvent.Name, fileName: filepath.Base(fileEvent.Name), isBusy: false}
		fileWatcher.triggerInstsMap[fileEvent.Name] = triggerInst
	}
	//start a thread to handle call back function.
	go fileWatcher.handleCallback(triggerInst, fileEvent)
}

func (fileWatcher *FileWatcher) handleCallback(triggerInst *TriggerInst, fileEvent *fsnotify.Event) {
	//cannot run this at this point, do nothing.
	if !triggerInst.canrun() {
		return
	}

	defer triggerInst.setLastUpdate()

	//wait...
	timeChannel := time.Tick(minInterval)
	<-timeChannel

	var eventType EventType

	//check different file types:
	if fileEvent.Op&fsnotify.Write == fsnotify.Write {
		eventType = EvTypeWrite
	} else if fileEvent.Op&fsnotify.Create == fsnotify.Create {
		eventType = EvTypeCreate
		fileInfo, err := os.Stat(fileEvent.Name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "checking fileinfo err:", err)
			return
		}
		//new directory, watch it too.
		if fileInfo.IsDir() {
			fileWatcher.AddAll(fileEvent.Name)
		}
	} else if fileEvent.Op&fsnotify.Remove == fsnotify.Remove {
		eventType = EvTypeRemove
	} else if fileEvent.Op&fsnotify.Rename == fsnotify.Rename {
		eventType = EvTypeRename
	} else if fileEvent.Op&fsnotify.Chmod == fsnotify.Chmod {
		eventType = EvTypeChmod
	} else {
		fmt.Fprintln(os.Stderr, "unknown file event...")
		return
	}

	fileWatcher.callback(&FileWatchEvent{FilePath: fileEvent.Name, Type: eventType})
}
