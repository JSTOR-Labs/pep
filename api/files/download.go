package files

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	CachePath     = "/mnt/cache/"
	ElasticURLFmt = "https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-%s-no-jdk-windows-x86_64.zip"
	JavaURL       = "https://github.com/adoptium/temurin11-binaries/releases/download/jdk-11.0.13%2B8/OpenJDK11U-jre_x64_windows_hotspot_11.0.13_8.zip"
)

func GetElasticURL(version string) string {
	return fmt.Sprintf(ElasticURLFmt, version)
}

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

	dialer := &net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: 5 * time.Second,
				}
				return d.DialContext(ctx, "udp", "8.8.8.8:53")
			},
		},
	}
	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, addr)
	}

	http.DefaultTransport.(*http.Transport).DialContext = dialContext

	client := http.Client{Timeout: time.Minute * 5}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create request")
		return "", err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to download file")
		return "", err
	}
	defer res.Body.Close()

	out, err := os.Create(filepath.Join(CachePath, filename))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create file")
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	return filepath.Join(CachePath, filename), err
}

func Untar(dst string, r io.Reader) ([]string, error) {
	var filenames []string
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return filenames, err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return filenames, nil

		// return any other error
		case err != nil:
			return filenames, err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		filenames = append(filenames, header.Name)

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return filenames, err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return filenames, err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return filenames, err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
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
