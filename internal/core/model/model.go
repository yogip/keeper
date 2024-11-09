package model

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	ID           int64   `json:"id"`
	Login        string  `json:"login"`
	PasswordHash *[]byte `json:"password,omitempty"`
}

type SecretType string

const (
	SecretTypePassword SecretType = "password"
	SecretTypeOTP      SecretType = "otp"
	SecretTypeNote     SecretType = "note"
	SecretTypeFile     SecretType = "file"
)

type Tag struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Folder struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type SecretMeta struct {
	ID     int64      `json:"id"`
	Name   string     `json:"name"`
	Tags   []*Tag     `json:"tags"`
	Folder Folder     `json:"folder"`
	Type   SecretType `json:"type"`
}

type DataKey struct {
	Version int64
	Key     []byte
}

type SecretList struct {
	Secrets []*SecretMeta `json:"secrets"`
}

type SecretListRequest struct {
	UserID int64 `json:"user_id"`
}

type SecretRequest struct {
	ID     int64      `json:"id"`
	UserID int64      `json:"user_id"`
	Type   SecretType `json:"type"`
}

type Password struct {
	SecretMeta
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Note struct {
	SecretMeta
	Note string `json:"Note"`
}
