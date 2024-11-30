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
	getDomains := flag.Bool("getdomains", false, "show all domains")
	getIPs := flag.Bool("getips", false, "show all allowd ips")

	debug := flag.Bool("debug", false, "enable debug")

	flag.Parse()

	if *debug {
		config.Debug = true
	}
	if *getDomains {
		fmt.Println(database.GetDomains())
	}
	if *getIPs {
		fmt.Println(database.GetIPs())
	}
	if *showHelp {
		flag.Usage()
		return
	}

	if *showVersion || *showVersionShort {
		logger.Infof("%v %v\n", config.GetName(), config.GetVersion())
		return
	}

	if *dnsServerIP != "" {
		database.UpdateConfig("server", *dnsServerIP)
		logger.Infof("DNS server IP updated to: %s", *dnsServerIP)
		return
	}
	if *reversProxyIP != "" {
		database.UpdateConfig("revers_proxy_ip", *reversProxyIP)
		logger.Infof("Reverse proxy IP updated to: %s", *reversProxyIP)
		return
	}
	if *dnsServerPort != 0 {
		database.UpdateConfig("port", strconv.Itoa(*dnsServerPort))
		logger.Infof("DNS server port updated to: %d", *dnsServerPort)
		return
	}
	if *domainToAdd != "" {
		database.AddDomain(*domainToAdd)
		logger.Infof("Domain added: %s", *domainToAdd)
		return
	}
	if *ipToAdd != "" {
		database.AllowIP(*ipToAdd)
		logger.Infof("IP allowed: %s", *ipToAdd)
		return
	}
	if *domainToDelete != "" {
		database.RemoveDomain(*domainToDelete)
		logger.Infof("Domain removed: %s", *domainToDelete)
		return
	}
	if *ipToDelete != "" {
		database.RemoveIP(*ipToDelete)
		logger.Infof("IP removed: %s", *ipToDelete)
		return
	}
	if *startDnsServer {
		dnsserver.StartDnsServer()
		logger.Info("DNS server started")
		return
	}

	if *domainRestrictions == "yes" || *domainRestrictions == "no" {
		database.UpdateConfig("domain_restrictions", *domainRestrictions)
		logger.Infof("Domain Restriction changed to %s", *domainRestrictions)

	} else if *domainRestrictions != "" {
		logger.Fatal("Invalid value for domain_restrictions. Please use 'yes' or 'no'.")
	}

	if *ipRestrictions == "yes" || *ipRestrictions == "no" {
		database.UpdateConfig("ip_restrictions", *ipRestrictions)
		logger.Infof("IP Restriction changed to %s", *ipRestrictions)

	} else if *ipRestrictions != "" {
		logger.Fatal("Invalid value for ip_restrictions. Please use 'yes' or 'no'.")
	}

	if len(flag.Args()) == 0 {
		logger.Fatal("No flags provided. Use -help for usage information.")
		return
	}
}
