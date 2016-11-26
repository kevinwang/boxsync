package filemonitor

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	minInterval = (time.Second)
)

type TriggerInst struct {
	filepath       string
	filename       string
	mutexLock      sync.Mutex
	isBusy         bool
	lastUpdateTime time.time
	callback       onFileEventCallback
}

func (inst *TriggerInst) canrun() bool {
	//Todo, implement this method
	return
}

func (inst *TriggerInst) setLastUpdate() {
	//Todo, implement this method
	return
}
