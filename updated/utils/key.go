package utils

import (
	"crypto/ed25519"
	"encoding/gob"
	"io"
)

// SaveKey saves a private or public key
func SaveKey(w io.Writer, key interface{}) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(key)
}

// LoadPrivKey loads an ed25519 private key
func LoadPrivKey(r io.Reader) (ed25519.PrivateKey, error) {
	var key ed25519.PrivateKey
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&key)
	return key, err
}

// LoadPubKey loads an ed25519 public key
func LoadPubKey(r io.Reader) (ed25519.PublicKey, error) {
	var key ed25519.PublicKey
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&key)
	return key, err
}
