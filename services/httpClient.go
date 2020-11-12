package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/StevenLeclerc/mayday-client/config"
	logType "github.com/StevenLeclerc/mayday-client/types/log"
	crunchyTools "github.com/crunchy-apps/crunchy-tools"
)

func SendLog(logs []logType.Log) bool {
	logger := crunchyTools.FetchLogger()
	if logs != nil {
		appConfig := config.FetchAppConfig()
		logJson, errJson := json.Marshal(logs)
		crunchyTools.HasError(errJson, "Client-MayDay - JsonUnMarshal", true)
		if errJson == nil {
			clientHttp := http.DefaultClient
			clientHttp.Timeout = time.Minute + 10
			r, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/log", appConfig.ServerURL), strings.NewReader(string(logJson)))
			r.Close = true
			r.Header.Add("API_KEY", appConfig.APIKey)
			config.Debug(fmt.Sprintf("[SendLog] headers used: %s", r.Header))
			res, errDo := clientHttp.Do(r)
			crunchyTools.HasError(errDo, "Client-MayDay - Do Request", true)
			if errDo != nil {
				return false
			}
			logger.Info.Printf("[SendLog] Status: %s\n", res.Status)
			res.Body.Close()
			return true
		}
	} else {
		logger.Warn.Printf("[SendLog] No Logs to send.\n")
	}
	return true
}
