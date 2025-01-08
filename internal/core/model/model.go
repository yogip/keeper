package model

type SecretType string

const (
	SecretTypePassword SecretType = "password"
	SecretTypeCard     SecretType = "card"
	SecretTypeNote     SecretType = "note"
	SecretTypeFile     SecretType = "file"
)

type DataKey struct {
	Version int64
	Key     string
}

type SecretList struct {
	Secrets []*SecretMeta `json:"secrets"`
}

type SecretListRequest struct {
	UserID int64
	Name   string
}

type SecretRequest struct {
	ID     int64
	UserID int64
	Type   SecretType
}

type SecretUpdateRequest struct {
	ID      int64
	UserID  int64
	Type    SecretType
	Name    string
	Note    string
	Payload []byte
}

type SecretCreateRequest struct {
	UserID  int64
	Type    SecretType
	Name    string
	Note    string
	Payload []byte
}

type CreateFileRequest struct {
	UserID   int64
	Name     string
	FileName string
	Note     string
	Payload  []byte
}
