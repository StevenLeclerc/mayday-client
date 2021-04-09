package services

import (
	"time"

	"github.com/StevenLeclerc/mayday-client/config"
	logType "github.com/StevenLeclerc/mayday-client/types/log"
	crunchyTools "github.com/crunchy-apps/crunchy-tools"
)

type QueueHandler struct {
	ChanMessage chan logType.Log
	Queue       []logType.Log
	Status      bool
	Paused      bool
}

var queueHandler QueueHandler

func FetchQueueHandler() *QueueHandler {
	if !queueHandler.Status {
		queueHandler.ChanMessage = make(chan logType.Log)
		queueHandler.Status = true
		queueHandler.Paused = false
	}
	return &queueHandler
}

func (receiver *QueueHandler) InsertPostMessage(log logType.Log) {
	receiver.ChanMessage <- log
	return
}

//Supervisor should be launched one within a go routine
// It will listen to the main Log chan, and every new log pushed, check if the threshold his reached.
// In this case, the api call will be made, and the queued purged.
func (receiver *QueueHandler) Supervisor(chanApiIsOnline chan bool) {
	logger := crunchyTools.FetchLogger()
	logger.Info.Println("[Supervisor] Started!")
	for logFetch := range receiver.ChanMessage {
		if len(receiver.Queue) >= 1000 && receiver.Paused == false {
			config.Debug("[Supervisor] Queue reach 1000 elements")
			if !SendLog(receiver.Queue) && receiver.Paused == false {
				logger.Warn.Printf("[Supervisor] API Problem, sending state to Stabilizer.")
				chanApiIsOnline <- false
			} else {
				chanApiIsOnline <- true
				receiver.CleanQueue()
				config.Debug("[Supervisor] Queue Cleaned")
			}
		}
		config.Debug("[Supervisor] Append message to Queue")
		receiver.Queue = append(receiver.Queue, logFetch)
	}
}

//WakeUpQueue should be launched once within a go routine.
// It will send logs within the queue every x seconds
func (receiver *QueueHandler) WakeUpQueue(chanApiIsOnline chan bool) {
	logger := crunchyTools.FetchLogger()
	logger.Info.Println("[WakeUpQueue] Started!")
	for {
		time.Sleep(time.Second * 2)
		countQueue := len(receiver.Queue)
		if countQueue > 0 && receiver.Paused == false {
			logger.Info.Printf("[WakeUpQueue] Clean needed for: %d messages\n", countQueue)
			if !SendLog(receiver.Queue) && receiver.Paused == false {
				logger.Warn.Printf("[WakeUpQueue] API Problem, sending state to Stabilizer.")
				chanApiIsOnline <- false
			} else {
				chanApiIsOnline <- true
				receiver.CleanQueue()
			}
		}
	}
}

func (receiver *QueueHandler) CleanQueue() {
	receiver.Queue = []logType.Log{}
}

func (receiver *QueueHandler) ForceSendMessages() bool {
	status := false
	countQueue := len(receiver.Queue)
	logger := crunchyTools.FetchLogger()
	logger.Info.Printf("[ForceSendMessages] %d messages pushed\n", countQueue)
	if countQueue > 0 {
		status = SendLog(receiver.Queue)
		if status {
			receiver.CleanQueue()
		}
	}
	return status
}
