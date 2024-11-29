package dnsserver

import (
	"freeDNS/config"
	"freeDNS/database"
	"net"

	"github.com/imafaz/logger"
	"github.com/miekg/dns"
)

func HandleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	clientIP, _, _ := net.SplitHostPort(w.RemoteAddr().String())
	response := new(dns.Msg)
	response.SetReply(r)

	if config.Debug {
		logger.Debugf("Received request from %s for %s", clientIP, r.Question[0].Name)
	}

	if !database.IPExists(clientIP) && database.GetConfig("ip_restrictions") == "yes" {
		logger.Debugf("Rejected request from %s, reason: not in allowed IPs", clientIP)
		return
	}

	for _, question := range r.Question {
		if question.Qtype == dns.TypeA {
			domain := question.Name
			if database.DomainExists(domain) || database.GetConfig("domain_restrictions") == "yes" {
				rr := &dns.A{
					Hdr: dns.RR_Header{
						Name:   domain,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    3600,
					},
					A: net.ParseIP(database.GetConfig("revers_proxy_ip")),
				}
				response.Answer = append(response.Answer, rr)

				if config.Debug {
					logger.Debugf("Responding to %s from %s with ip %s", domain, clientIP, database.GetConfig("revers_proxy_ip"))
				}
			} else {
				if config.Debug {
					logger.Debugf("Rejected request from %s for %s , reason: not in allowed domains", clientIP, domain)
				}
			}
		}
	}

	w.WriteMsg(response)
}

func StartDnsServer() {
	dns.HandleFunc(".", HandleDNSRequest)
	server := &dns.Server{Addr: net.JoinHostPort(database.GetConfig("server"), database.GetConfig("port")), Net: "udp"}
	logger.Infof("Starting DNS server on %s:%s", database.GetConfig("server"), database.GetConfig("server"))
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Failed to start dns server: ", err.Error())
		return
	}
}
