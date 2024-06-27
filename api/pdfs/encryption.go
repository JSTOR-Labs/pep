package pdfs

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/JSTOR-Labs/pep/api/utils"
	"github.com/manifoldco/promptui"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/scrypt"

	_ "embed"
)

type X509 struct {
	Certificate *bytes.Buffer
	PrivateKey  *bytes.Buffer
	PublicKey   *rsa.PublicKey
	Password    string
}

//go:embed keys/cert.pem
var cert []byte

//go:embed keys/ciphertext
var ciphertext []byte

//go:embed keys/key.pem
var privateKey []byte

func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	ct, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		return nil, err
	}
	return ct, err
}

func DecryptWithPrivateKey(ct []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ct, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, err
}

func GetKey(password string, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	key, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}

func PrepareAES(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return gcm, nil
}
func AESEncryptBytes(plaintext []byte, pw string) ([]byte, error) {

	key, salt, err := GetKey(pw, nil)
	if err != nil {
		return nil, err
	}

	gcm, err := PrepareAES(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	ciphertext = append(ciphertext, salt...)

	return ciphertext, nil
}

func AESDecryptBytes(ciphertext []byte, password string) ([]byte, error) {
	salt, data := ciphertext[len(ciphertext)-32:], ciphertext[:len(ciphertext)-32]

	key, _, err := GetKey(password, salt)
	if err != nil {
		return nil, err
	}

	// Creating block of algorithm
	gcm, err := PrepareAES(key)
	if err != nil {
		return []byte{}, err
	}

	// Deattached nonce and decrypt
	nonce := data[:gcm.NonceSize()]
	ct := data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return []byte{}, err
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

func GetPrivateKey(password string) (*rsa.PrivateKey, error) {
	key, err := AESDecryptBytes(privateKey, password)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(key)

	pk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pk, err
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

func PasswordPrompt() (*string, error) {
	var empty *string

	// validate the input
	validate := func(input string) error {
		if input == "" {
			return errors.New("A password is required. If you were not provided with a PDF Access password, please contact an administrator.")
		}
		return nil
	}

	// Each template displays the formatted data received.
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | blue }} ",
		Invalid: "{{ . | yellow }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     "Enter your PDF Access Password:",
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()

	if err != nil {
		return empty, err
	}

	return &result, err
}

func GetPDFPassword(pw []byte) (string, error) {
	var err error
	if pw == nil {
		exPath, err := utils.GetExecutablePath()
		if err != nil {
			return "", err
		}
		path := filepath.Join(exPath, "content", "password.txt")
		pw, err = os.ReadFile(path)

		if err != nil {
			return "", err
		}
	}

	pk, err := GetPrivateKey(string(pw))
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
	if err != nil {
		return "", err
	}
	return string(plaintext), err
}

func GenerateCert(pw string) (X509, error) {
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
		NotAfter:              time.Now().AddDate(5, 0, 0),
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

	read, err := io.ReadAll(caPrivKeyPEM)
	if err != nil {
		return X509{}, err
	}

	encryptedPK, err := AESEncryptBytes(read, pw)
	if err != nil {
		return X509{}, err
	}
	buf := bytes.NewBuffer(encryptedPK)

	return X509{
		Certificate: caPEM,
		PrivateKey:  buf,
		PublicKey:   &caPrivKey.PublicKey,
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
	readBuf, err := io.ReadAll(content)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, readBuf, os.ModePerm)

	return err
}

func SaveCert(keypath string, c X509) error {
	err := SaveFile(filepath.Join(keypath, "cert.pem"), c.Certificate)
	if err != nil {
		return err
	}

	err = SaveFile(filepath.Join(keypath, "key.pem"), c.PrivateKey)
	if err != nil {
		return err
	}
	return err
}

func SaveEncryptionFiles(password string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	keypath := filepath.Join(wd, "pdfs/keys/")

	c, err := GenerateCert(password)
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

	ct, err := EncryptWithPublicKey([]byte(userPW), c.PublicKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate encrypt user password")
		return err
	}
	err = os.WriteFile(filepath.Join(keypath, "ciphertext"), ct, 0644)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save ciphertext")
		return err
	}
	return nil
}

func HasPDFs() (bool, error) {
	hasPDFs := false
	path, err := utils.GetPDFPath()
	if err != nil {
		return hasPDFs, err
	}

	if _, err := os.Stat(path); err == nil {
		hasPDFs = true
	}
	return hasPDFs, err
}
func PromptUser(save bool) (string, error) {
	hasPW := false
	hasPDFs, err := HasPDFs()
	if err != nil {
		return "", err
	}

	exPath, err := utils.GetExecutablePath()
	if err != nil {
		return "", err
	}

	pwPath := filepath.Join(exPath, "content", "password.txt")
	if _, err := os.Stat(pwPath); err == nil {
		hasPW = true
	}

	var pw string
	if hasPDFs && !hasPW {
		password, err := PasswordPrompt()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get required PDF password")
			return "", err
		}
		if save {
			err = os.WriteFile(pwPath, []byte(*password), os.ModePerm)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to save PDF password")
				return "", err
			}
		}
		pw = *password

	}
	return pw, nil
}
func EncryptPDFDirectory(p string, pw string) error {
	// Extracting the password from a saved file is time consuming when repeated. Doing it once here
	// saves a lot of time.
	userPW, err := GetPDFPassword([]byte(pw))
	if err != nil {
		log.Error().Err(err).Msg("Failed to decrypt pdf user password")
		return err
	}

	err = filepath.Walk(p,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && info.Name() != ".DS_Store" {
				return EncryptPDF(path, userPW)
			}
			return err
		})
	return err
}

func EncryptPDF(path string, userPW string) error {
	// The owner password is generated and used here, then ignored. We will only ever decrypt with the user password.
	ownerPW, err := GetPassword()
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
	if err != nil {
		os.Remove(path)
		fmt.Println(path)
		fmt.Println("Error: ", err)
		err = nil
	}
	return err
}
