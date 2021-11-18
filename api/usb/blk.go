package usb

import (
	"os/exec"
	"strings"

	"github.com/JSTOR-Labs/pep/api/which"
	"github.com/rs/zerolog/log"
)

type BlkIDData struct {
	UUID     string
	Type     string
	Label    string
	PartUUID string
}

func BlkID(path string) (*BlkIDData, error) {
	cmd := exec.Command(which.LookupExecutable("blkid"), path)
	out, err := cmd.Output()
	if err != nil {
		log.Warn().Err(err).Str("out", string(out)).Msg("blkid failed")
		return nil, err
	}
	var b BlkIDData
	rawData := make(map[string]string)
	fields := strings.Split(string(out), " ")
	for _, field := range fields[0:] {
		values := strings.Split(field, "=")
		if len(values) < 2 {
			continue
		}
		rawData[values[0]] = strings.Trim(values[1], "\"")
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
