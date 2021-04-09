package services

import (
	"reflect"
	"sync"
	"time"

	crunchyTools "github.com/crunchy-apps/crunchy-tools"
)

//Stabilizer will listen to chanApi, this chan will received false if any problem occurs during the http.post
//in this case, the Stabilizer will block all file readers, set Queue.Paused to true then set the retrying process.
func Stabilizer(chanApi chan bool, readerMutexes []*sync.Mutex) {
	logger := crunchyTools.FetchLogger()
	logger.Info.Println("[Stabilizer] Started !")
	for chanApiIsResponding := range chanApi {
		if chanApiIsResponding == false {
			logger.Warn.Println("[Stabilizer] Problem detected... Locking File Readers.")
			queueHandler := FetchQueueHandler()
			queueHandler.Paused = true
			for _, mutex := range readerMutexes {
				mutex.Lock()
			}
			logger.Warn.Println("[Stabilizer] Blocking File Readers Done.")
			go RetryingSendQueue(readerMutexes)
		}
	}
}

//isMutexLocked check if the mutex is locked or not, using reflection
func IsMutexLocked(m *sync.Mutex) bool {
	const mutexLocked = 1
	state := reflect.ValueOf(m).Elem().FieldByName("state")
	return state.Int()&mutexLocked == mutexLocked
}

//RetryingSendQueue will be triggered by the Stabilizer. It will retry every n seconds (default 10) to push the queue.
//If it succeed, it will unlock file readers. And set the Queue.Paused to false.
func RetryingSendQueue(readerMutexes []*sync.Mutex) {
	logger := crunchyTools.FetchLogger()
	queueHandler := FetchQueueHandler()
	timerTick := time.Tick(10 * time.Second)
	for range timerTick {
		logger.Warn.Println("[Stabilizer][RetryingSendQueue] Sending...")
		if queueHandler.ForceSendMessages() {
			logger.Warn.Println("[Stabilizer][RetryingSendQueue] Queue Sent.")
			logger.Warn.Println("[Stabilizer][RetryingSendQueue] Unlocking File Readers...")
			for _, mutex := range readerMutexes {
				if IsMutexLocked(mutex) {
					mutex.Unlock()
				}
			}
			logger.Warn.Println("[Stabilizer][RetryingSendQueue] Unlocking File Readers Done.")
			queueHandler.Paused = false
			break
		} else {
			logger.Warn.Println("[Stabilizer][RetryingSendQueue] API still down, retrying within 10 secs")
		}
	}
}
