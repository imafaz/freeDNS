package flags

import (
	"flag"
	"fmt"
	"freeDNS/config"
	"freeDNS/database"
	"freeDNS/dnsserver"
	"os"

	"github.com/imafaz/logger"
)

func ParseFlags() {
	help := flag.Bool("help", false, "Show help")
	version := flag.Bool("version", false, "Show version")
	shortVersion := flag.Bool("v", false, "Show version (short)")
	debugMode := flag.Bool("debug", false, "enable debug mode")

	dnsServerIP := flag.String("dns-server-ip", "", "Set DNS server listen IP")
	dnsServerPort := flag.Int("dns-server-port", 0, "Set DNS server listen port")
	addDomain := flag.String("add-domain", "", "Add domain")
	addIP := flag.String("add-ip", "", "Add IP")
	deleteDomain := flag.String("delete-domain", "", "Delete domain")
	deleteIP := flag.String("delete-ip", "", "Delete IP")
	startServer := flag.Bool("start-server", false, "Start DNS server")
	reverseProxyIP := flag.String("reverse-proxy-ip", "", "Reverse proxy nginx IP")
	enableSpecificDomains := flag.String("enable-specific-domains", "", "Enable specific domains (yes/no)")
	enableIPRestrictions := flag.String("enable-ip-restrictions", "", "Enable IP restrictions (yes/no)")
	listDomains := flag.Bool("list-domains", false, "Show all domains")
	listIPs := flag.Bool("list-ips", false, "Show all allowed IPs")
	listConfigs := flag.Bool("list-configs", false, "Show all configs")

	flag.Parse()

	if *debugMode {
		config.Debug = true
	}
	if *listDomains {
		domains := database.GetDomains()
		fmt.Println("Domains:")
		for _, domain := range domains {
			fmt.Println(" -", domain)
		}
		os.Exit(0)
	}
	if *listIPs {
		ips := database.GetIPs()
		fmt.Println("Allowed IPs:")
		for _, ip := range ips {
			fmt.Println(" -", ip)
		}
		os.Exit(0)
	}
	if *listConfigs {
		configs := database.GetAllConfig()
		fmt.Println("Configurations:")
		for key, value := range configs {
			fmt.Printf(" - %s: %s\n", key, value)
		}
		os.Exit(0)
	}
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *version || *shortVersion {
		logger.Infof("%v %v\n", config.GetName(), config.GetVersion())
		os.Exit(0)
	}

	if *dnsServerIP != "" {
		database.UpdateServerIP(*dnsServerIP)
		logger.Infof("DNS server IP updated to: %s", *dnsServerIP)
		os.Exit(0)
	}
	if *reverseProxyIP != "" {
		database.UpdateProxyIP(*reverseProxyIP)
		logger.Infof("Proxy IP updated to: %s", *reverseProxyIP)
		os.Exit(0)
	}
	if *dnsServerPort != 0 {
		database.UpdateServerPort(*dnsServerPort)
		logger.Infof("DNS server port updated to: %d", *dnsServerPort)
		os.Exit(0)
	}
	if *addDomain != "" {
		database.AddDomain(*addDomain)
		logger.Infof("Domain added: %s", *addDomain)
		os.Exit(0)
	}
	if *addIP != "" {
		database.AllowIP(*addIP)
		logger.Infof("IP allowed: %s", *addIP)
		os.Exit(0)
	}
	if *deleteDomain != "" {
		database.RemoveDomain(*deleteDomain)
		logger.Infof("Domain removed: %s", *deleteDomain)
		os.Exit(0)
	}
	if *deleteIP != "" {
		database.RemoveIP(*deleteIP)
		logger.Infof("IP removed: %s", *deleteIP)
		os.Exit(0)
	}

	if *enableSpecificDomains == "yes" {
		database.UpdateEnableSpecificDomains(true)
		logger.Infof("Specific domains enabled")
		os.Exit(0)
	} else if *enableSpecificDomains == "no" {
		database.UpdateEnableSpecificDomains(false)
		logger.Infof("Specific domains disabled")
		os.Exit(0)
	} else if *enableSpecificDomains != "" {
		logger.Fatal("Invalid value for enable-specific-domains. Please use 'yes' or 'no'.")
	}

	if *enableIPRestrictions == "yes" {
		database.UpdateEnableIPRestrictions(true)
		logger.Infof("IP restrictions enabled")
		os.Exit(0)
	} else if *enableIPRestrictions == "no" {
		database.UpdateEnableIPRestrictions(false)
		logger.Infof("IP restrictions disabled")
		os.Exit(0)
	} else if *enableIPRestrictions != "" {
		logger.Fatal("Invalid value for enable-ip-restrictions. Please use 'yes' or 'no'.")
	}
	if *startServer {
		dnsserver.StartDnsServer()
	}
	if len(flag.Args()) == 0 {

		logger.Fatal("No flags provided. Use -help for usage information.")
		os.Exit(0)
	}
}
