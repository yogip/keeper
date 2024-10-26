package encryption

import (
	"fmt"
)

// Raw data key encrypted by EncryptionKey.
type DataKey []byte

// EncryptionManager is a manager for encrypting data with uniq data key.
// Plaintext data encrypts by the data key and the data key encrypts by the encryption key.
type EncryptionManager struct {
	version int64
	key     PrivateKey
}

// NewEncryptionManager loads encryption key from file and creates new encryption manager.
func NewEncryptionManager(version int64, keyDir string) (*EncryptionManager, error) {
	keyPath := fmt.Sprintf("%s/encription_key_v_%d.key", keyDir, version)
	privateKey, err := LoadPrivateKeyFromFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("encryption manager creating error: %w", err)
	}
	manager := EncryptionManager{version: version, key: *privateKey}
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
func (e *EncryptionManager) Decrypt(chipertext []byte, rawKey DataKey) ([]byte, error) {

	dataKey, err := LoadPrivateKey(rawKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load data key: %w", err)
	}

	plaintext, err := dataKey.Decrypt(chipertext)
	if err != nil {
		return nil, fmt.Errorf("failed to Decrypt chipertext by data key: %w", err)
	}

	return plaintext, nil
}
