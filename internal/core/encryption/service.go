package encryption

import (
	"encoding/hex"
	"keeper/internal/core/model"

	"github.com/pkg/errors"
)

// EncryptionService is a service to encrypt and decrypt data by keys according to the keys version.
type EncryptionService struct {
	masterKey *PrivateKey
	managers  map[int64]*EncryptionManager
	dir       string
}

func NewEncryptionService(keysDir string, masterKey *PrivateKey) *EncryptionService {
	return &EncryptionService{
		dir:       keysDir,
		masterKey: masterKey,
		managers:  make(map[int64]*EncryptionManager),
	}
}

// Encrypt plaintext by encryption key version.
func (s *EncryptionService) Encrypt(plaintext []byte, version int64) (string, string, error) {
	manager, err := s.getManager(version)
	if err != nil {
		return "", "", err
	}

	encData, key, err := manager.Encrypt(plaintext)
	if err != nil {
		return "", "", err
	}
	return hex.EncodeToString(encData), hex.EncodeToString(key), nil
}

// Decrypt dycrypts ciphertext using data key.
func (s *EncryptionService) Decrypt(chipertext string, dataKey *model.DataKey) ([]byte, error) {
	manager, err := s.getManager(dataKey.Version)
	if err != nil {
		return nil, err
	}
	ch, err := hex.DecodeString(chipertext)
	if err != nil {
		return nil, err
	}
	key, err := hex.DecodeString(dataKey.Key)
	if err != nil {
		return nil, err
	}

	return manager.Decrypt(ch, DataKeyRaw(key))
}

// getManager gets encryption manager by key version.
func (s *EncryptionService) getManager(version int64) (*EncryptionManager, error) {
	manager, ok := s.managers[version]
	if ok {
		return manager, nil
	}

	manager, err := NewEncryptionManager(version, s.dir, s.masterKey.Decrypt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create encryption manager")
	}
	s.managers[version] = manager
	return manager, nil
}
