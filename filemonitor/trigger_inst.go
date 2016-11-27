package filemonitor

import (
	"sync"
	"time"
)

//set it for 1 second now
const (
	minInterval = (time.Second)
)

type TriggerInst struct {
	filePath       string
	fileName       string
	mutexLock      sync.Mutex
	isBusy         bool
	lastUpdateTime time.Time
	callback       onFileEventCallback
}

func (inst *TriggerInst) canrun() bool {
	inst.mutexLock.Lock()
	defer inst.mutexLock.Unlock()
	if inst.isBusy || time.Now().Sub(inst.lastUpdateTime) < minInterval {
		return false
	}

	inst.isBusy = true
	return true
}

func (inst *TriggerInst) setLastUpdate() {
	inst.mutexLock.Lock()
	defer inst.mutexLock.Unlock()
	inst.lastUpdateTime = time.Now()
	inst.isBusy = false
	return
}
