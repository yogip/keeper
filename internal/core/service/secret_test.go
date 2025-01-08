package service

import (
	"context"
	"keeper/internal/core/model"
)

func (s *TestSuite) TestListSecret() {
	ctx := context.Background()

	type args struct {
		req model.SecretListRequest
	}
	tests := []struct {
		name     string
		args     args
		expected *model.SecretList
	}{
		{
			name: "Test All",
			args: args{
				req: model.SecretListRequest{UserID: 1},
			},
			expected: &model.SecretList{
				Secrets: []*model.SecretMeta{
					{ID: 1, Type: model.SecretTypePassword, Name: "pwd-1"},
					{ID: 2, Type: model.SecretTypePassword, Name: "pwd-2", Note: "note..."},
					{ID: 3, Type: model.SecretTypeNote, Name: "note-1", Note: "note..."},
					{ID: 4, Type: model.SecretTypeNote, Name: "note-2"},
					{ID: 5, Type: model.SecretTypeCard, Name: "card-1", Note: "some note ..."},
				},
			},
		},
	}

	for _, tt := range tests {
		actual, err := s.secretSrv.ListSecretsMeta(ctx, tt.args.req)
		s.Require().NoError(err)
		s.EqualValues(tt.expected, actual)
	}
}

func (s *TestSuite) TestGetSecret() {
	ctx := context.Background()

	type args struct {
		req model.SecretRequest
	}
	tests := []struct {
		name     string
		args     args
		expected *model.Secret
	}{
		{
			name: "Test Password - 1",
			args: args{
				req: model.SecretRequest{ID: 1, UserID: 1, Type: model.SecretTypePassword},
			},
			expected: model.NewSecret(
				1,
				"pwd-1",
				model.SecretTypePassword,
				[]byte(`{"login": "stub-login", "password": "stub-password"}`),
				"",
			),
		},
		{
			name: "Test Password - 2",
			args: args{
				req: model.SecretRequest{ID: 2, UserID: 1, Type: model.SecretTypePassword},
			},
			expected: model.NewSecret(
				2,
				"pwd-2",
				model.SecretTypePassword,
				[]byte(`{"login": "stub-login", "password": "stub-password"}`),
				"note...",
			),
		},
		{
			name: "Test Note - 3",
			args: args{
				req: model.SecretRequest{ID: 3, UserID: 1, Type: model.SecretTypeNote},
			},
			expected: model.NewSecret(
				3,
				"note-1",
				model.SecretTypeNote,
				[]byte(`{"text": "stub-text"}`),
				"note...",
			),
		},
		{
			name: "Test Note - 4",
			args: args{
				req: model.SecretRequest{ID: 4, UserID: 1, Type: model.SecretTypeNote},
			},
			expected: model.NewSecret(
				4,
				"note-2",
				model.SecretTypeNote,
				[]byte(`{"text": "stub-text-2"}`),
				"",
			),
		},
		{
			name: "Test Card - 5",
			args: args{
				req: model.SecretRequest{ID: 5, UserID: 1, Type: model.SecretTypeCard},
			},
			expected: model.NewSecret(
				5,
				"card-1",
				model.SecretTypeCard,
				[]byte(`{"number": "1234 1234 1234 1234", "month": 11, "year": 25, "holder_name": "Holder Name", "cvc": 123}`),
				"some note ...",
			),
		},
	}

	for _, tt := range tests {
		actual, err := s.secretSrv.GetSecret(ctx, tt.args.req)
		s.Require().NoError(err)

		switch tt.args.req.Type {
		case model.SecretTypePassword:
			actPwd, err := actual.AsPassword()
			s.Require().NoError(err)

			pwd, err := tt.expected.AsPassword()
			s.Require().NoError(err)
			s.Assert().Equal(pwd, actPwd)
		case model.SecretTypeNote:
			actNote, err := actual.AsNote()
			s.Require().NoError(err)

			note, err := tt.expected.AsNote()
			s.Require().NoError(err)
			s.Assert().Equal(note, actNote)
		case model.SecretTypeCard:
			actCard, err := actual.AsCard()
			s.Require().NoError(err)

			card, err := tt.expected.AsCard()
			s.Require().NoError(err)
			s.Assert().Equal(card, actCard)
		}

	}
}

func (s *TestSuite) TestFile() {
	ctx := context.Background()

	crReq := model.CreateFileRequest{
		UserID:   2,
		Name:     "test file",
		FileName: "test.txt",
		Note:     "some note ...",
		Payload:  []byte("file content"),
	}

	fileID, err := s.secretSrv.CreateFile(ctx, crReq)
	s.Require().NoError(err)
	s.Assert().GreaterOrEqual(fileID, int64(1))

	req := model.SecretRequest{ID: fileID, UserID: 2, Type: model.SecretTypeFile}
	secret, err := s.secretSrv.GetSecret(ctx, req)
	s.Require().NoError(err)
	s.Assert().GreaterOrEqual(fileID, secret.ID)
	s.Assert().Equal(crReq.Name, secret.Name)

	file, err := secret.AsFile()
	// fmt.Println("---")
	// p, _ := file.GetPayload()
	// fmt.Println(string(p))
	// fmt.Println("---")
	s.Require().NoError(err)

	s.Assert().Equal(crReq.FileName, file.FileName)
	s.Assert().Equal(string(crReq.Payload), string(file.File))
}

