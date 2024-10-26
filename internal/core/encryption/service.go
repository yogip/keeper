package encryption

import "github.com/pkg/errors"

// EncryptionService is a service to encrypt and decrypt data by keys according to the keys version.
type EncryptionService struct {
	managers map[int64]*EncryptionManager
	dir      string
}

func NewEncryptionService(keysDir string) *EncryptionService {
	return &EncryptionService{
		dir:      keysDir,
		managers: make(map[int64]*EncryptionManager),
	}
}

// Encrypt plaintext by encryption key version.
func (s *EncryptionService) Encrypt(plaintext []byte, version int64) ([]byte, DataKey, error) {
	manager, err := s.getManager(version)
	if err != nil {
		return nil, nil, err
	}

	return manager.Encrypt(plaintext)
}

// Decrypt dycrypts ciphertext using data key.
func (s *EncryptionService) Decrypt(chipertext []byte, dataKey DataKey, version int64) ([]byte, error) {
	manager, err := s.getManager(version)
	if err != nil {
		return nil, err
	}

	return manager.Decrypt(chipertext, dataKey)
}

// getManager gets encryption manager by key version.
func (s *EncryptionService) getManager(version int64) (*EncryptionManager, error) {
	manager, ok := s.managers[version]
	if ok {
		return manager, nil
	}

	manager, err := NewEncryptionManager(version, s.dir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create encryption manager")
	}
	s.managers[version] = manager
	return manager, nil
}
