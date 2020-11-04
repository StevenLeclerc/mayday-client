package logType

import "time"

type Log struct {
	Message          string    `json:"message"`
	Hostname         string    `json:"hostname"`
	Channels         []string  `json:"channels"`
	LoggedAt         time.Time `json:"loggedAt"`
	LogFetcherApiKey string    `json:"LogFetcherApiKey"`
	Category         string    `json:"category"`
}
