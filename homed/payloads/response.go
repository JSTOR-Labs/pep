package payloads

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/JSTOR-Labs/pep/homed/constants"
)

type ActionOp int

const (
	// Reboot the NUC
	OpReboot ActionOp = iota
	// Fetch a file from an internet source
	OpFetch
	// Verify the checksum of a file
	OpVerify
	// Extract a tar archive
	OpExtract
	// Delete files or directories from disk
	OpDelete
	// Execute a command on the NUC
	OpExec
	// Put a file to an internet destination
	OpPut
	// Perform an Update
	OpUpdate
)

type ResponseAction struct {
	Op ActionOp `json:"op"`
	// This is why we don't just submodule the payloads package
	Data json.RawMessage `json:"data,omitempty"`
}

type PingResponse struct {
	Actions []ResponseAction `json:"actions,omitempty"`
}

type FetchData struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

type VerifyData struct {
	Path     string `json:"path"`
	Checksum []byte `json:"checksum"`
}

type ExtractData struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type DeleteData struct {
	Path string `json:"path"`
}

type ExecData struct {
	Path string   `json:"path"`
	Args []string `json:"args"`
}

type PutData struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

type UpdateData struct {
	Version Version `json:"version"`
}

type Version struct {
	ID      uint      `gorm:"primaryKey"`
	Version string    `json:"version"`
	Date    time.Time `json:"date"`
	Assets  []Asset   `json:"assets"`
}

type Asset struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	Size     uint   `json:"size"`
	Dst      string `json:"dst"`

	CreatedAt time.Time `json:"created_at"`
}

func (a *Asset) Process(c *http.Client) error {
	resp, err := c.Get(fmt.Sprintf(constants.GetAssetEp, a.ID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch a.FileType {
	case "tar":
		body, err := gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}

		tarReader := tar.NewReader(body)
		for {
			header, err := tarReader.Next()

			if err == io.EOF {
				break
			}

			if err != nil {
				return err
			}

			dstPath := filepath.Join(a.Dst, header.Name)

			switch header.Typeflag {
			case tar.TypeDir:
				if err := os.Mkdir(dstPath, fs.FileMode(header.Mode)); err != nil {
					return err
				}
			case tar.TypeReg:
				outFile, err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fs.FileMode(header.Mode))
				if err != nil {
					return err
				}

				if _, err := io.Copy(outFile, tarReader); err != nil {
					return err
				}

				if err := outFile.Close(); err != nil {
					return err
				}
			}
		}
	default:
		f, err := os.Create(a.Dst)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, resp.Body)
		if err != nil {
			return err
		}
	}

	return nil
}
