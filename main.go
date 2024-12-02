package main

import (
	"freeDNS/config"
	"freeDNS/database"
	"freeDNS/flags"
	"freeDNS/logging"
)

func main() {
	logging.Init()
	database.InitDB(config.GetDBPath())
	flags.ParseFlags()
}
