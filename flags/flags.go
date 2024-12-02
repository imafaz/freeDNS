package flags

import (
	"flag"
	"fmt"
	"freeDNS/config"
	"freeDNS/database"
	"freeDNS/dnsserver"
	"os"
	"strconv"

	"github.com/imafaz/logger"
)

func ParseFlags() {
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
	proxyIP := flag.String("proxyIP", "", "reverse proxy nginx IP")
	specificDomain := flag.String("specific_domains", "", "Set specific domains (yes/no)")
	ipRestrictions := flag.String("ip_restrictions", "", "Set IP restrictions (yes/no)")
	getDomains := flag.Bool("getdomains", false, "show all domains")
	getIPs := flag.Bool("getips", false, "show all allowd ips")
	getConfigs := flag.Bool("getconfigs", false, "show all configs")

	debug := flag.Bool("debug", false, "enable debug")

	flag.Parse()

	if *debug {
		config.Debug = true
	}
	if *getDomains {
		fmt.Println(database.GetDomains())
		os.Exit(0)
	}
	if *getIPs {
		fmt.Println(database.GetIPs())
		os.Exit(0)
	}
	if *getConfigs {
		fmt.Println(database.GetAllConfig())
		os.Exit(0)
	}
	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *showVersion || *showVersionShort {
		logger.Infof("%v %v\n", config.GetName(), config.GetVersion())
		os.Exit(0)
	}

	if *dnsServerIP != "" {
		database.UpdateConfig("server", *dnsServerIP)
		logger.Infof("DNS server IP updated to: %s", *dnsServerIP)
		os.Exit(0)
	}
	if *proxyIP != "" {
		database.UpdateConfig("proxy_ip", *proxyIP)
		logger.Infof("proxy IP updated to: %s", *proxyIP)
		os.Exit(0)
	}
	if *dnsServerPort != 0 {
		database.UpdateConfig("port", strconv.Itoa(*dnsServerPort))
		logger.Infof("DNS server port updated to: %d", *dnsServerPort)
		os.Exit(0)
	}
	if *domainToAdd != "" {
		database.AddDomain(*domainToAdd)
		logger.Infof("Domain added: %s", *domainToAdd)
		os.Exit(0)
	}
	if *ipToAdd != "" {
		database.AllowIP(*ipToAdd)
		logger.Infof("IP allowed: %s", *ipToAdd)
		os.Exit(0)
	}
	if *domainToDelete != "" {
		database.RemoveDomain(*domainToDelete)
		logger.Infof("Domain removed: %s", *domainToDelete)
		os.Exit(0)
	}
	if *ipToDelete != "" {
		database.RemoveIP(*ipToDelete)
		logger.Infof("IP removed: %s", *ipToDelete)
		os.Exit(0)
	}

	if *specificDomain == "yes" || *specificDomain == "no" {
		database.UpdateConfig("specific_domains", *specificDomain)
		logger.Infof("Domain Restriction changed to %s", *specificDomain)
		os.Exit(0)

	} else if *specificDomain != "" {
		logger.Fatal("Invalid value for specific_domains. Please use 'yes' or 'no'.")
	}

	if *ipRestrictions == "yes" || *ipRestrictions == "no" {
		database.UpdateConfig("ip_restrictions", *ipRestrictions)
		logger.Infof("IP Restriction changed to %s", *ipRestrictions)
		os.Exit(0)

	} else if *ipRestrictions != "" {
		logger.Fatal("Invalid value for ip_restrictions. Please use 'yes' or 'no'.")
	}
	if *startDnsServer {
		dnsserver.StartDnsServer()
	}
	if len(flag.Args()) == 0 {
		logger.Fatal("No flags provided. Use -help for usage information.")
		os.Exit(0)
	}
}
