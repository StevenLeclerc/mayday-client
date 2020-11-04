package config

import (
	"encoding/json"

	configLogType "github.com/StevenLeclerc/mayday-client/types/configLog"
	crunchyTools "github.com/crunchy-apps/crunchy-tools"
)

func FetchAppConfig() configLogType.AppConfig {
	var appConfig configLogType.AppConfig
	appRootPath := crunchyTools.GetApplicationRootFolder()
	//appRootPath := "/Users/stevenleclerc/Documents/DEV/Mayday/mayday-client"
	file := crunchyTools.OpenFile(appRootPath + "/config.json")
	errUnMarshal := json.Unmarshal(crunchyTools.FileToByte(file), &appConfig)
	crunchyTools.HasError(errUnMarshal, "config - FetchAppConfig - UnMarshal", false)
	return appConfig
}

//Debug Will print debug a formatted message if AppConfig.Debug is true
func Debug(message string) {
	if FetchAppConfig().Debug {
		logger := crunchyTools.FetchLogger()
		logger.Warn.Printf("[DEBUG]%s", message)
	}
}
