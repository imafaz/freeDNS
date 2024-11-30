package main

import (
	"flag"
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
	showVersionShort := flag.Bool("v", false, "Show version (short)")
	dnsServerIP := flag.String("server", "", "Set DNS server listen IP")
	dnsServerPort := flag.Int("port", 0, "Set DNS server listen port")
	domainToAdd := flag.String("adddomain", "", "Add domain")
	ipToAdd := flag.String("addip", "", "Add IP")
	domainToDelete := flag.String("deldomain", "", "Delete domain")
	ipToDelete := flag.String("delip", "", "Delete IP")
	startDnsServer := flag.Bool("start", false, "Start DNS server")
	reversProxyIP := flag.String("proxyIP", "", "reverse proxy nginx IP")
	domainRestrictions := flag.String("domain_restrictions", "", "Set domain restrictions (yes/no)")
	ipRestrictions := flag.String("ip_restrictions", "", "Set IP restrictions (yes/no)")

	debug := flag.Bool("debug", false, "enable debug")

	flag.Parse()

	if *debug {
		config.Debug = true
	}

	if *showHelp {
		flag.Usage()
		return
	}

	if *showVersion || *showVersionShort {
		logger.Debugf("%v %v\n", config.GetName(), config.GetVersion())
		return
	}

	if *dnsServerIP != "" {
		database.UpdateConfig("server", *dnsServerIP)
		logger.Debugf("DNS server IP updated to: %s", *dnsServerIP)
	}
	if *reversProxyIP != "" {
		database.UpdateConfig("revers_proxy_ip", *reversProxyIP)
		logger.Debugf("Reverse proxy IP updated to: %s", *reversProxyIP)
	}
	if *dnsServerPort != 0 {
		database.UpdateConfig("port", strconv.Itoa(*dnsServerPort))
		logger.Debugf("DNS server port updated to: %d", *dnsServerPort)
	}
	if *domainToAdd != "" {
		database.AddDomain(*domainToAdd)
		logger.Debugf("Domain added: %s", *domainToAdd)
	}
	if *ipToAdd != "" {
		database.AllowIP(*ipToAdd)
		logger.Debugf("IP allowed: %s", *ipToAdd)
	}
	if *domainToDelete != "" {
		database.RemoveDomain(*domainToDelete)
		logger.Debugf("Domain removed: %s", *domainToDelete)
	}
	if *ipToDelete != "" {
		database.RemoveIP(*ipToDelete)
		logger.Debugf("IP removed: %s", *ipToDelete)
	}
	if *startDnsServer {
		dnsserver.StartDnsServer()
		logger.Debug("DNS server started")
	}

	if *domainRestrictions == "yes" || *domainRestrictions == "no" {
		database.UpdateConfig("domain_restrictions", *domainRestrictions)
	} else if *domainRestrictions != "" {
		logger.Fatal("Invalid value for domain_restrictions. Please use 'yes' or 'no'.")
	}

	if *ipRestrictions == "yes" || *ipRestrictions == "no" {
		database.UpdateConfig("ip_restrictions", *ipRestrictions)
	} else if *ipRestrictions != "" {
		logger.Fatal("Invalid value for ip_restrictions. Please use 'yes' or 'no'.")
	}

	if len(flag.Args()) == 0 {
		logger.Fatal("No flags provided. Use -help for usage information.")
	}
}
