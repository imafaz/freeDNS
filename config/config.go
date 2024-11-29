package config

import (
	_ "embed"
	"fmt"
	"strings"
)

var (
	Debug bool = false
	//go:embed version
	version string

	//go:embed name
	name string
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
