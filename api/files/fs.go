package files

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rs/zerolog/log"
)

var (
	DownloadPath           = "./downloads"
	Bucket                 = ""
	Prefix                 = ""
	WindowsDir             = "./JSTOR-Windows"
	ChromebookDir          = "./JSTOR-Chromebook"
	MacDir                 = "./JSTOR-Mac"
	JSTORDir               = "JSTOR"
	ContentDir             = "content"
	PDFDir                 = "pdfs"
	InstallGuideDir        = "./install_guides"
	ChromebookInstallGuide = filepath.Join(InstallGuideDir, "chromebook/README.pdf")
	WindowsInstallGuide    = filepath.Join(InstallGuideDir, "windows/README.pdf")
	ShellDir               = "./shell"
	WindowsShell           = filepath.Join(ShellDir, "start.bat")
	ChromebookShell        = filepath.Join(ShellDir, "start.sh")
	MacShell               = filepath.Join(ShellDir, "start.command")
	FrontendDir            = "dist"
)

type S3Downloader struct {
	DLManager *s3manager.Downloader
	Bucket    string
	Dir       string
}

func CreateS3Session() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create session")
		return nil, err
	}
	return sess, nil
}
func DownloadBucket(bucket string, prefix string, dir string) error {
	sess, err := CreateS3Session()
	if err != nil {
		return err
	}

	manager := s3manager.NewDownloader(sess)
	d := S3Downloader{Bucket: bucket, Dir: dir, DLManager: manager}

	sess, err = CreateS3Session()
	if err != nil {
		return err
	}

	client := s3.New(sess)
	params := &s3.ListObjectsInput{Bucket: &bucket, Prefix: &prefix}
	client.ListObjectsPages(params, d.EachPage)
	return nil
}

func (d *S3Downloader) EachPage(page *s3.ListObjectsOutput, more bool) bool {
	for _, obj := range page.Contents {
		d.DownloadToFile(*obj.Key)
	}

	return true
}

func (d *S3Downloader) DownloadToFile(key string) error {
	file := filepath.Join(d.Dir, key)
	path := filepath.Dir(file)
	if strings.HasSuffix(key, "/") {
		path = file
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	// Set up the local file
	fd, err := os.Create(file)
	if err != nil {
		log.Error().Err(err).Msg("failed to create file")
		return err
	}
	defer fd.Close()

	// Download the file using the AWS SDK for Go
	fmt.Printf("Downloading s3://%s/%s to %s...\n", d.Bucket, key, file)
	params := &s3.GetObjectInput{Bucket: &d.Bucket, Key: &key}
	_, err = d.DLManager.Download(fd, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to download file")
	}
	return err
}

func CopyFile(src, dst string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	return err
}

func CopyRecursive(src, dst string) error {
	if s, err := os.Stat(src); err == nil && s.IsDir() {
		_ = os.Mkdir(fmt.Sprintf("%s/%s", dst, s.Name()), s.Mode())
		contents, err := ioutil.ReadDir(src)
		if err != nil {
			return err
		}

		for _, f := range contents {
			err = CopyRecursive(fmt.Sprintf("%s/%s", src, f.Name()), fmt.Sprintf("%s/%s", dst, s.Name()))
			if err != nil {
				return err
			}
		}
	} else if err == nil {
		err = CopyFile(src, fmt.Sprintf("%s/%s", dst, s.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
