package database

import (
	"database/sql"
	"freeDNS/config"
	"freeDNS/logging"
	"strings"

	"github.com/imafaz/logger"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB(dbPath string) {
	logging.Debug("initial database")
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Fatal("Error opening database: ", err.Error())
	}
	createTables()
	insertConfig()
}

func createTables() {
	logging.Debug("creating tables")
	var err error
	var query string
	query = `CREATE TABLE IF NOT EXISTS specific_domains (
        domain TEXT NOT NULL PRIMARY KEY
    );`
	_, err = db.Exec(query)
	if err != nil {
		logger.Fatal("Error creating specific_domains table: ", err.Error())
	}

	query = `CREATE TABLE IF NOT EXISTS config (
    key TEXT NOT NULL PRIMARY KEY,
    value TEXT NOT NULL
);`
	_, err = db.Exec(query)
	if err != nil {
		logger.Fatal("Error creating config table: ", err.Error())
	}

	query = `CREATE TABLE IF NOT EXISTS ip_whitelist (
		ip TEXT NOT NULL PRIMARY KEY
		);`
	_, err = db.Exec(query)
	if err != nil {
		logger.Fatal("Error creating ip_whitelist table: ", err.Error())
	}
}

func insertConfig() {
	logging.Debug("inserting configs")

	var query string
	var err error

	configs := map[string]string{
		"enable_ip_restrictions":  "yes",
		"enable_specific_domains": "yes",
		"proxy_ip":                config.GetServerIP(),
		"server_ip":               "0.0.0.0",
		"server_port":             "53",
	}

	for key, value := range configs {
		query = `SELECT COUNT(*) FROM config WHERE key = ?`
		var count int
		err = db.QueryRow(query, key).Scan(&count)
		if err != nil {
			logger.Fatal("Error checking config key: ", err.Error())
		}
		if count == 0 {
			query = `INSERT INTO config (key, value) VALUES (?, ?)`
			_, err = db.Exec(query, key, value)
			if err != nil {
				logger.Fatal("Error inserting config key: ", err.Error())
			}
		}
	}
}

func UpdateServerPort(value int) {
	logging.Debugf("updating server port to %s", value)
	query := `UPDATE config SET value = ? WHERE key = 'server_port'`
	_, err := db.Exec(query, value)
	if err != nil {
		logger.Fatal("Error updating server port: ", err.Error())
	}
}

func UpdateServerIP(value string) {
	logging.Debugf("updating server IP to %s", value)
	query := `UPDATE config SET value = ? WHERE key = 'server_ip'`
	_, err := db.Exec(query, value)
	if err != nil {
		logger.Fatal("Error updating server IP: ", err.Error())
	}
}

func UpdateProxyIP(value string) {
	logging.Debugf("updating proxy IP to %s", value)
	query := `UPDATE config SET value = ? WHERE key = 'proxy_ip'`
	_, err := db.Exec(query, value)
	if err != nil {
		logger.Fatal("Error updating proxy IP: ", err.Error())
	}
}

func UpdateEnableSpecificDomains(value bool) {
	logging.Debugf("updating enable_specific_domains to %t", value)
	var strValue string
	if value {
		strValue = "yes"
	} else {
		strValue = "no"
	}
	query := `UPDATE config SET value = ? WHERE key = 'enable_specific_domains'`
	_, err := db.Exec(query, strValue)
	if err != nil {
		logger.Fatal("Error updating enable_specific_domains: ", err.Error())
	}
}

func UpdateEnableIPRestrictions(value bool) {
	logging.Debugf("updating enable_ip_restrictions to %t", value)
	var strValue string
	if value {
		strValue = "yes"
	} else {
		strValue = "no"
	}
	query := `UPDATE config SET value = ? WHERE key = 'enable_ip_restrictions'`
	_, err := db.Exec(query, strValue)
	if err != nil {
		logger.Fatal("Error updating enable_ip_restrictions: ", err.Error())
	}
}

func AddDomain(domain string) {
	logging.Debugf("adding domain %s to specific domains", domain)
	query := `INSERT INTO specific_domains (domain) VALUES (?)`
	_, err := db.Exec(query, domain)
	if err != nil {
		logger.Fatal("Error adding domain: ", err.Error())
	}
}
func AllowIP(IP string) {
	logging.Debugf("allowing IP %s", IP)
	query := `INSERT INTO ip_whitelist (ip) VALUES (?)`
	_, err := db.Exec(query, IP)
	if err != nil {
		logger.Fatal("Er IP: ", err.Error())
	}
}

