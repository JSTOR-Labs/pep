package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/JSTOR-Labs/pep/updated/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	manifestFile = "manifest.json"
)

// File is a file in the update package
type File struct {
	Src       string // Location in the update package
	Dst       string // Location to store on the machine
	Signature string
	Mode      os.FileMode
	User      int
	Group     int
}

type Versions struct {
	Latest    string
	Stable    string
	Manifests map[string]string // key is the version, value is the manifest filename
}

// Manifest represents the update manifest
type Manifest struct {
	Updated time.Time
	Version string
	Package string
	Files   map[string]File
}

// Manager provides access to the update repository
type Manager struct {
	sess       *session.Session
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
	bucket     string
}

// NewManager creates a new manager
func NewManager(bucket string) (*Manager, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return nil, err
	}

	downloader := s3manager.NewDownloader(sess)
	uploader := s3manager.NewUploader(sess)

	return &Manager{
		sess:       sess,
		downloader: downloader,
		uploader:   uploader,
		bucket:     bucket,
	}, nil
}

// Check pings the CDN to check for an update
func (m *Manager) Check() (bool, *Manifest, error) {
	buf := aws.NewWriteAtBuffer([]byte{})

	_, err := m.downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(manifestFile),
	})
	if err != nil {
		return false, nil, err
	}

	f, err := os.Open(fmt.Sprintf("/%s", manifestFile))
	if err != nil {
		return false, nil, err
	}
	dec := json.NewDecoder(f)
	localManifest := Manifest{}

	if err := dec.Decode(&localManifest); err != nil {
		return false, nil, err
	}

	remoteManifest := Manifest{}
	if err := json.Unmarshal(buf.Bytes(), &remoteManifest); err != nil {
		return false, nil, err
	}

	return localManifest.Updated.Before(remoteManifest.Updated), &remoteManifest, nil
}

// Download pulls the update specified by the provided manifest
func (m *Manager) Download(manifest *Manifest) error {
	f, err := os.Create(fmt.Sprintf("/mnt/%s.tar.gz", manifest.Updated.Format("2006-02-01_15-04-05")))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = m.downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(manifest.Package),
	})
	if err != nil {
		return err
	}

	return nil
}

// Install installs the update specifed by the manifest on the system
func (m *Manager) Install(manifest *Manifest) error {
	dst := "/"
	src := fmt.Sprintf("/mnt/%s.tar.gz", manifest.Updated.Format("2006-02-01_15-04-05"))

	f, err := os.Open(src)
	if err != nil {
		return err
	}

	if err := utils.Untar(dst, f); err != nil {
		return err
	}

	return os.Remove(src)
}
