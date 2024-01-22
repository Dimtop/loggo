package main

import (
	"fmt"
	"loggo/harvester"
	"loggo/log"
)

func main() {
	printFunc := func(logs []log.Log) {
		fmt.Println(logs)
	}
	printLengthFunc := func(logs []log.Log) {
		fmt.Println(len(logs))
	}
	funcs := []func(logs []log.Log){printFunc, printLengthFunc}

	h := harvester.InitializeHarvester("C:\\Users\\topal\\Desktop\\Code\\Orderit\\tools\\logs\\config.json", funcs)
	h.Run()
}
