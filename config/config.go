package config

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/Dimtop/loggo/utils"
)

type Config struct {
	DirPath      string
	LastUsedFile string
	Files        map[string]FileConfig
	MaxSize      int64
}

type FileConfig struct {
	Size    int64
	SeekPos int64
}
type ConfigFileWritesChanMessage struct {
	Action     CONFIG_UPDATE_ACTION
	ConfigData Config
	Identifier string
}
type CONFIG_UPDATE_ACTION string

const (
	CHANGE_LAST_USED_FILE CONFIG_UPDATE_ACTION = "CHANGE_LAST_USED_FILE"
	CHANGE_SEEK_POS       CONFIG_UPDATE_ACTION = "CHANGE_SEEK_POS"
	ADD_TO_FILES_LIST     CONFIG_UPDATE_ACTION = "ADD_TO_FILE_LIST"
	UPDATE_FILE_SIZE      CONFIG_UPDATE_ACTION = "UPDATE_FILE_SIZE"
)

var ConfigFileWritesRWMutex sync.RWMutex

func ReadConfigFile(configPath string) Config {
	content, err := utils.ReadFileAsString(configPath, false, -1)
	if err != nil {
		panic(err)
	}
	var config Config
	json.Unmarshal([]byte(content), &config)
	return config
}

func (c *Config) EmitChangeFileSeekposMessage(configPath string, filePath string, seekPos int64) ConfigFileWritesChanMessage {
	configData := ReadConfigFile(configPath)

	configData.Files[filePath] = FileConfig{
		Size:    configData.Files[filePath].Size,
		SeekPos: seekPos,
	}

	return ConfigFileWritesChanMessage{
		ConfigData: configData,
		Action:     CHANGE_SEEK_POS,
		Identifier: filePath,
	}

}

func (c *Config) EmitChangeLastUsedFileMessage(configPath string, lastUsedFile string) ConfigFileWritesChanMessage {

	configData := ReadConfigFile(configPath)

	configData.LastUsedFile = lastUsedFile

	return ConfigFileWritesChanMessage{
		ConfigData: configData,
		Action:     CHANGE_LAST_USED_FILE,
	}

}

func (c *Config) EmitAddToFilesListMessage(configPath string, fileToAdd string) ConfigFileWritesChanMessage {

	configData := ReadConfigFile(configPath)

	configData.Files[fileToAdd] = FileConfig{
		Size: 0,
	}

	return ConfigFileWritesChanMessage{
		ConfigData: configData,
		Action:     ADD_TO_FILES_LIST,
		Identifier: fileToAdd,
	}

}

func (c *Config) EmitUpdateFileSizeMessage(configPath string, filePath string, size int64) ConfigFileWritesChanMessage {

	configData := ReadConfigFile(configPath)

	configData.Files[filePath] = FileConfig{
		Size: size,
	}

	return ConfigFileWritesChanMessage{
		ConfigData: configData,
		Action:     UPDATE_FILE_SIZE,
		Identifier: filePath,
	}

}

func UpdateFileSize(currentConfigData Config, alteredConfigData Config, filePath string) Config {
	currentConfigData.Files[filePath] = alteredConfigData.Files[filePath]
	return currentConfigData
}

func AddToFilesList(currentConfigData Config, alteredConfigData Config, filePath string) Config {
	currentConfigData.Files[filePath] = FileConfig{
		Size: 0,
	}
	return currentConfigData
}

func ChangeLastUsedFileMessage(currentConfigData Config, alteredConfigData Config) Config {
	currentConfigData.LastUsedFile = alteredConfigData.LastUsedFile
	return currentConfigData
}

func ChangeSeekPos(currentConfigData Config, alteredConfigData Config, filePath string) Config {
	currentConfigData.Files[filePath] = alteredConfigData.Files[filePath]
	return currentConfigData
}

func (c *Config) UpdateConfig(configPath string, updateConfigMessage ConfigFileWritesChanMessage) {

	configData := ReadConfigFile(configPath)

	ConfigFileWritesRWMutex.Lock()
	defer ConfigFileWritesRWMutex.Unlock()
	f, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if updateConfigMessage.Action == UPDATE_FILE_SIZE {
		configData = UpdateFileSize(configData, updateConfigMessage.ConfigData, updateConfigMessage.Identifier)
	}
	if updateConfigMessage.Action == ADD_TO_FILES_LIST {
		configData = AddToFilesList(configData, updateConfigMessage.ConfigData, updateConfigMessage.Identifier)
	}
	if updateConfigMessage.Action == CHANGE_LAST_USED_FILE {
		configData = ChangeLastUsedFileMessage(configData, updateConfigMessage.ConfigData)
	}
	if updateConfigMessage.Action == CHANGE_SEEK_POS {
		configData = ChangeSeekPos(configData, updateConfigMessage.ConfigData, updateConfigMessage.Identifier)
	}
	configStr, err := json.Marshal(configData)
	if err != nil {
		panic(err)
	}

	if err := f.Truncate(0); err != nil {
		panic(err)
	}
	f.Seek(0, 0)

	if _, err = f.WriteString(string(configStr)); err != nil {
		panic(err)
	}

}
