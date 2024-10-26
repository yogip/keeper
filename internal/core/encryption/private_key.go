package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"

	"os"
)

type PrivateKey struct {
	key *rsa.PrivateKey
}

func NewPrivateKey() (*PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("failed to generate rsa key: %w", err)
	}

	return &PrivateKey{key: key}, nil
}

func LoadPrivateKeyFromFile(file string) (*PrivateKey, error) {
	if file == "" {
		return nil, nil
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("reading private key from file (%s) error: %w", file, err)
	}

	return LoadPrivateKey(data)
}

func LoadPrivateKey(rawKey []byte) (*PrivateKey, error) {
	block, _ := pem.Decode(rawKey)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing private key error: %w", err)
	}

	return &PrivateKey{key: privateKey}, nil
}

func (p *PrivateKey) Dump() []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(p.key),
		},
	)
}

func (p *PrivateKey) Encrypt(plaintext []byte) ([]byte, error) {
	var encData []byte

	pubKey := p.PublicKey()
	if pubKey == nil {
		return nil, errors.New("private key is not intialized")
	}
	dataLen := len(plaintext)
	hash := sha512.New()
	step := chunkSize(pubKey.Size(), hash.Size())
	for begin := 0; begin < dataLen; begin += step {
		end := begin + step
		if end > dataLen {
			end = dataLen
		}

		encChunk, err := encryptChunk(plaintext[begin:end], hash, pubKey)
		if err != nil {
			return nil, fmt.Errorf("encrypt chunk error:%w", err)
		}

		encData = append(encData, encChunk...)
	}

	return encData, nil
}

func (p *PrivateKey) Decrypt(chipertext []byte) ([]byte, error) {
	var decData []byte

	dataLen := len(chipertext)
	hash := sha512.New()
	step := 512
	for begin := 0; begin < dataLen; begin += step {
		end := begin + step
		if end > dataLen {
			end = dataLen
		}

		decChunk, err := decryptChunk(chipertext[begin:end], hash, p.key)
		if err != nil {
			return nil, fmt.Errorf("decrypt chunk error:%w", err)
		}

		decData = append(decData, decChunk...)
	}

	return decData, nil
}

func (p *PrivateKey) PublicKey() *rsa.PublicKey {
	if p.key == nil {
		return nil
	}
	return &p.key.PublicKey
}

func encryptChunk(data []byte, hash hash.Hash, pupKey *rsa.PublicKey) ([]byte, error) {
	b, err := rsa.EncryptOAEP(hash, rand.Reader, pupKey, data, nil)
	if err != nil {
		return nil, fmt.Errorf("rsa OAEP encrypt error:%w", err)
	}

	return b, nil
}

func decryptChunk(data []byte, hash hash.Hash, privKey *rsa.PrivateKey) ([]byte, error) {
	b, err := rsa.DecryptOAEP(hash, rand.Reader, privKey, data, nil)
	if err != nil {
		return nil, fmt.Errorf("rsa OAEP decrypt error:%w", err)
	}

	return b, nil
}

// The message must be no longer than the length of the public modulus minus
// twice the hash length, minus a further 2.
func chunkSize(keySize int, hashSize int) int {
	// https://cs.opensource.google/go/go/+/refs/tags/go1.23.1:src/crypto/rsa/rsa.go;l=527
	return keySize - 2*hashSize - 2
}
