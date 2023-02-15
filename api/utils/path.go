package utils

import (
	"os"

	"github.com/spf13/viper"
)

func GetRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return wd + "/" + viper.GetString("web.root"), err
}
