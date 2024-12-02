package config

import (
	_ "embed"
	"fmt"
	"net"
	"strings"

	"github.com/imafaz/logger"
)

var (
	//go:embed version
	version string

	//go:embed name
	name string

	Debug bool
)

func GetVersion() string {
	return strings.TrimSpace(version)
}

func GetName() string {
	return strings.TrimSpace(name)
}
func GetDBPath() string {
	return fmt.Sprintf("/etc/%s/%s.db", GetName(), GetName())
}
func GetLogPath() string {
	return fmt.Sprintf("/var/log/%s.log", GetName())
}
func GetServerIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.Fatal(err.Error())
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}
