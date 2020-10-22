package config

import (
	"fmt"
	"testing"
)

func TestFetchAppConfig(t *testing.T) {
	config := FetchAppConfig()
	for _, logConf := range config.LogConfigs {
		fmt.Printf("Path: %s | Channel: %s\n", logConf.LogFilePath, logConf.Channels)
	}
}
