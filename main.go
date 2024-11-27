package main

import (
	"freeDNS/config"
	"net"
	"strconv"

	"github.com/imafaz/logger"
	"github.com/miekg/dns"
)

func main() {
	config.ParseConfig()
	dns.HandleFunc(".", config.HandleDNSRequest)
	server := &dns.Server{Addr: net.JoinHostPort(config.ServerIP, strconv.Itoa(config.Port)), Net: "udp"}
	logger.Infof("Starting DNS server on %s:%d", config.ServerIP, config.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Failed to start dns server: ", err.Error())
		return
	}
}
