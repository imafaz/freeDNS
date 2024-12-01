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
	var query string
	query = `CREATE TABLE IF NOT EXISTS domains (
        domain TEXT NOT NULL PRIMARY KEY
    );`
	_, err = db.Exec(query)
	if err != nil {
		logger.Fatal(err.Error())
	}

	query = `CREATE TABLE IF NOT EXISTS config (
    key TEXT NOT NULL PRIMARY KEY,
    value TEXT NOT NULL
);`
	_, err = db.Exec(query)
	if err != nil {
		logger.Fatal(err.Error())
	}

	query = `CREATE TABLE IF NOT EXISTS allowIP (
		ip TEXT NOT NULL PRIMARY KEY
		);`
	_, err = db.Exec(query)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func insertConfig() {
	var query string
	var err error

	configs := map[string]string{
		"ip_restrictions":     "yes",
		"domain_restrictions": "yes",
		"revers_proxy_ip":     "127.0.0.1",
		"server":              "0.0.0.0",
		"port":                "53",
	}

	for key, value := range configs {
		query = `SELECT COUNT(*) FROM config WHERE key = ?`
		var count int
		err = db.QueryRow(query, key).Scan(&count)
		if err != nil {
			logger.Fatal(err.Error())
		}
		if count == 0 {
			query = `INSERT INTO config (key, value) VALUES (?, ?)`
			_, err = db.Exec(query, key, value)
			if err != nil {
				logger.Fatal(err.Error())
			}
		}
	}
}

func UpdateConfig(key string, value string) {
	query := `UPDATE config SET value = ? WHERE key = ?`
	_, err := db.Exec(query, value, key)
	if err != nil {
		logger.Fatal(err.Error())
	}
}
func AddDomain(domain string) {
	query := `INSERT INTO domains (domain) VALUES (?)`
	_, err := db.Exec(query, domain)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func AllowIP(IP string) {
	query := `INSERT INTO allowIP (ip) VALUES (?)`
	_, err := db.Exec(query, IP)
	if err != nil {
		logger.Fatal(err.Error())
	}
}
func RemoveDomain(domain string) {
	query := `DELETE FROM domains WHERE domain = ?`
	_, err := db.Exec(query, domain)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func RemoveIP(IP string) {
	query := `DELETE FROM allowIP WHERE ip = ?`
	_, err := db.Exec(query, IP)
	if err != nil {
		logger.Fatal(err.Error())
	}
}
func DomainExists(domain string) bool {
	query := `SELECT domain FROM domains WHERE domain = ?`
	err := db.QueryRow(query, domain).Scan(&domain)
	if err != nil {

		if err != sql.ErrNoRows {
			logger.Fatal(err.Error())
		}

		return false
	}

	return true
}

// func DomainExists(domain string) bool {
// 	var exists bool
// 	query := `SELECT EXISTS(SELECT 1 FROM domains WHERE domain = ?)`
// 	err := db.QueryRow(query, domain).Scan(&exists)
// 	if err != nil {
// 		logger.Fatal(err.Error())
// 	}
// 	return exists
// }

func IPExists(IP string) bool {
	query := `SELECT ip FROM allowIP WHERE ip = ?`
	err := db.QueryRow(query, IP).Scan(&IP)
	if err != nil {

		if err != sql.ErrNoRows {
			logger.Fatal(err.Error())
		}

		return false
	}

	return true
}

//	func IPExists(IP string) bool {
//		var exists bool
//		query := `SELECT EXISTS(SELECT 1 FROM allowIP WHERE ip = ?)`
//		err := db.QueryRow(query, IP).Scan(&exists)
//		if err != nil {
//			logger.Fatal(err.Error())
//		}
//		return exists
//	}
func GetAllConfig() (map[string]string, error) {
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
	var domains []string
	query := `SELECT domain FROM domains`
	rows, err := db.Query(query)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var domain string
		if err := rows.Scan(&domain); err != nil {
			logger.Fatal(err.Error())
		}
		domains = append(domains, domain)
	}

	if err := rows.Err(); err != nil {
		logger.Fatal(err.Error())
	}

	return domains
}

func GetIPs() []string {
	var ips []string
	query := `SELECT ip FROM allowIP`
	rows, err := db.Query(query)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var ip string
		if err := rows.Scan(&ip); err != nil {
			logger.Fatal(err.Error())
		}
		ips = append(ips, ip)
	}

	if err := rows.Err(); err != nil {
		logger.Fatal(err.Error())
	}

	return ips
}
func GetConfig(key string) string {
	var value string
	query := `SELECT value FROM config WHERE key = ?`
	err := db.QueryRow(query, key).Scan(&value)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return value
}
