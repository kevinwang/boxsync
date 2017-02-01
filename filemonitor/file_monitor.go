package filemonitor

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	//"time"

	"github.com/fsnotify/fsnotify"
)

//set it for 1 now
const (
	maxEventCount = 1
)

type onFileEventCallback func(*FileWatchEvent)

type FileWatcherState int

type EventType int

const (
	EvTypeCreate EventType = iota
	EvTypeWrite
	EvTypeRemove
	EvTypeRename
	EvTypeChmod
)

type FileWatchEvent struct {
	FilePath string
	Type     EventType
}

type FileWatcher struct {
	//public
	FileEventC chan FileWatchEvent

	//private
	watcher         *fsnotify.Watcher
	triggerInstsMap map[string]*TriggerInst
	mutexLock       sync.Mutex
	callback        onFileEventCallback
	quitC           chan int
	exclude         *Exclude
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

	//public members
	fileWatcher.FileEventC = make(chan FileWatchEvent, maxEventCount+1)

	//private members
	fileWatcher.triggerInstsMap = map[string]*TriggerInst{}
	fileWatcher.callback = callback
	fileWatcher.quitC = make(chan int)
	fileWatcher.exclude = &Exclude{Patterns: make(map[string]bool), Files: make(map[string]bool)}

	//start a thread to watch
	go fileWatcher.startRunning()

	return fileWatcher
}

func (fileWatcher *FileWatcher) AddExcludePatterns(patternsToAdd ...string) {
	for _, toAdd := range patternsToAdd {
		toAdd = filepath.Clean(toAdd)
		fileWatcher.exclude.Patterns[toAdd] = true
	}
}

func (fileWatcher *FileWatcher) AddExcludeFiles(filesToAdd ...string) {
	for _, toAdd := range filesToAdd {
		toAdd = filepath.Clean(toAdd)
		fileWatcher.exclude.Files[toAdd] = true
	}
}

func (fileWatcher *FileWatcher) RemoveExcludePatterns(patternsToRemove ...string) {
	for _, toRemove := range patternsToRemove {
		toRemove = filepath.Clean(toRemove)
		if _, ok := fileWatcher.exclude.Patterns[toRemove]; ok {
			delete(fileWatcher.exclude.Patterns, toRemove)
		}
	}
}

func (fileWatcher *FileWatcher) RemoveExcludeFiles(filesToRemove ...string) {
	for _, toRemove := range filesToRemove {
		toRemove = filepath.Clean(toRemove)
		if _, ok := fileWatcher.exclude.Files[toRemove]; ok {
			delete(fileWatcher.exclude.Files, toRemove)
		}
	}
}

func (fileWatcher *FileWatcher) AddAll(filePath string) {
	//Has not been initialized, do nothing now.
	if fileWatcher.watcher == nil {
		fmt.Fprintln(os.Stderr, "filewatcher is nil, return")
		return
	}

	dirs, err := getSubFolders(filePath)
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

	dirs, err := getSubFolders(filePath)
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

func (fileWatcher *FileWatcher) Close() {
	fileWatcher.quitC <- 0
	err := fileWatcher.watcher.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "watcher Close error:", err)
	}
	fileWatcher.watcher = nil
}

//------------------------------------------
//Private methods start here
//------------------------------------------

func (fileWatcher *FileWatcher) startRunning() {
	for {
		select {
		case <-fileWatcher.quitC:
			fmt.Fprintln(os.Stdout, "watcher closed")
			return
		case fileEvent, ok := <-fileWatcher.watcher.Events:
			if !ok {
				fmt.Fprintln(os.Stderr, "unknown errors")
				return
			}

			//trigger an event only if the file is not excluded
			//if !fileWatcher.exclude.IsMatch(fileEvent.Name) {
			fileWatcher.triggerEvent(&fileEvent)
			//}
		case errorEvent, ok := <-fileWatcher.watcher.Errors:
			if !ok {
				fmt.Fprintln(os.Stderr, errorEvent.Error())
			}
		}
	}
}

func (fileWatcher *FileWatcher) triggerEvent(fileEvent *fsnotify.Event) {
	//fileWatcher.mutexLock.Lock()
	//defer fileWatcher.mutexLock.Unlock()

	/*
		var triggerInst *TriggerInst
		var ok bool
		triggerInst, ok = fileWatcher.triggerInstsMap[fileEvent.Name]
		//first time event firing
		if !ok {
			triggerInst = &TriggerInst{filePath: fileEvent.Name, fileName: filepath.Base(fileEvent.Name), isBusy: false}
			fileWatcher.triggerInstsMap[fileEvent.Name] = triggerInst
		}
	*/
	//start a thread to handle file function.
	fileWatcher.handleFileEvent(nil, fileEvent)
}

func (fileWatcher *FileWatcher) handleFileEvent(triggerInst *TriggerInst, fileEvent *fsnotify.Event) {
	//cannot run this at this point, do nothing.
	//if !triggerInst.canrun() {
	//	return
	//}

	//defer triggerInst.setLastUpdate()

	//wait...
	//timeChannel := time.Tick(minInterval)
	//<-timeChannel

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
		return
	} else if fileEvent.Op&fsnotify.Chmod == fsnotify.Chmod {
		eventType = EvTypeChmod
	} else {
		fmt.Fprintln(os.Stderr, "unknown file event...")
		return
	}

	/*
		We are passing by pointers, but the caller might change the event values
		Therefore, create different FileWatchEvent for each exposured member
	*/

	//handle callback function
	fileWatcher.callback(&FileWatchEvent{FilePath: fileEvent.Name, Type: eventType})

	//push to channel
	if len(fileWatcher.FileEventC) == maxEventCount {
		//discard oldest element in the channel
		<-fileWatcher.FileEventC
	}
	fileWatcher.FileEventC <- FileWatchEvent{FilePath: fileEvent.Name, Type: eventType}
}
