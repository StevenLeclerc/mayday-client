package services

import (
	"time"

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
			SendLog(receiver.Queue)
			receiver.Queue = []logType.Log{}
		}
		receiver.Queue = append(receiver.Queue, logFetch)
	}
}

//WakeUpQueue should be launched once within a go routine.
// It will check every second if the number of log pushed in Log chan is the same or not.
// If the number is still the same between two ticks, it will clean the queue by sending to the api the stuck logs
func (receiver *QueueHandler) WakeUpQueue() {
	logger := crunchyTools.FetchLogger()
	logger.Info.Println("[WakeUpQueue] Started!")

	for {
		countQueue := len(receiver.Queue)
		time.Sleep(time.Second * 1)
		if countQueue > 0 && countQueue == len(receiver.Queue) {
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
