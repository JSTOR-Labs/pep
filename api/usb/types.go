package usb

import "time"

type USBDrive struct {
	Name       string      `json:"name"`
	Size       string      `json:"size"`
	Model      string      `json:"model"`
	Partitions []Partition `json:"partitions"`
}

type Partition struct {
	Name       string             `json:"name"`
	Label      string             `json:"label"`
	Size       string             `json:"size"`
	SizeKB     int32              `json:"sizeKB"`
	Filesystem string             `json:"filesystem"`
	Config     FlashDriveManifest `json:"config,omitempty"`
}

type FlashDriveManifest struct {
	Indices []string  `json:"indices"`
	Created time.Time `json:"created"`
	Name    string    `json:"name"`
}

type driveInfo struct {
	APIPath  string    `msgpack:"api"`
	JavaHome string    `msgpack:"java"`
	ESRoot   string    `msgpack:"es"`
	Created  time.Time `msgpack:"created"`
	Hostname string    `msgpack:"source"`
}

type elasticConfig struct {
	Port                 int      `yaml:"http.port"`
	Repo                 []string `yaml:"path.repo"`
	XPackMachineLearning bool     `yaml:"xpack.ml.enabled"`
	DiscoveryType        string   `yaml:"discovery.type"`
}
