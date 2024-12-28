package model

const UserCtxKey string = "user"

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	ID           int64   `json:"id"`
	Login        string  `json:"login"`
	PasswordHash *[]byte `json:"password,omitempty"`
}

type Token string

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

type UpdatePasswordRequest struct {
	UserID int64
	Data   *Password
	Key    *DataKey
}

type UpdateNoteRequest struct {
	UserID int64
	Data   *Note
	Key    *DataKey
}

type EncryptedCard struct {
	Payload string
	Meta    *SecretMeta
	DataKey *DataKey
}

type UpdateCardRequest struct {
	UserID int64
	Card   Card
	Key    *DataKey
}

type FileMeta struct {
	SecretMeta
	FileName string
}

type CreateFileRequest struct {
	UserID   int64
	Name     string
	FileName string
	Note     string
	Payload  []byte
}
