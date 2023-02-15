package pdfs

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"time"

	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-password/password"

	_ "embed"
)

type X509 struct {
	Certificate *bytes.Buffer
	PrivateKey  *bytes.Buffer
	PublicKey   *rsa.PublicKey
}

//go:embed keys/cert.pem
var cert []byte

//go:embed keys/ciphertext
var ciphertext []byte

func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, err
}

func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, err
}

func GetReadSeeker(file *os.File) *bytes.Reader {
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()

	buffer := make([]byte, size)

	// read file content to buffer
	file.Read(buffer)

	return bytes.NewReader(buffer)
}

func GetPrivateKey() (*rsa.PrivateKey, error) {
	key, err := os.ReadFile("./content/key.pem")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(key)

	pk, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pk.(*rsa.PrivateKey), err
}

func GetCert() (*x509.Certificate, error) {
	block, _ := pem.Decode(cert)
	c, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return c, err
}

func IsExpired(cert *x509.Certificate) bool {
	now := time.Now()
	return now.After(cert.NotAfter)
}
func IsEarly(cert *x509.Certificate) bool {
	now := time.Now()
	return now.Before(cert.NotBefore)
}

func GetPDFPassword() (string, error) {

	pk, err := GetPrivateKey()
	if err != nil {
		return "", err
	}

	cert, err := GetCert()
	if err != nil {
		return "", err
	}

	if IsExpired(cert) {
		return "", errors.New("certificate expired")
	}
	if IsEarly(cert) {
		return "", errors.New("certificate not yet valid")
	}

	plaintext, err := DecryptWithPrivateKey(ciphertext, pk)
	return string(plaintext), err
}

func GenerateCert() (X509, error) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Ithaka"},
			Country:       []string{"US"},
			Province:      []string{"MI"},
			Locality:      []string{"Ann Arbor"},
			StreetAddress: []string{"301 E Liberty St"},
			PostalCode:    []string{"48104"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 6, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return X509{}, err
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return X509{}, err
	}
	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	c, err := x509.ParseCertificate(caBytes)
	if err != nil {
		return X509{}, err
	}

	return X509{
		Certificate: caPEM,
		PrivateKey:  caPrivKeyPEM,
		PublicKey:   c.PublicKey.(*rsa.PublicKey),
	}, nil
}

func GetPassword() (string, error) {
	pw, err := password.Generate(64, 10, 10, false, true)
	if err != nil {
		return "", err
	}
	return pw, err
}

func SaveFile(path string, content *bytes.Buffer) error {
	readBuf, err := ioutil.ReadAll(content)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, readBuf, os.ModePerm)

	return err
}

func SaveCert(keypath string, cert X509) error {
	err := SaveFile(keypath+"cert.pem", cert.Certificate)
	if err != nil {
		return err
	}
	err = SaveFile(keypath+"key.pem", cert.PrivateKey)
	if err != nil {
		return err
	}
	return err
}

func SaveEncryptionFiles() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	keypath := wd + "/pdfs/keys/"

	c, err := GenerateCert()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate certificate")
		return err
	}
	SaveCert(keypath, c)

	userPW, err := GetPassword()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate PDF password")
		return err
	}

	ciphertext, err := EncryptWithPublicKey([]byte(userPW), c.PublicKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate encrypt user password")
		return err
	}
	err = SaveFile(keypath+"ciphertext", bytes.NewBuffer(ciphertext))
	if err != nil {
		log.Error().Err(err).Msg("Failed to save ciphertext")
		return err
	}
	return nil
}
func EncryptPDFDirectory(p string) error {
	err := filepath.Walk(p,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && info.Name() != ".DS_Store" {
				return EncryptPDF(path)
			}
			return err
		})
	return err
}

func EncryptPDF(path string) error {
	// The owner password is generated and used here, then ignored. We will only ever decrypt with the user password.
	ownerPW, err := GetPassword()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate PDF owner password")
		return err
	}

	userPW, err := GetPDFPassword()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate PDF owner password")
		return err
	}

	config := pdfcpu.LoadConfiguration()
	config.UserPW = userPW
	config.OwnerPW = ownerPW
	config.EncryptUsingAES = true
	config.EncryptKeyLength = 256

	err = pdfcpu.EncryptFile(path, path, config)
	return err
}
