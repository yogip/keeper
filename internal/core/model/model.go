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

type Tag struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SecretMeta struct {
	ID   int64
	Name string
	// Tags []*Tag     `json:"tags"` todo
	Type SecretType
}

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
	Payload []byte
}

type SecretCreateRequest struct {
	UserID  int64
	Type    SecretType
	Name    string
	Payload []byte
}

type Password struct {
	SecretMeta
	Login    string `json:"login"`
	Password string `json:"password"`
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

type Note struct {
	SecretMeta
	Note string `json:"Note"`
}

type CardData struct {
	Number     string `json:"number"`
	Month      int    `json:"month"`
	Year       int    `json:"year"`
	HolderName string `json:"holder_name"`
	CVC        int    `json:"cvc"`
}

type Card struct {
	SecretMeta
	CardData
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
	Path string
}

type File struct {
	FileMeta
	Body []byte
}

type CreateFileRequest struct {
	UserID int64
	Name   string
	Body   []byte
}
