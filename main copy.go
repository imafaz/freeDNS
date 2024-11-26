package main

import (
	"flag"
	"freeDNS/config"
	"net"
	"os"
	"strconv"

	"github.com/imafaz/logger"
	"github.com/miekg/dns"
)

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	clientIP, _, _ := net.SplitHostPort(w.RemoteAddr().String())
	response := new(dns.Msg)
	response.SetReply(r)

	if config.Debug {
		logger.Log(logger.DEBUG, "Received request from ", clientIP, " for \n", r.Question[0].Name)
	}

	if _, exists := config.AllowedIPs[clientIP]; !exists && config.AllowedIPFile != "ALL" {
		if config.Debug {
			logger.Log(logger.DEBUG, "Rejected request from ", clientIP, ": not in allowed IPs")
		}
		return
	}

	for _, question := range r.Question {
		if question.Qtype == dns.TypeA {
			domain := question.Name
			if _, exists := config.AllowedDomains[domain]; exists || config.DomainsFile == "ALL" {
				rr := &dns.A{
					Hdr: dns.RR_Header{
						Name:   domain,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    3600,
					},
					A: net.ParseIP(config.DnsIP),
				}
				response.Answer = append(response.Answer, rr)

				if config.Debug {
					logger.Log(logger.INFO, "Responding to ", domain, " from ", clientIP, " with IP ", config.DnsIP)
				}
			} else {
				if config.Debug {
					logger.Log(logger.DEBUG, "Rejected request for ", domain, ": not in allowed domains")
				}
			}
		}
	}

	w.WriteMsg(response)
}

func main() {
	logger.Init("app.log", logger.CONSOLE_ONLY)

	_, err := os.Stat("config.ConfigFile")
	if err != nil {
		if err := config.LoadConfig(config.ConfigFile); err != nil {
			logger.Log(logger.FATAL, "Error loading config: \n", err.Error())
			return
		}
	} else {
		logger.Log(logger.INFO, "config file not found cheking flags")

	}

	flag.StringVar(&config.ServerIP, "server", config.ServerIP, "Server IP")
	flag.IntVar(&config.Port, "port", config.Port, "Port")
	flag.StringVar(&config.DnsIP, "dns", config.DnsIP, "DNS response IP")
	flag.StringVar(&config.AllowedIPFile, "allowedip", config.AllowedIPFile, "Allowed IPs file")
	flag.StringVar(&config.DomainsFile, "domains", config.DomainsFile, "Allowed domains file")
	flag.BoolVar(&config.Debug, "debug", config.Debug, "Enable debug mode")
	flag.BoolVar(&config.Help, "help", config.Help, "Show help message")
	flag.Parse()

	if config.Help {
		config.PrintUsage()
		return
	}

	if config.DnsIP == "" {
		config.PrintUsage()
		return
	}

	if config.AllowedIPFile != "ALL" {
		if err := config.LoadAllowedIPs(config.AllowedIPFile); err != nil {
			logger.Log(logger.FATAL, "Error loading allowed IPs: ", err.Error())
			return
		}
	}

	if config.DomainsFile != "ALL" {
		if err := config.LoadDomains(config.DomainsFile); err != nil {
			logger.Log(logger.FATAL, "Error loading allowed domains: ", err.Error())
			return
		}
	}

	dns.HandleFunc(".", handleDNSRequest)

	server := &dns.Server{Addr: net.JoinHostPort(config.ServerIP, strconv.Itoa(config.Port)), Net: "udp"}
	logger.Log(logger.INFO, "Starting DNS server on ", config.ServerIP, ":", strconv.Itoa(config.Port))

	if err := server.ListenAndServe(); err != nil {
		logger.Log(logger.FATAL, "Failed to start server: ", err.Error())
		return
	}
}
