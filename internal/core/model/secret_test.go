package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSecret(t *testing.T) {
	type args struct {
		id         int64
		name       string
		secretType SecretType
		data       []byte
		note       string
	}
	tests := []struct {
		name string
		args args
		want *Secret
	}{
		{
			name: "Test with valid arguments",
			args: args{
				id:         1,
				name:       "test_name",
				secretType: SecretTypePassword,
				data:       []byte("test_data"),
				note:       "test_note",
			},
			want: &Secret{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypePassword,
					Note: "test_note",
				},
				Payload: []byte("test_data"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSecret(tt.args.id, tt.args.name, tt.args.secretType, tt.args.data, tt.args.note)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSecret_AsPassword(t *testing.T) {
	type fields struct {
		SecretMeta SecretMeta
		Payload    []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Password
		wantErr bool
	}{
		{
			name: "Test with valid arguments",
			fields: fields{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypePassword,
					Note: "test_note",
				},
				Payload: []byte(`{"login":"test_login","password":"test_password"}`),
			},
			want: &Password{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypePassword,
					Note: "test_note",
				},
				Login:    "test_login",
				Password: "test_password",
			},
			wantErr: false,
		},
		{
			name: "Test with invalid secret type",
			fields: fields{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypeNote,
					Note: "test_note",
				},
				Payload: []byte(`{"text":"test_text","password":"test_password"}`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test with invalid payload",
			fields: fields{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypePassword,
					Note: "test_note",
				},
				Payload: []byte(`{"login":"test_login","password":"test_password"`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Secret{
				SecretMeta: tt.fields.SecretMeta,
				Payload:    tt.fields.Payload,
			}
			got, err := s.AsPassword()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSecret_AsNote(t *testing.T) {
	type fields struct {
		SecretMeta SecretMeta
		Payload    []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Note
		wantErr bool
	}{
		{
			name: "Test with valid arguments",
			fields: fields{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypeNote,
					Note: "test_note",
				},
				Payload: []byte(`{"text":"test_text"}`),
			},
			want: &Note{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypeNote,
					Note: "test_note",
				},
				Text: "test_text",
			},
			wantErr: false,
		},
		{
			name: "Test with invalid secret type",
			fields: fields{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypePassword,
					Note: "test_note",
				},
				Payload: []byte(`{"text":"test_text"}`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test with invalid payload",
			fields: fields{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypeNote,
					Note: "test_note",
				},
				Payload: []byte(`{"text":"test_text"`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Secret{
				SecretMeta: tt.fields.SecretMeta,
				Payload:    tt.fields.Payload,
			}
			got, err := s.AsNote()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSecret_AsCard(t *testing.T) {
	type fields struct {
		SecretMeta SecretMeta
		Payload    []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Card
		wantErr bool
	}{
		{
			name: "Test with valid arguments",
			fields: fields{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypeCard,
					Note: "test_note",
				},
				Payload: []byte(`{"number":"test_number","month":1,"year":2022,"holder_name":"test_holder_name","cvc":123}`),
			},
			want: &Card{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypeCard,
					Note: "test_note",
				},
				Number:     "test_number",
				Month:      1,
				Year:       2022,
				HolderName: "test_holder_name",
				CVC:        123,
			},
			wantErr: false,
		},
		{
			name: "Test with invalid secret type",
			fields: fields{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypePassword,
					Note: "test_note",
				},
				Payload: []byte(`{"number":"test_number","month":1,"year":2022,"holder_name":"test_holder_name","cvc":123}`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test with invalid payload",
			fields: fields{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypeCard,
					Note: "test_note",
				},
				Payload: []byte(`{"number":"test_number","month":1,"year":2022,"holder_name":"test_holder_name","cvc":123"`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Secret{
				SecretMeta: tt.fields.SecretMeta,
				Payload:    tt.fields.Payload,
			}
			got, err := s.AsCard()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSecret_AsFile(t *testing.T) {
	type fields struct {
		SecretMeta SecretMeta
		Payload    []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    *File
		wantErr bool
	}{
		{
			name: "Test with valid arguments",
			fields: fields{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypeFile,
					Note: "test_note",
				},
				Payload: []byte(`{"s3_name":"test_s3_name","file_name":"test_file_name","file":"dGVzdF9maWxl"}`),
			},
			want: &File{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypeFile,
					Note: "test_note",
				},
				S3Name:   "test_s3_name",
				FileName: "test_file_name",
				File:     []byte("test_file"),
			},
			wantErr: false,
		},
		{
			name: "Test with invalid secret type",
			fields: fields{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypePassword,
					Note: "test_note",
				},
				Payload: []byte(`{"s3_name":"test_s3_name","file_name":"test_file_name","file":"dGVzdF9maWxl"}`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test with invalid payload",
			fields: fields{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypeFile,
					Note: "test_note",
				},
				Payload: []byte(`{"s3_name":"test_s3_name","file_name":"test_file_name","file":"dGVzdF9maWxl"`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Secret{
				SecretMeta: tt.fields.SecretMeta,
				Payload:    tt.fields.Payload,
			}
			got, err := s.AsFile()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPassword_GetPayload(t *testing.T) {
	type fields struct {
		SecretMeta SecretMeta
		Login      string
		Password   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Test with valid arguments",
			fields: fields{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypePassword,
					Note: "test_note",
				},
				Login:    "test_login",
				Password: "test_password",
			},
			want:    []byte(`{"ID":1,"Name":"test_name","Type":"password","note":"test_note","login":"test_login","password":"test_password"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Password{
				SecretMeta: tt.fields.SecretMeta,
				Login:      tt.fields.Login,
				Password:   tt.fields.Password,
			}
			got, err := p.GetPayload()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNote_GetPayload(t *testing.T) {
	type fields struct {
		SecretMeta SecretMeta
		Text       string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Test with valid arguments",
			fields: fields{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypeNote,
					Note: "test_note",
				},
				Text: "test_text",
			},
			want:    []byte(`{"ID":1,"Name":"test_name","Type":"note","note":"test_note","text":"test_text"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Note{
				SecretMeta: tt.fields.SecretMeta,
				Text:       tt.fields.Text,
			}
			got, err := n.GetPayload()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCard_GetPayload(t *testing.T) {
	type fields struct {
		SecretMeta SecretMeta
		Number     string
		Month      int
		Year       int
		HolderName string
		CVC        int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Test with valid arguments",
			fields: fields{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypeCard,
					Note: "test_note",
				},
				Number:     "test_number",
				Month:      1,
				Year:       2022,
				HolderName: "test_holder_name",
				CVC:        123,
			},
			want:    []byte(`{"ID":1,"Name":"test_name","Type":"card","note":"test_note","number":"test_number","month":1,"year":2022,"holder_name":"test_holder_name","cvc":123}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Card{
				SecretMeta: tt.fields.SecretMeta,
				Number:     tt.fields.Number,
				Month:      tt.fields.Month,
				Year:       tt.fields.Year,
				HolderName: tt.fields.HolderName,
				CVC:        tt.fields.CVC,
			}
			got, err := c.GetPayload()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFile_GetPayload(t *testing.T) {
	type fields struct {
		SecretMeta SecretMeta
		S3Name     string
		FileName   string
		File       []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Test with valid arguments",
			fields: fields{
				SecretMeta: SecretMeta{
					ID:   1,
					Name: "test_name",
					Type: SecretTypeFile,
					Note: "test_note",
				},
				S3Name:   "test_s3_name",
				FileName: "test_file_name",
				File:     []byte("test_file"),
			},
			want:    []byte(`{"ID":1,"Name":"test_name","Type":"file","note":"test_note","s3_name":"test_s3_name","file_name":"test_file_name","file":"dGVzdF9maWxl"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				SecretMeta: tt.fields.SecretMeta,
				S3Name:     tt.fields.S3Name,
				FileName:   tt.fields.FileName,
				File:       tt.fields.File,
			}
			got, err := f.GetPayload()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
