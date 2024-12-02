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

	if !database.IPExists(clientIP) && database.GetEnableIPRestrictions() {
		logging.Debugf("Rejected request from %s, reason: not in allowed IPs", clientIP)
		return
	}

	for _, question := range r.Question {
		if question.Qtype == dns.TypeA {
			domain := question.Name
			if database.DomainExists(domain) || !database.GetEnableSpecificDomains() {
				rr := &dns.A{
					Hdr: dns.RR_Header{
						Name:   domain,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    3600,
					},
					A: net.ParseIP(database.GetProxyIP()),
				}
				response.Answer = append(response.Answer, rr)

				logging.Debugf("Responding to %s from %s with ip %s", domain, clientIP, database.GetProxyIP())

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
	server := &dns.Server{Addr: net.JoinHostPort(database.GetServerIP(), database.GetServerPort()), Net: "udp"}
	logger.Infof("Starting DNS server on %s:%s", database.GetServerIP(), database.GetServerPort())
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Failed to start dns server: ", err.Error())
		return
	}
}
