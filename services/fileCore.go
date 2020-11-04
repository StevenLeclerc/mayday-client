package services

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/StevenLeclerc/mayday-client/config"
	configLogType "github.com/StevenLeclerc/mayday-client/types/configLog"
	"github.com/StevenLeclerc/mayday-client/types/messageQueue"
	crunchyTools "github.com/crunchy-apps/crunchy-tools"
)

//ReadFile will, depend on logConfig, read and push all the file to mayday backend via the inner queuing system,
// Then, it will check if the file has been modified (every 1 second).
//In this case, everything will be push in the inner queuing system
func ReadFile(chanLog chan messageQueue.MessageQueue, logConfig configLogType.LogConfig) {
	file, errOpen := os.Open(logConfig.LogFilePath)
	defer file.Close()
	if errOpen != nil {
		crunchyTools.HasError(errOpen, "[FileCore]", false)
	}
	logger := crunchyTools.FetchLogger()
	logger.Info.Printf("[FileCore] Start treating: %s\n", logConfig.LogFilePath)

	lastFileSize := getStatOfFile(logConfig.LogFilePath).Size()
	if logConfig.LogAllFile {
		readAllFile(file, chanLog, logConfig)
	}

	readTimer := time.Tick(1 * time.Second)
	for _ = range readTimer {
		config.Debug("[FileCore] Checking File...")
		actualFileSize := getStatOfFile(logConfig.LogFilePath).Size()
		if lastFileSize < actualFileSize {
			config.Debug("[FileCore] File changed")
			buf := make([]byte, actualFileSize-lastFileSize)
			_, errReadAt := file.ReadAt(buf, lastFileSize)
			_ = crunchyTools.HasError(errReadAt, "FileCore - ReadFile - ReadAt", true)
			bufString := string(buf)
			bufStrings := strings.Split(bufString, "\n")
			for _, log := range bufStrings {
				if log != "\n" && log != "" {
					config.Debug(fmt.Sprintf("[FileCore] Message %s inserted to queue", log))
					pushToChan(chanLog, log, logConfig)
				}
			}
			lastFileSize = actualFileSize
		}
	}
}

func pushToChan(chanLog chan messageQueue.MessageQueue, lastLine string, logConfig configLogType.LogConfig) {
	chanLog <- messageQueue.MessageQueue{
		Message:  strings.Split(lastLine, "\n")[0],
		Channels: logConfig.Channels,
		Category: logConfig.Category,
	}
}

func getStatOfFile(filePath string) os.FileInfo {
	stat, errStat := os.Stat(filePath)
	_ = crunchyTools.HasError(errStat, "FileCore - getStatOfFile - Stat", false)
	return stat
}

func readAllFile(file *os.File, chanLog chan messageQueue.MessageQueue, logConfig configLogType.LogConfig) {
	reader := bufio.NewReader(file)
	logger := crunchyTools.FetchLogger()
	config.Debug(fmt.Sprintf("[FileCore][readAllFile] activated for %s", logConfig.LogFilePath))
	for {
		line, errRead := reader.ReadString('\n')
		if errRead != nil {
			if errRead == io.EOF {
				logger.Info.Printf("[FileCore] LogAllFile '%s' Done...\n", logConfig.LogFilePath)
			}
			break
		}
		pushToChan(chanLog, line, logConfig)
	}
	reader.Reset(file)
}
