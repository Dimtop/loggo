package harvester

import (
	"time"

	"github.com/Dimtop/loggo/config"
	"github.com/Dimtop/loggo/file"
	"github.com/Dimtop/loggo/log"
	"github.com/Dimtop/loggo/utils"
)

type Harvester struct {
	Ticker            *time.Ticker
	ConfigPath        string
	OnLogsParsedFuncs []func(logs []log.Log)
}

func InitializeHarvester(configPath string, onLogsParsedFuncs []func(logs []log.Log)) *Harvester {
	h := &Harvester{
		Ticker:            time.NewTicker(time.Second * 5),
		ConfigPath:        configPath,
		OnLogsParsedFuncs: onLogsParsedFuncs,
	}

	return h
}

func (h *Harvester) Run() {

	for {
		select {
		case <-h.Ticker.C:

			cfg := config.ReadConfigFile(h.ConfigPath)
			var logs []log.Log
			logsAsStrings := file.ParseFile(cfg.LastUsedFile, cfg.SeekPos)

			for _, logStr := range logsAsStrings {
				logs = append(logs, log.ParseLogFromString(logStr))
			}

			cfg.UpdateConfig(h.ConfigPath, cfg.EmitChangeSeekposMessage(h.ConfigPath, utils.GetFileSize(cfg.LastUsedFile)))

			for _, function := range h.OnLogsParsedFuncs {
				function(logs)
			}

		}
	}
}
