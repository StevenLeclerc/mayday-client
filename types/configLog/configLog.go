package configLogType

type AppConfig struct {
	APIKey          string      `json:"apiKey"`
	ServerURL       string      `json:"serverUrl"`
	DefaultHostname string      `json:"defaultHostname"`
	LogConfigs      []LogConfig `json:"logConfigs"`
}

type LogConfig struct {
	LogFilePath string   `json:"logFilePath"`
	Channels    []string `json:"channels"`
	LogAllFile  bool     `json:"logAllFile"`
	Category    string   `json:"category"`
}
