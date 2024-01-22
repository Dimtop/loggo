package main

import (
	"fmt"

	"github.com/Dimtop/loggo/harvester"
	"github.com/Dimtop/loggo/log"
)

func main() {
	/*for i := 0; i < 100; i++ {
		fmt.Println(i)
		file.WriteLog(`{
		    "Timestamp":1705662385,
		    "Context":"platform",
			"Environment":"PRODUCTION",
		    "LogType": "INFO",
		    "Service":"buyers-app",
		    "Feature":"cart",
			"Message" : "Basket was emptied",
		    "Code":"basket-was-emptied",
		    "Data":{
		        "cartId":"134"
		    }

		}`, "C:\\Users\\topal\\Desktop\\Code\\Orderit\\tools\\logs\\config.json")
	}*/

	printLengthFunc := func(logs []log.Log) {
		fmt.Println(len(logs))
	}
	funcs := []func(logs []log.Log){printLengthFunc}

	h := harvester.InitializeHarvester("C:\\Users\\topal\\Desktop\\Code\\Orderit\\tools\\logs\\config.json", 5, funcs)
	h.Run()
}
