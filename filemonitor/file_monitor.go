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

type FileWatcherState int

const (
	StateWatching = iota
	StateClosed
)

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
	state           FileWatcherState
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
	fileWatcher.state = StateWatching
	go fileWatcher.startRunning()

	return fileWatcher
}

func (fileWatcher *FileWatcher) AddAll(filePath string) {
	//Has not been initialized, do nothing now.
	if fileWatcher.watcher == nil {
		fmt.Fprintln(os.Stderr, "filewatcher is nil, return")
		return
	}

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
	//Has not been initialized, do nothing now.
	if fileWatcher.watcher == nil {
		fmt.Fprintln(os.Stderr, "filewatcher is nil, return")
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

func (fileWatcher *FileWatcher) RemoveAll(filePath string) {
	//Has not been initialized, do nothing now.
	if fileWatcher.watcher == nil {
		fmt.Fprintln(os.Stderr, "filewatcher is nil, return")
		return
	}

	dirs, err := fileWatcher.getSubFolders(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "directory traversing err:", err)
		return
	}

	for _, dir := range dirs {
		fileWatcher.Remove(dir)
	}
}

func (fileWatcher *FileWatcher) Remove(filePath string) {
	//Has not been initialized, do nothing now.
	if fileWatcher.watcher == nil {
		fmt.Fprintln(os.Stderr, "filewatcher is nil, return")
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

	if fileWatcher.state == StateWatching {
		fileWatcher.Close()
	}
}

func (fileWatcher *FileWatcher) Close() {
	err := fileWatcher.watcher.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "watcher Close error:", err)
	}
	fileWatcher.state = StateClosed
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
		case fileEvent, ok := <-fileWatcher.watcher.Events:
			if !ok {
				if fileWatcher.state == StateClosed {
					fmt.Fprintln(os.Stdout, "watcher Closed")
				} else {
					fmt.Fprintln(os.Stderr, "unknown errors")
				}
				return
			}
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
