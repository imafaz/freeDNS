package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	ConfigFile     string = "/etc/freeDNS/freeDNS.conf"
	ServerIP       string = "0.0.0.0"
	Port           int    = 53
	DnsIP          string
	AllowedIPFile  string = "ALL"
	DomainsFile    string = "ALL"
	Help           bool   = false
	Debug          bool   = false
	AllowedIPs     map[string]struct{}
	AllowedDomains map[string]struct{}
)

func LoadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "server":
			ServerIP = value
		case "port":
			Port, _ = strconv.Atoi(value)
		case "dns":
			DnsIP = value
		case "allowedip":
			AllowedIPFile = value
		case "domains":
			DomainsFile = value
		case "debug":
			if value == "1" || value == "true" {
				Debug = true
			} else {
				Debug = false
			}

		}
	}

	return scanner.Err()
}

func LoadAllowedIPs(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	AllowedIPs = make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip := strings.TrimSpace(scanner.Text())
		if ip == "" || strings.HasPrefix(ip, "#") {
			continue
		}
		AllowedIPs[ip] = struct{}{}

	}

	return scanner.Err()
}

func LoadDomains(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	AllowedDomains = make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if domain == "" || strings.HasPrefix(domain, "#") {
			continue
		}
		AllowedDomains[domain] = struct{}{}
	}

	return scanner.Err()
}

func PrintUsage() {
	fmt.Println("Usage: freeDNS [options]")
	fmt.Println("Options:")
	fmt.Println("  -server string   Server IP address (default: 0.0.0.0)")
	fmt.Println("  -port int        Port number (default: 53)")
	fmt.Println("  -dns string      IP address to respond with for DNS queries (required)")
	fmt.Println("  -allowedip string  Path to a file containing allowed IP addresses (default: ALL)")
	fmt.Println("  -domains string  Path to a file containing allowed domains (default: ALL)")
	fmt.Println("  -debug           Enable debug mode for verbose output")
	fmt.Println("  -help            Show this help message")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  freeDNS -server 192.168.1.1 -port 53 -dns 172.16.1.1 -allowdip allowedips.list -domains domains.list")

}
