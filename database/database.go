package database

import (
	"database/sql"

	"github.com/imafaz/logger"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB(dbPath string) {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Fatal(err.Error())
	}
	createTables()
	insertConfig()
}

func createTables() {
	var err error
	var sql string
	sql = `CREATE TABLE IF NOT EXISTS domains (
        domain TEXT NOT NULL PRIMARY KEY
    );`
	_, err = db.Exec(sql)
	if err != nil {
		logger.Fatal(err.Error())
	}

	sql = `CREATE TABLE IF NOT EXISTS config (
    key TEXT NOT NULL PRIMARY KEY,
    value TEXT NOT NULL
);`
	_, err = db.Exec(sql)
	if err != nil {
		logger.Fatal(err.Error())
	}

	sql = `CREATE TABLE IF NOT EXISTS allowIP (
		ip TEXT NOT NULL PRIMARY KEY
		);`
	_, err = db.Exec(sql)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func insertConfig() {
	var sql string
	var err error

	configs := map[string]string{
		"ip_restrictions":     "no",
		"domain_restrictions": "no",
		"revers_proxy_ip":     "",
		"server":              "0.0.0.0",
		"port":                "53",
	}

	for key, value := range configs {
		sql = `SELECT COUNT(*) FROM config WHERE key = ?`
		var count int
		err = db.QueryRow(sql, key).Scan(&count)
		if err != nil {
			logger.Fatal(err.Error())
		}
		if count == 0 {
			sql = `INSERT INTO config (key, value) VALUES (?, ?)`
			_, err = db.Exec(sql, key, value)
			if err != nil {
				logger.Fatal(err.Error())
			}
		}
	}
}

func UpdateConfig(key string, value string) {
	sql := `UPDATE config SET value = ? WHERE key = ?`
	_, err := db.Exec(sql, value, key)
	logger.Fatal(err.Error())
}
func AddDomain(domain string) {
	sql := `INSERT INTO domains (domain) VALUES (?, ?)`
	_, err := db.Exec(sql, domain)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func AllowIP(IP string) {
	sql := `INSERT INTO allowIP (ip) VALUES (?, ?)`
	_, err := db.Exec(sql, IP)
	if err != nil {
		logger.Fatal(err.Error())
	}
}
func RemoveDomain(domain string) {
	sql := `DELETE FROM domains WHERE domain = ?`
	_, err := db.Exec(sql, domain)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func RemoveIP(IP string) {
	sql := `DELETE FROM allowIP WHERE ip = ?`
	_, err := db.Exec(sql, IP)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func DomainExists(domain string) bool {
	var exists bool
	sql := `SELECT EXISTS(SELECT 1 FROM domains WHERE domain = ?)`
	err := db.QueryRow(sql, domain).Scan(&exists)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return exists
}

func IPExists(IP string) bool {
	var exists bool
	sql := `SELECT EXISTS(SELECT 1 FROM allowIP WHERE ip = ?)`
	err := db.QueryRow(sql, IP).Scan(&exists)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return exists
}
func GetAllConfig() (map[string]string, error) {
	configs := make(map[string]string)
	sql := `SELECT key, value FROM config`
	rows, err := db.Query(sql)
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

func GetDomains() ([]string, error) {
	var domains []string
	sql := `SELECT domain FROM domains`
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var domain string
		if err := rows.Scan(&domain); err != nil {
			return nil, err
		}
		domains = append(domains, domain)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return domains, nil
}

func GetIPs() ([]string, error) {
	var ips []string
	sql := `SELECT ip FROM allowIP`
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ip string
		if err := rows.Scan(&ip); err != nil {
			return nil, err
		}
		ips = append(ips, ip)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ips, nil
}
func GetConfig(key string) string {
	var value string
	sql := `SELECT value FROM config WHERE key = ?`
	err := db.QueryRow(sql, key).Scan(&value)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return value
}
