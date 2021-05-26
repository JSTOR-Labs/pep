package usb

import (
	"log"
	"os/exec"
	"strings"
)

type BlkIDData struct {
	UUID     string
	Type     string
	Label    string
	PartUUID string
}

func BlkID(path string) (*BlkIDData, error) {
	cmd := exec.Command("/sbin/blkid", path, "-o", "export")
	out, err := cmd.Output()
	if err != nil {
		log.Println(string(out))
		return nil, err
	}
	var b BlkIDData
	rawData := make(map[string]string)
	fields := strings.Split(string(out), "\n")
	for _, field := range fields {
		values := strings.Split(field, "=")
		if len(values) < 2 {
			continue
		}
		rawData[values[0]] = values[1]
	}
	setValue(rawData, "UUID", &b.UUID)
	setValue(rawData, "TYPE", &b.Type)
	setValue(rawData, "LABEL", &b.Label)
	setValue(rawData, "PARTUUID", &b.PartUUID)
	return &b, nil
}

func setValue(rawData map[string]string, key string, value *string) {
	if v, ok := rawData[key]; ok {
		*value = v
	}
}
