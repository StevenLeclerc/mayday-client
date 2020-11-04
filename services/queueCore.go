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
}

var queueHandler QueueHandler

func FetchQueueHandler() QueueHandler {
	if !queueHandler.Status {
		queueHandler.ChanMessage = make(chan logType.Log)
		queueHandler.Status = true
	}
	return queueHandler
}

func (receiver *QueueHandler) InsertPostMessage(log logType.Log) {
	receiver.ChanMessage <- log
	return
}

//Supervisor should be launched one within a go routine
// It will listen to the main Log chan, and every new log pushed, check if the threshold his reached.
// In this case, the api call will be made, and the queued purged.
func (receiver *QueueHandler) Supervisor() {
	logger := crunchyTools.FetchLogger()
	logger.Info.Println("[Supervisor] Started!")
	for logFetch := range receiver.ChanMessage {
		if len(receiver.Queue) >= 1000 {
			config.Debug("[Supervisor] Queue reach 1000 elements")
			SendLog(receiver.Queue)
			receiver.Queue = []logType.Log{}
			config.Debug("[Supervisor] Queue Cleaned")
		}
		config.Debug("[Supervisor] Append message to Queue")
		receiver.Queue = append(receiver.Queue, logFetch)
	}
}

//WakeUpQueue should be launched once within a go routine.
// It will send logs within the queue every x seconds
func (receiver *QueueHandler) WakeUpQueue() {
	logger := crunchyTools.FetchLogger()
	logger.Info.Println("[WakeUpQueue] Started!")
	for {
		time.Sleep(time.Second * 2)
		countQueue := len(receiver.Queue)
		if countQueue > 0 {
			logger.Info.Printf("[WakeUpQueue] Clean needed for: %d messages\n", countQueue)
			SendLog(receiver.Queue)
			receiver.Queue = []logType.Log{}
		}
	}
}
func (receiver *QueueHandler) ForceSendMessages() {
	countQueue := len(receiver.Queue)
	logger := crunchyTools.FetchLogger()
	logger.Info.Printf("[ForceSendMessages] %d messages pushed\n", countQueue)
	if countQueue > 0 {
		SendLog(receiver.Queue)
	}
}
