package logType

import "time"

type Log struct {
	Message        string    `json:"message"`
	Hostname       string    `json:"hostname"`
	Channels       []string  `json:"channels"`
	LoggedAt       time.Time `json:"loggedAt"`
	FetchLogApiKey string    `json:"fetchLogApiKey"`
	Category       string    `json:"category"`
}
