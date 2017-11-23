# peanut_core
A common logging package for Golang

level
=
* 5 levels: info trace error warning success

error
=
* logging when an error occurs and make it colored

color
=
* support 'set default color' function
* make msg colored

write file
=
* by time
* by size(default 1024)
* write file at regular intervals
* support rotate

config
=
* get config from file

Getting Started
=
```go
package main

import (
	"github.com/Jason916/peanut_core/log"
	"fmt"
)

func init() {

	// set default color for logging
	log.SetDefaultColor(log.DefaultColor)

}

func main() {
	// make msg colored
	log.Trace("this is trace tag")

	log.Info("this is info tag")

	log.Success("this is success tag")

	// write log file
	plogWriteConfig := log.NewPLogWriterConfig()
	plogWriter, err := log.NewPLogWriter(".", "test_log", true, plogWriteConfig)

	if err != nil{
		log.Error("write log failed, cause by: ", err)
		return
	}

	plog := log.NewPLogger(plogWriter)
	plog.PInfo("test info ")
	plog.PError("test write error")
	plog.PTrace("test write trace")

	if err := plog.PClose(); err != nil {
		fmt.Println("log close failed, cause by: ", err)
	}
}