func (s *TestSuite) TestCreateSecret() {
	ctx := context.Background()

	type args struct {
		req model.SecretCreateRequest
	}
	tests := []struct {
		name     string
		args     args
		expected *model.Secret
	}{
		{
			name: "Test Password",
			args: args{
				req: model.SecretCreateRequest{
					UserID:  2,
					Type:    model.SecretTypePassword,
					Name:    "pwd-name",
					Note:    "note...",
					Payload: []byte(`{"login": "stub-login", "password": "stub-password"}`),
				},
			},
			expected: model.NewSecret(
				10001,
				"pwd-name",
				model.SecretTypePassword,
				[]byte(`{"login": "stub-login", "password": "stub-password"}`),
				"note...",
			),
		},
		{
			name: "Test Note",
			args: args{
				req: model.SecretCreateRequest{
					UserID:  2,
					Type:    model.SecretTypeNote,
					Name:    "note-1",
					Note:    "some note ...",
					Payload: []byte(`{"text": "stub-text"}`),
				},
			},
			expected: model.NewSecret(
				10002,
				"note-1",
				model.SecretTypeNote,
				[]byte(`{"text": "stub-text"}`),
				"some note ...",
			),
		},
		{
			name: "Test Card",
			args: args{
				req: model.SecretCreateRequest{
					UserID:  2,
					Type:    model.SecretTypeCard,
					Name:    "card-1",
					Note:    "some note ...",
					Payload: []byte(`{"number": "1234 1234 1234 1234", "month": 11, "year": 25, "holder_name": "Holder Name", "cvc": 123}`),
				},
			},
			expected: model.NewSecret(
				10003,
				"card-1",
				model.SecretTypeCard,
				[]byte(`{"number": "1234 1234 1234 1234", "month": 11, "year": 25, "holder_name": "Holder Name", "cvc": 123}`),
				"some note ...",
			),
		},
	}

	for _, tt := range tests {
		actual, err := s.secretSrv.CreateSecret(ctx, tt.args.req)
		s.Require().NoError(err)

		switch tt.args.req.Type {
		case model.SecretTypePassword:
			actPwd, err := actual.AsPassword()
			s.Require().NoError(err)

			pwd, err := tt.expected.AsPassword()
			s.Require().NoError(err)
			s.Assert().Equal(pwd, actPwd)
		case model.SecretTypeNote:
			actNote, err := actual.AsNote()
			s.Require().NoError(err)

			note, err := tt.expected.AsNote()
			s.Require().NoError(err)
			s.Assert().Equal(note, actNote)
		case model.SecretTypeCard:
			actCard, err := actual.AsCard()
			s.Require().NoError(err)

			card, err := tt.expected.AsCard()
			s.Require().NoError(err)
			s.Assert().Equal(card, actCard)
		}

	}
}

func (s *TestSuite) TestUpdateSecret() {
	ctx := context.Background()

	type args struct {
		req model.SecretUpdateRequest
	}
	tests := []struct {
		name     string
		args     args
		expected *model.Secret
	}{
		{
			name: "Test Password",
			args: args{
				req: model.SecretUpdateRequest{
					ID:      6,
					UserID:  2,
					Type:    model.SecretTypePassword,
					Name:    "pwd-name NEW",
					Note:    "note... NEW",
					Payload: []byte(`{"login": "stub-login-NEW", "password": "stub-password-NEW"}`),
				},
			},
			expected: model.NewSecret(
				6,
				"pwd-name NEW",
				model.SecretTypePassword,
				[]byte(`{"login": "stub-login-NEW", "password": "stub-password-NEW"}`),
				"note... NEW",
			),
		},
		{
			name: "Test Note",
			args: args{
				req: model.SecretUpdateRequest{
					ID:      7,
					UserID:  2,
					Type:    model.SecretTypeNote,
					Name:    "note-1 NEW",
					Note:    "some note ... NEW",
					Payload: []byte(`{"text": "stub-text NEW NEW NEW"}`),
				},
			},
			expected: model.NewSecret(
				7,
				"note-1 NEW",
				model.SecretTypeNote,
				[]byte(`{"text": "stub-text NEW NEW NEW"}`),
				"some note ... NEW",
			),
		},
		{
			name: "Test Card",
			args: args{
				req: model.SecretUpdateRequest{
					ID:      8,
					UserID:  2,
					Type:    model.SecretTypeCard,
					Name:    "card-1 NEW",
					Note:    "some note ... NEW",
					Payload: []byte(`{"number": "1234 0000 1234 1230", "month": 10, "year": 26, "holder_name": "Holder Name NEW", "cvc": 321}`),
				},
			},
			expected: model.NewSecret(
				8,
				"card-1 NEW",
				model.SecretTypeCard,
				[]byte(`{"number": "1234 0000 1234 1230", "month": 10, "year": 26, "holder_name": "Holder Name NEW", "cvc": 321}`),
				"some note ... NEW",
			),
		},
	}

	for _, tt := range tests {
		actual, err := s.secretSrv.UpdateSecret(ctx, tt.args.req)
		s.Require().NoError(err)

		switch tt.args.req.Type {
		case model.SecretTypePassword:
			actPwd, err := actual.AsPassword()
			s.Require().NoError(err)

			pwd, err := tt.expected.AsPassword()
			s.Require().NoError(err)
			s.Assert().Equal(pwd, actPwd)
		case model.SecretTypeNote:
			actNote, err := actual.AsNote()
			s.Require().NoError(err)

			note, err := tt.expected.AsNote()
			s.Require().NoError(err)
			s.Assert().Equal(note, actNote)
		case model.SecretTypeCard:
			actCard, err := actual.AsCard()
			s.Require().NoError(err)

			card, err := tt.expected.AsCard()
			s.Require().NoError(err)
			s.Assert().Equal(card, actCard)
		}

	}
}
