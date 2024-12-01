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
		response.SetRcode(r, dns.RcodeRefused)
		w.WriteMsg(response)
		return
	}

	for _, question := range r.Question {
		if question.Qtype == dns.TypeA {
			domain := question.Name
			logger.Infof("domain exists: %t", database.DomainExists(domain))

			if database.DomainExists(domain) {
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
				ip := resolveDomain(domain)
				rr := &dns.A{
					Hdr: dns.RR_Header{
						Name:   domain,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    3600,
					},
					A: ip,
				}
				response.Answer = append(response.Answer, rr)

				if config.Debug {
					logger.Debugf("Resolved %s to %s", domain, ip.String())
				}

			}
		}
	}

	w.WriteMsg(response)
}

func resolveDomain(domain string) net.IP {
	ips, err := net.LookupIP(domain)
	if err != nil {
		logger.Fatal("Failed to resolve domain: ", err.Error())
	}
	for _, ip := range ips {
		if ip.To4() != nil {
			return ip
		}
	}
	return nil
}

func StartDnsServer() {
	dns.HandleFunc(".", HandleDNSRequest)
	server := &dns.Server{Addr: net.JoinHostPort(database.GetConfig("server"), database.GetConfig("port")), Net: "udp"}
	logger.Infof("Starting DNS server on %s:%s", database.GetConfig("server"), database.GetConfig("port"))
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Failed to start dns server: ", err.Error())
		return
	}
}
