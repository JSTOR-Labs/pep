package files

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	CachePath  = "cache/"
	ElasticURL = "https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-7.10.2-no-jdk-windows-x86_64.zip"
	JavaURL    = "https://github.com/AdoptOpenJDK/openjdk11-binaries/releases/download/jdk-11.0.10%2B9/OpenJDK11U-jre_x64_windows_hotspot_11.0.10_9.zip"
)

func checkCache(filename string) (string, bool) {
	if _, err := os.Stat(CachePath); err != nil {
		_ = os.Mkdir(CachePath, 0777)
	}
	if _, err := os.Stat(CachePath + filename); err == nil {
		// file exists
		return CachePath + filename, true
	}

	return "", false
}

func DownloadFile(url, filename string) (string, error) {
	filePath, ok := checkCache(filename)
	if ok {
		return filePath, nil
	}

	client := http.Client{Timeout: time.Minute * 5}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	out, err := os.Create(CachePath + filename)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	return CachePath + filename, err
}

// Unzip will take a zip file path as src, and unzip it to the dest
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, f.Name)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		os.Chown(fpath, 1000, 1000)

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
