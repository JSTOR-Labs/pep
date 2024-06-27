package utils

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func GetRoot() (string, error) {
	exPath, err := GetExecutablePath()
	if err != nil {
		return "", err
	}
	path := filepath.Join(exPath, viper.GetString("web.root"))
	return path, err
}

func GetExecutablePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}
func GetPDFPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(ex)
	return filepath.Join(exPath, "pdfs"), nil
}
func GetManifestPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(ex)
	return filepath.Join(exPath, "manifest.json"), nil
}
