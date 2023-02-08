package pdfs

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"time"

	_ "embed"
)

//go:embed encryption/cert.pem
var cert []byte

//go:embed encryption/ciphertext
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
