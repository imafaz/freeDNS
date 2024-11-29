package main

import (
	"flag"
	"fmt"
	"freeDNS/config"
	"freeDNS/database"
	"freeDNS/dnsserver"
	"strconv"

	"github.com/imafaz/logger"
)

func main() {
	database.InitDB(config.GetDBPath())
	logger.Infof("%v %v", config.GetName(), config.GetVersion())

	showHelp := flag.Bool("help", false, "Show help")
	showVersion := flag.Bool("version", false, "Show version")
	dnsServerIP := flag.String("server", "", "Set DNS server listen IP")
	dnsServerPort := flag.Int("port", 0, "Set DNS server listen port")
	domainToAdd := flag.String("adddomain", "", "Add domain")
	ipToAdd := flag.String("addip", "", "Add IP")
	domainToDelete := flag.String("deldomain", "", "Delete domain")
	ipToDelete := flag.String("delip", "", "Delete IP")
	startDnsServer := flag.Bool("start", false, "Start DNS server")
	reversProxyIP := flag.String("proxy", "", "reverse proxy nginx IP")

	debug := flag.Bool("debug", false, "enable debug")

	flag.Parse()

	if *debug {
		config.Debug = true
	}

	if *showHelp {
		flag.Usage()
		return
	}

	if *showVersion {
		fmt.Printf("%v %v", config.GetName(), config.GetVersion())
		return
	}

	if *dnsServerIP != "" {
		database.UpdateConfig("server", *dnsServerIP)
		return
	}
	if *reversProxyIP != "" {
		database.UpdateConfig("revers_proxy_ip", *reversProxyIP)
		return
	}
	if *dnsServerPort != 0 {
		database.UpdateConfig("port", strconv.Itoa(*dnsServerPort))
		return
	}
	if *domainToAdd != "" {
		database.AddDomain(*domainToAdd)
		return
	}
	if *ipToAdd != "" {
		database.AllowIP(*ipToAdd)
		return
	}
	if *domainToDelete != "" {
		database.RemoveDomain(*domainToDelete)
		return
	}
	if *ipToDelete != "" {
		database.RemoveIP(*ipToDelete)
		return
	}
	if *startDnsServer {
		dnsserver.StartDnsServer()
	} else {
		flag.Usage()
		return
	}

}
