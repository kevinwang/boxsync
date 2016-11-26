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
	//Todo, implement this method
	return
}

func (fileWatcher *FileWatcher) Watch(filepath string) {
	//Todo, implement this method
	return
}

func (fileWatch *FileWatcher) RemoveWatch(filepath string) {
	//Todo, implement this method
	return
}

func (fileWatcher *FileWatcher) WaitForKill() {
	//Todo, implement this method
	return
}

func (fileWatcher *FileWatcher) Close() {
	//Todo, implement this method
	return
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
