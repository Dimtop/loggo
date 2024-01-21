package main

import "loggo/harvester"

func main() {
	h := harvester.InitializeHarvester("config.json")
	h.Run()
}
