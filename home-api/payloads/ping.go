package payloads

import (
	"time"
)

type VersionInfo struct {
	Version string    `json:"version"`
	Date    time.Time `json:"date"`
}

type Ping struct {
	MACAddress string      `json:"mac_addr"`
	Version    VersionInfo `json:"version_info"`
	HasErrors  bool        `json:"has_errors"`
}
