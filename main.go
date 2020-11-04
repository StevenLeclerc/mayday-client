package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/StevenLeclerc/mayday-client/config"
	"github.com/StevenLeclerc/mayday-client/services"
	logType "github.com/StevenLeclerc/mayday-client/types/log"
	"github.com/StevenLeclerc/mayday-client/types/messageQueue"
	crunchyTools "github.com/crunchy-apps/crunchy-tools"
)

//TODO Add to config.json the log threshold
func main() {
	logger := crunchyTools.FetchLogger()
	logger.Info.Println("MayDay Client - v1.1.0")

	var chanLog chan messageQueue.MessageQueue
	chanLog = make(chan messageQueue.MessageQueue)
	queueH := services.FetchQueueHandler()

	go queueH.Supervisor()
	go queueH.WakeUpQueue()
	go run(chanLog, 0)

	appConfig := config.FetchAppConfig()
	for _, logConfig := range appConfig.LogConfigs {
		go services.ReadFile(chanLog, logConfig)
	}
	gracefulHandler()
}

//run should be launched once within a routine.
// It will prepare all the information needed by the api.
// Then, listen the Log chan, received the raw, inject some metadata and push it to the inner queuing system
func run(chanLog chan messageQueue.MessageQueue, id int8) {
	appConf := config.FetchAppConfig()
	hostname, errHost := os.Hostname()
	if appConf.DefaultHostname != "" || errHost != nil {
		hostname = appConf.DefaultHostname
	}
	crunchyTools.HasError(errHost, "Client-MayDay - GetHostname", true)
	if errHost != nil {
		hostname = "NoHostFound"
	}

	queueH := services.FetchQueueHandler()
	for log := range chanLog {
		logToPush := logType.Log{
			Message:          log.Message,
			Hostname:         hostname,
			Channels:         log.Channels,
			LoggedAt:         time.Now(),
			LogFetcherApiKey: appConf.APIKey,
			Category:         log.Category,
		}
		queueH.InsertPostMessage(logToPush)
	}
}

//gracefulHandler will listen for SIGINT/SIGTERM. In this case,
//it will push the last message in queue before leave
func gracefulHandler() {
	log := crunchyTools.FetchLogger()
	log.Info.Println("Graceful handler settled")
	chanSigClose := make(chan os.Signal)
	signal.Notify(chanSigClose, syscall.SIGINT, syscall.SIGTERM)
	for sign := range chanSigClose {
		log.Warn.Printf("Signal: '%s' received.\n", sign)
		queueH := services.FetchQueueHandler()
		queueH.ForceSendMessages()
		close(chanSigClose)
		os.Exit(0)
	}
}
