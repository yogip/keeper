package encryption

import (
	"errors"
	"keeper/internal/core/model"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptionService(t *testing.T) {
	privKeyData := `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`
	masterKey, err := NewPrivateKey()
	require.NoError(t, err)

	encPrivateKey, err := masterKey.Encrypt([]byte(privKeyData))
	require.NoError(t, err)

	type args struct {
		plainText     []byte
		encPrivateKey []byte
		keyDir        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Test success",
			args: args{
				plainText:     []byte("plain text"),
				encPrivateKey: encPrivateKey,
				keyDir:        "/tmp",
			},
			wantErr: nil,
		},
		{
			name: "Test failed: encryption key in wrong format",
			args: args{
				plainText:     []byte("plain text"),
				encPrivateKey: []byte(privKeyData), // not encrypted
				keyDir:        "/tmp",
			},
			wantErr: errors.New("failed to create encryption manager: decoding private key (/tmp/encryption_key_v_10.pem) error: decrypt chunk error:rsa OAEP decrypt error:crypto/rsa: decryption error"),
		},
		{
			name: "Test failed: There is no key file",
			args: args{
				plainText:     []byte("plain text"),
				encPrivateKey: []byte(privKeyData),
				keyDir:        "/not_exists",
			},
			wantErr: errors.New("failed to create encryption manager: encryption manager creating error: reading private key from file (/not_exists/encryption_key_v_10.pem) error: open /not_exists/encryption_key_v_10.pem: no such file or directory"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyPath := "/tmp/encryption_key_v_10.pem"
			f, err := os.OpenFile(keyPath, os.O_CREATE|os.O_WRONLY, 0777)
			require.NoError(t, err)
			defer os.Remove(keyPath)

			f.Write(tt.args.encPrivateKey)
			f.Close()

			encService := NewEncryptionService(tt.args.keyDir, masterKey)

			chiphertext, dataKey, err := encService.Encrypt(tt.args.plainText, 10)
			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
				return
			} else {
				require.NoError(t, err)
			}

			got, err := encService.Decrypt(chiphertext, &model.DataKey{Key: dataKey, Version: 10})
			require.NoError(t, err)
			assert.Equal(t, []byte("plain text"), got)
		})
	}
}
