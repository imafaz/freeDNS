package config

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/imafaz/logger"
	"github.com/miekg/dns"
)

var (
	ConfigFile     string = ""
	ServerIP       string
	Port           int
	DnsIP          string
	AllowedIPFile  string
	DomainsFile    string
	Help           bool
	Debug          bool = false
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
	fmt.Println("  -conf            Configuration file path")
	fmt.Println("  -debug           Enable debug mode for verbose output")
	fmt.Println("  -help            Show this help message")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  freeDNS -conf /etc/freeDNS/freeDNS.conf")
	fmt.Println("Document:")
	fmt.Println("  https://github.com/imafaz/freeDNS")

}

func ParseConfig() {
	flag.StringVar(&ConfigFile, "conf", ConfigFile, "Config File Path")
	flag.BoolVar(&Debug, "debug", Debug, "Enable debug mode")
	flag.BoolVar(&Help, "help", Help, "Show help message")
	flag.Parse()

	if Help {
		PrintUsage()
		os.Exit(0)
	}

	if ConfigFile == "" {
		logger.Fatalf("Error: The -conf flag is required. Please provide a configuration path file. Use 'freeDNS -help' for assistance.")
		os.Exit(0)
	}

	if err := LoadConfig(ConfigFile); err != nil {
		logger.Fatalf("Error loading configuration file: %s\nError: %s", ConfigFile, err.Error())
	}

	if net.ParseIP(ServerIP).To4() == nil {
		logger.Fatalf("Server ip '%s' not valid!", ServerIP)
	}

	if AllowedIPFile != "ALL" {
		if err := LoadAllowedIPs(AllowedIPFile); err != nil {
			logger.Fatalf("Error loading allowed IPs file: %s\nError: %s", AllowedIPFile, err.Error())
			return
		}
	}

	if DomainsFile != "ALL" {
		if err := LoadDomains(DomainsFile); err != nil {
			logger.Fatalf("Error loading allowed domains file: %s\nError: %s", DomainsFile, err.Error())
			return
		}
	}
}
func HandleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	clientIP, _, _ := net.SplitHostPort(w.RemoteAddr().String())
	response := new(dns.Msg)
	response.SetReply(r)

	if Debug {
		logger.Debugf("Received request from %s for %s", clientIP, r.Question[0].Name)
	}

	if _, exists := AllowedIPs[clientIP]; !exists && AllowedIPFile != "ALL" {
		if Debug {
			logger.Debugf("Rejected request from %s, reason: not in allowed IPs", clientIP)
		}
		return
	}

	for _, question := range r.Question {
		if question.Qtype == dns.TypeA {
			domain := question.Name
			if _, exists := AllowedDomains[domain]; exists || DomainsFile == "ALL" {
				rr := &dns.A{
					Hdr: dns.RR_Header{
						Name:   domain,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    3600,
					},
					A: net.ParseIP(DnsIP),
				}
				response.Answer = append(response.Answer, rr)

				if Debug {
					logger.Debugf("Responding to %s from %s with ip %s", domain, clientIP, DnsIP)
				}
			} else {
				if Debug {
					logger.Debugf("Rejected request from %s for %s , reason: not in allowed domains", clientIP, domain)
				}
			}
		}
	}

	w.WriteMsg(response)
}
