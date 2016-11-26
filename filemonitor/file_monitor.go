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
	Create EventType = iota
	Write
	Remove
	Rename
	Attrib
)

type FileWatcher struct {
	watcher         *fsnotify.Watcher
	triggerInstsMap map[string]*TriggerInst
	mutexLock       sync.Mutex
	callback        onFileEventCallback
}

type FileWatchEvent struct {
	filepath  string
	eventType EventType
}

//Public methods start here

func NewWatcher(callback onFileEventCallback) *FileWatcher {
	fileWatcher := new(FileWatcher)

	var err error
	fileWatcher.watcher, err = fsnotify.NewWatcher()

	//fsnotify fails...
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fileWatcher = nil
	} else {
		fileWatcher.triggerInstsMap = map[string]*triggerInst{}
		fileWatcher.callback = callback
		//start a thread to watch
		go fileWatcher.startRunning()
	}
}

func (fileWatcher *FileWatcher) Watch(filepath string) {
	//Has not been initlized, do nothing now.
	if fileWatcher.watcher == nil {
		return
	}

	fileWatcher.mutexLock.Lock()
	defer fileWatcher.mutexLock.Unlock()

	err := fileWatcher.watcher.Watch(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tried to watch, but err:", err)
	}
}

func (fileWatcher *FileWatcher) RemoveWatch(filepath string) {
	//Has not been initilized, do nothing now.
	if fileWatcher.watcher == nil {
		return
	}

	fileWatcher.mutexLock.Lock()
	def fileWatcher.mutexLock.Unlock()

	err := fileWatcher.watcher.RemoveWatch(filepath);
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func (fileWatcher *FileWatcher) WaitForKill() {
	onkillChannel := make(chan os.Signal, 1)

	//just try to capture kill, and interrupt signals now...
	signal.Notify(onkillChannel, os.Interrupt, os.Kill)
	<-onkillChannel
	fmt.Fprintln(os.Stderr, "\n kill signal triggerred, exiting...")
}

func (fileWatcher *FileWatcher) Close() {
	//do clean up here
	fileWatcher.watcher.Close()
	fileWatcher.watcher = nil
}

//Private methods start here

func (fileWatcher *FileWatcher) startRunning() {
	//Todo, implement this method
	return
}

func (fileWatcher *FileWatcher) triggerEvent(fileEvent *fsnotify.FileEvent) {
	//Todo, implement this method
	return
}

func (fileWatcher *FileWatcher) handleCallback(triggerInst *TriggerInst, fileEvent *fsnotify.FileEvent) {
	//Todo, implement this method
	return
}
