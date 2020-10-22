package services

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	configLogType "github.com/StevenLeclerc/mayday-client/types/configLog"
	"github.com/StevenLeclerc/mayday-client/types/messageQueue"
	crunchyTools "github.com/crunchy-apps/crunchy-tools"
)

//TODO USE other method than open file... too much cpu consumption
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

	tiktok := time.Tick(1 * time.Second)
	for _ = range tiktok {
		actualFileSize := getStatOfFile(logConfig.LogFilePath).Size()
		if lastFileSize < actualFileSize {
			fmt.Println("Old", lastFileSize)
			fmt.Println("Actual", actualFileSize)
			buf := make([]byte, actualFileSize-lastFileSize)
			_, errReadAt := file.ReadAt(buf, lastFileSize)
			_ = crunchyTools.HasError(errReadAt, "FileCore - ReadFile - ReadAt", true)
			bufString := string(buf)
			bufStrings := strings.Split(bufString, "\n")
			fmt.Println(bufStrings)
			for _, log := range bufStrings {
				if log != "\n" && log != "" {
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
