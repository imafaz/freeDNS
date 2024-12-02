package dnsserver

import (
	"freeDNS/database"
	"freeDNS/logging"
	"net"

	"github.com/imafaz/logger"
	"github.com/miekg/dns"
)

func HandleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	clientIP, _, _ := net.SplitHostPort(w.RemoteAddr().String())
	response := new(dns.Msg)
	response.SetReply(r)

	logging.Debugf("Received request from %s for %s", clientIP, r.Question[0].Name)

	if !database.IPExists(clientIP) && database.GetConfig("ip_restrictions") == "yes" {
		logging.Debugf("Rejected request from %s, reason: not in allowed IPs", clientIP)
		return
	}

	for _, question := range r.Question {
		if question.Qtype == dns.TypeA {
			domain := question.Name
			if database.DomainExists(domain) || database.GetConfig("specific_domains") == "no" {
				rr := &dns.A{
					Hdr: dns.RR_Header{
						Name:   domain,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    3600,
					},
					A: net.ParseIP(database.GetConfig("proxy_ip")),
				}
				response.Answer = append(response.Answer, rr)

				logging.Debugf("Responding to %s from %s with ip %s", domain, clientIP, database.GetConfig("proxy_ip"))

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

				logging.Debugf("Resolved %s to %s", domain, ip.String())

			}
		}
	}

	w.WriteMsg(response)
}

func resolveDomain(domain string) net.IP {
	ips, err := net.LookupIP(domain)
	if err != nil {
		if _, ok := err.(*net.DNSError); ok {
			logging.Debug("No such host: ", domain)
		} else {
			logger.Fatal("Error resolving domain: ", err.Error())
		}
		return nil
	}

	for _, ip := range ips {
		if ip.To4() != nil {
			return ip
		}
	}
	logging.Debug("No A record found for domain: ", domain)
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