func RemoveDomain(domain string) {
	logging.Debugf("removing domain %s", domain)
	query := `DELETE FROM specific_domains WHERE domain = ?`
	_, err := db.Exec(query, domain)
	if err != nil {
		logger.Fatal("Error removing domain: ", err.Error())
	}
}

func RemoveIP(IP string) {
	logging.Debugf("removing IP %s", IP)
	query := `DELETE FROM ip_whitelist WHERE ip = ?`
	_, err := db.Exec(query, IP)
	if err != nil {
		logger.Fatal("Error removing IP: ", err.Error())
	}
}

func DomainExists(domain string) bool {
	domain = strings.TrimSuffix(domain, ".")
	logging.Debugf("checking if domain %s exists", domain)
	query := `SELECT EXISTS(SELECT 1 FROM specific_domains WHERE domain = ?)`
	var exists bool
	err := db.QueryRow(query, domain).Scan(&exists)
	if err != nil {
		logger.Fatal("Error checking if domain exists: ", err.Error())
	}
	return exists
}

func IPExists(IP string) bool {
	logging.Debugf("checking if IP %s exists", IP)
	query := `SELECT EXISTS(SELECT 1 FROM ip_whitelist WHERE ip = ?)`
	var exists bool
	err := db.QueryRow(query, IP).Scan(&exists)
	if err != nil {
		logger.Fatal("Error checking if IP exists: ", err.Error())
	}
	return exists
}

func GetAllConfig() map[string]string {
	logging.Debug("getting all configurations")
	configs := make(map[string]string)
	query := `SELECT key, value FROM config`
	rows, err := db.Query(query)
	if err != nil {
		logger.Fatal("Error getting configs: ", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			logger.Fatal("Error getting configs row: ", err.Error())

		}
		configs[key] = value
	}

	if err := rows.Err(); err != nil {
		logger.Fatal("Error getting configs raws: ", err.Error())
	}

	return configs
}

func GetDomains() []string {
	logging.Debug("getting all domains")
	var domains []string
	query := `SELECT domain FROM specific_domains`
	rows, err := db.Query(query)
	if err != nil {
		logger.Fatal("Error getting domains: ", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var domain string
		if err := rows.Scan(&domain); err != nil {
			logger.Fatal("Error scanning domain: ", err.Error())
		}
		domains = append(domains, domain)
	}

	if err := rows.Err(); err != nil {
		logger.Fatal("Error iterating over domains: ", err.Error())
	}

	return domains
}

func GetIPs() []string {
	logging.Debug("getting all allowed IPs")
	var ips []string
	query := `SELECT ip FROM specific_domains`
	rows, err := db.Query(query)
	if err != nil {
		logger.Fatal("Error getting allowed IPs: ", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var ip string
		if err := rows.Scan(&ip); err != nil {
			logger.Fatal("Error scanning IP: ", err.Error())
		}
		ips = append(ips, ip)
	}

	if err := rows.Err(); err != nil {
		logger.Fatal("Error iterating over allowed IPs: ", err.Error())
	}

	return ips
}

func GetServerPort() string {
	logging.Debug("getting server port")
	var value string
	query := `SELECT value FROM config WHERE key = 'server_port'`
	err := db.QueryRow(query).Scan(&value)
	if err != nil {
		logger.Fatal("Error getting server port: ", err.Error())
	}
	return value
}
func GetServerIP() string {
	logging.Debug("getting server port")
	var value string
	query := `SELECT value FROM config WHERE key = 'server_ip'`
	err := db.QueryRow(query).Scan(&value)
	if err != nil {
		logger.Fatal("Error getting server ip: ", err.Error())
	}
	return value
}
func GetProxyIP() string {
	logging.Debug("getting proxy IP")
	var value string
	query := `SELECT value FROM config WHERE key = 'proxy_ip'`
	err := db.QueryRow(query).Scan(&value)
	if err != nil {
		logger.Fatal("Error getting proxy IP: ", err.Error())
	}
	return value
}

func GetEnableSpecificDomains() bool {
	logging.Debug("getting enable_specific_domains")
	var value string
	query := `SELECT value FROM config WHERE key = 'enable_specific_domains'`
	err := db.QueryRow(query).Scan(&value)
	if err != nil {
		logger.Fatal("Error getting enable_specific_domains: ", err.Error())
	}
	return value == "yes"
}

func GetEnableIPRestrictions() bool {
	logging.Debug("getting enable_ip_restrictions")
	var value string
	query := `SELECT value FROM config WHERE key = 'enable_ip_restrictions'`
	err := db.QueryRow(query).Scan(&value)
	if err != nil {
		logger.Fatal("Error getting enable_ip_restrictions: ", err.Error())
	}
	return value == "yes"
}
