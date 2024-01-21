package harvester

import (
	"fmt"
	"loggo/config"
	"loggo/file"
	"loggo/utils"
	"time"
)

type Harvester struct {
	Ticker     *time.Ticker
	ConfigPath string
}

func InitializeHarvester(configPath string) *Harvester {
	h := &Harvester{
		Ticker:     time.NewTicker(time.Second * 5),
		ConfigPath: configPath,
	}

	return h
}

func (h *Harvester) Run() {

	for {
		select {
		case <-h.Ticker.C:

			cfg := config.ReadConfigFile(h.ConfigPath)
			logs := file.ParseFile(cfg.LastUsedFile, cfg.SeekPos)
			fmt.Println(logs)
			cfg.UpdateConfig(h.ConfigPath, cfg.EmitChangeSeekposMessage(h.ConfigPath, utils.GetFileSize(cfg.LastUsedFile)))

		}
	}
}
