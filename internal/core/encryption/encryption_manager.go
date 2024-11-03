package encryption

import (
	"fmt"
	"os"
)

// Raw data key encrypted by EncryptionKey.
type DataKey []byte

// EncryptionManager is a manager for encrypting data with uniq data key.
// Plaintext data encrypts by the data key and the data key encrypts by the encryption key.
type EncryptionManager struct {
	version int64
	key     *PrivateKey
}

func readPrivateKey(file string) ([]byte, error) {
	if file == "" {
		return nil, nil
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("reading private key from file (%s) error: %w", file, err)
	}

	return data, nil
}

// NewEncryptionManager loads encryption key from file and creates new encryption manager.
func NewEncryptionManager(version int64, keyDir string, decoder func([]byte) ([]byte, error)) (*EncryptionManager, error) {
	keyPath := fmt.Sprintf("%s/encription_key_v_%d.pem", keyDir, version)
	encPrivateKey, err := readPrivateKey(keyPath)
	if err != nil {
		return nil, fmt.Errorf("encryption manager creating error: %w", err)
	}

	rawPrivateKey, err := decoder(encPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("decoding private key (%s) error: %w", keyPath, err)
	}

	privateKey, err := LoadPrivateKey(rawPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("parsing private key (%s) error: %w", keyPath, err)
	}

	manager := EncryptionManager{version: version, key: privateKey}
	return &manager, nil
}

// Encrypt creates new data key and encrypts data using it.
func (e *EncryptionManager) Encrypt(data []byte) ([]byte, DataKey, error) {
	dataKey, err := NewPrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate data key: %w", err)
	}

	encryptedData, err := dataKey.Encrypt(data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	rawDataKey := dataKey.Dump()
	encryptedDataKey, err := e.key.Encrypt(rawDataKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt data key: %w", err)
	}

	return encryptedData, encryptedDataKey, nil
}

// Decrypt load given data key and decrypts data using it.
func (e *EncryptionManager) Decrypt(chipertext []byte, encDataKey DataKey) ([]byte, error) {
	rawDK, err := e.key.Decrypt(encDataKey)
	if err != nil {
		return nil, fmt.Errorf("failed to Decrypt data key: %w", err)
	}

	dataKey, err := LoadPrivateKey(rawDK)
	if err != nil {
		return nil, fmt.Errorf("failed to load data key: %w", err)
	}

	plaintext, err := dataKey.Decrypt(chipertext)
	if err != nil {
		return nil, fmt.Errorf("failed to Decrypt chipertext by data key: %w", err)
	}

	return plaintext, nil
}
