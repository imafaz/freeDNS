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
	query = `CREATE TABLE IF NOT EXISTS domains (
        domain TEXT NOT NULL PRIMARY KEY
    );`
	_, err = db.Exec(query)
	if err != nil {
		logger.Fatal("Error creating domains table: ", err.Error())
	}

	query = `CREATE TABLE IF NOT EXISTS config (
    key TEXT NOT NULL PRIMARY KEY,
    value TEXT NOT NULL
);`
	_, err = db.Exec(query)
	if err != nil {
		logger.Fatal("Error creating config table: ", err.Error())
	}

	query = `CREATE TABLE IF NOT EXISTS allowIP (
		ip TEXT NOT NULL PRIMARY KEY
		);`
	_, err = db.Exec(query)
	if err != nil {
		logger.Fatal("Error creating allowIP table: ", err.Error())
	}
}

func insertConfig() {
	logging.Debug("inserting configs")

	var query string
	var err error

	configs := map[string]string{
		"ip_restrictions":  "yes",
		"specific_domains": "yes",
		"proxy_ip":         config.GetServerIP(),
		"server":           "0.0.0.0",
		"port":             "53",
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

func UpdateConfig(key string, value string) {
	logging.Debugf("updating config %s to %s", key, value)
	query := `UPDATE config SET value = ? WHERE key = ?`
	_, err := db.Exec(query, value, key)
	if err != nil {
		logger.Fatal("Error updating config key: ", err.Error())
	}
}
func AddDomain(domain string) {
	logging.Debugf("adding domain %s to specific domains", domain)
	query := `INSERT INTO domains (domain) VALUES (?)`
	_, err := db.Exec(query, domain)
	if err != nil {
		logger.Fatal("Error adding domain: ", err.Error())
	}
}
func AllowIP(IP string) {
	logging.Debugf("allowing IP %s", IP)
	query := `INSERT INTO allowIP (ip) VALUES (?)`
	_, err := db.Exec(query, IP)
	if err != nil {
		logger.Fatal("Error allowing IP: ", err.Error())
	}
}

func RemoveDomain(domain string) {
	logging.Debugf("removing domain %s", domain)
	query := `DELETE FROM domains WHERE domain = ?`
	_, err := db.Exec(query, domain)
	if err != nil {
		logger.Fatal("Error removing domain: ", err.Error())
	}
}

func RemoveIP(IP string) {
	logging.Debugf("removing IP %s", IP)
	query := `DELETE FROM allowIP WHERE ip = ?`
	_, err := db.Exec(query, IP)
	if err != nil {
		logger.Fatal("Error removing IP: ", err.Error())
	}
}

func DomainExists(domain string) bool {
	domain = strings.TrimSuffix(domain, ".")
	logging.Debugf("checking if domain %s exists", domain)
	query := `SELECT EXISTS(SELECT 1 FROM domains WHERE domain = ?)`
	var exists bool
	err := db.QueryRow(query, domain).Scan(&exists)
	if err != nil {
		logger.Fatal("Error checking if domain exists: ", err.Error())
	}
	return exists
}

func IPExists(IP string) bool {
	logging.Debugf("checking if IP %s exists", IP)
	query := `SELECT EXISTS(SELECT 1 FROM allowIP WHERE ip = ?)`
	var exists bool
	err := db.QueryRow(query, IP).Scan(&exists)
	if err != nil {
		logger.Fatal("Error checking if IP exists: ", err.Error())
	}
	return exists
}

func GetAllConfig() (map[string]string, error) {
	logging.Debug("getting all configurations")
	configs := make(map[string]string)
	query := `SELECT key, value FROM config`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		configs[key] = value
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}

func GetDomains() []string {
	logging.Debug("getting all domains")
	var domains []string
	query := `SELECT domain FROM domains`
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
	query := `SELECT ip FROM allowIP`
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

func GetConfig(key string) string {
	logging.Debugf("getting config for key %s", key)
	var value string
	query := `SELECT value FROM config WHERE key = ?`
	err := db.QueryRow(query, key).Scan(&value)
	if err != nil {
		logger.Fatal("Error getting config value: ", err.Error())
	}
	return value
}
