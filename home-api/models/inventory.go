package models

import (
	"time"
)

type Device struct {
	ID          uint   `gorm:"primaryKey"`
	MACAddress  string `gorm:"uniqueIndex"`
	Institution string
	Shipping    ShippingInfo `gorm:"embedded;embeddedPrefix:shipping_"`
	LastSeen    time.Time
	Version     VersionInfo `gorm:"embedded;embeddedPrefix:version_"`
}

type ShippingInfo struct {
	Address        string
	City           string
	State          string
	Zip            string
	TrackingNumber string
}

type VersionInfo struct {
	Version string
	Date    time.Time
}
