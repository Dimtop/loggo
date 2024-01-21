package file

import (
	"loggo/config"
	"loggo/log"
	"loggo/utils"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"
)

type File struct {
	File *os.File
	Info os.FileInfo
	Path string
}

func CreateLogFile(configPath string) string {

	cfg := config.ReadConfigFile(configPath)
	var filePath string
	filePath = path.Join(cfg.DirPath, "loggo-"+strconv.FormatInt(time.Now().Unix(), 10)+".txt")
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	f.Close()

	return filePath
}

func WriteLog(logJson string, configPath string) {

	logStr := log.CreateLogFromJson(logJson)

	cfg := config.ReadConfigFile(configPath)

	var filePath string
	if cfg.LastUsedFile != "" {
		filePath = cfg.LastUsedFile
		if cfg.Files[filePath].Size+int64(len(logStr)) > cfg.MaxSize {

			filePath = CreateLogFile(configPath)
			cfg.UpdateConfig(configPath, cfg.EmitChangeLastUsedFileMessage(configPath, filePath))
			cfg.UpdateConfig(configPath, cfg.EmitAddToFilesListMessage(configPath, filePath))

		}
	} else {
		filePath = CreateLogFile(configPath)
		cfg.UpdateConfig(configPath, cfg.EmitChangeLastUsedFileMessage(configPath, filePath))
		cfg.UpdateConfig(configPath, cfg.EmitAddToFilesListMessage(configPath, filePath))

	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic("unable to read file")
	}

	defer f.Close()
	if _, err = f.WriteString(logStr); err != nil {
		panic(err)
	}

	fileStats, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}

	cfg.UpdateConfig(configPath, cfg.EmitUpdateFileSizeMessage(configPath, filePath, fileStats.Size()))

}

func ParseFile(path string, seekPos int64) []string {
	re := regexp.MustCompile(`(LOG)\s+[0-9]{10}\s+(ERROR|SUCCESS|INFO)\s+(PRODUCTION|DEVELOPMENT|STAGING|LOCAL)\s+\S+\s+\S+\s+\S+\s+\S+\s+"(.*?)"\s+"{(.*?)}"`)

	content, err := utils.ReadFileAsString(path, true, seekPos)

	if err != nil {
		panic(err)
	}

	split := re.FindAllString(content, -1)

	return split

}
