package mutex

import (
	"encoding/hex"
	"sync"
	"time"
)

type lockObj struct {
	mutex   chan bool
	timeout *time.Timer
}

var locks map[string]lockObj

var masterMutex *sync.Mutex

func init() {
	masterMutex = &sync.Mutex{}
	locks = make(map[string]lockObj)
}

func Lock(pkHash []byte) {
	hashString := hex.EncodeToString(pkHash)
	lock := getLock(hashString)
	select {
	case lock.mutex <- true:
		return
	}
	lock.timeout = time.NewTimer(15 * time.Second)
	go func() {
		<-lock.timeout.C
		Unlock(pkHash)
	}()
}

func Unlock(pkHash []byte) {
	hashString := hex.EncodeToString(pkHash)
	if ! hasLock(hashString) {
		return
	}
	lock := getLock(hashString)
	if lock.timeout != nil {
		lock.timeout.Stop()
		lock.timeout = nil
	}
	<-lock.mutex
	go func() {
		cleanupTimer := time.NewTimer(5 * time.Second)
		<-cleanupTimer.C
		lock = getLock(hashString)
		masterMutex.Lock()
		if lock.timeout == nil {
			delete(locks, hashString)
		}
		masterMutex.Unlock()
	}()
}

func hasLock(hashString string) bool {
	_, ok := locks[hashString]
	return ok
}

func getLock(hashString string) lockObj {
	masterMutex.Lock()
	if _, ok := locks[hashString]; !ok {
		locks[hashString] = lockObj{
			mutex: make(chan bool, 1),
		}
	}
	masterMutex.Unlock()
	return locks[hashString]
}
