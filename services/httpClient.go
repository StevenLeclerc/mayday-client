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

func SendLog(logs []logType.Log) {
	logger := crunchyTools.FetchLogger()
	if logs != nil {
		appConfig := config.FetchAppConfig()
		logJson, errJson := json.Marshal(logs)
		crunchyTools.HasError(errJson, "Client-MayDay - JsonUnMarshal", true)
		if errJson == nil {
			clientHttp := http.DefaultClient
			clientHttp.Timeout = time.Minute + 10
			r, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/log", appConfig.ServerURL), strings.NewReader(string(logJson)))
			r.Close = true
			r.Header.Add("API_KEY", appConfig.APIKey)
			r.Header.Add("Content-Length", string(rune(len(string(logJson))+200000)))
			res, errDo := clientHttp.Do(r)
			crunchyTools.HasError(errDo, "Client-MayDay - Do Request", false)
			res.Body.Close()
			logger.Info.Printf("[SendLog] Status: %s\n", res.Status)
		}
	} else {
		logger.Warn.Printf("[SendLog] Empty PayloadHMAC\n")
	}
}
