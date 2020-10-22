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
