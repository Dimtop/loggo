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

func InitializeHarvester(configPath string, interval int, onLogsParsedFuncs []func(logs []log.Log)) *Harvester {
	h := &Harvester{
		Ticker:            time.NewTicker(time.Second * time.Duration(interval)),
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

			for k, f := range cfg.Files {

				var logs []log.Log
				logsAsStrings := file.ParseFile(k, f.SeekPos)

				for _, logStr := range logsAsStrings {
					logs = append(logs, log.ParseLogFromString(logStr))
				}
				cfg.UpdateConfig(h.ConfigPath, cfg.EmitChangeFileSeekposMessage(h.ConfigPath, k, utils.GetFileSize(k)))

				for _, function := range h.OnLogsParsedFuncs {
					function(logs)
				}
			}

		}
	}
}
