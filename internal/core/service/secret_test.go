package service

import (
	"context"
	"keeper/internal/core/model"
)

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
		// {
		// 	name: "Test File - 6",
		// 	args: args{
		// 		req: model.SecretRequest{ID: 6, UserID: 1, Type: model.SecretTypeFile},
		// 	},
		// 	expected: model.NewSecret(
		// 		6,
		// 		"file-1",
		// 		model.SecretTypeCard,
		// 		[]byte(`{"number": "1234 1234 1234 1234", "month": 11, "year": 25, "holder_name": "Holder Name", "cvc": 123}`),
		// 		"some note ...",
		// 	),
		// },
	}

	// S3Name   string `json:"s3_name"`
	// FileName string `json:"file_name"`
	// File     []byte `json:"file"`

	for _, tt := range tests {
		// s.encript(`{"number": "1234 1234 1234 1234", "month": 11, "year": 25, "holder_name": "Holder Name", "cvc": 123}`)
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

// func (s *TestSuite) TestGetNote() {
// 	ctx := context.Background()

// 	type args struct {
// 		req model.SecretRequest
// 	}
// 	tests := []struct {
// 		name     string
// 		args     args
// 		expected model.Note
// 	}{
// 		{
// 			name: "Test success - 1",
// 			args: args{
// 				req: model.SecretRequest{ID: 1, UserID: 1, Type: model.SecretTypeNote},
// 			},
// 			expected: model.Note{
// 				SecretMeta: model.SecretMeta{
// 					ID:   1,
// 					Name: "note-1",
// 					Type: model.SecretTypeNote,
// 					Note: "stub-note",
// 				},
// 			},
// 		},
// 		{
// 			name: "Test success - 2",
// 			args: args{
// 				req: model.SecretRequest{ID: 2, UserID: 1, Type: model.SecretTypeNote},
// 			},
// 			expected: model.Note{
// 				SecretMeta: model.SecretMeta{
// 					ID:   2,
// 					Name: "note-2",
// 					Type: model.SecretTypeNote,
// 					Note: "stub-note",
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		actual, err := s.secretSrv.GetNote(ctx, tt.args.req)
// 		s.Require().NoError(err)
// 		s.Assert().Equal(tt.expected, *actual)
// 	}
// }

// func (s *TestSuite) TestGetCard() {
// 	ctx := context.Background()

// 	req := model.SecretRequest{ID: 1, UserID: 1, Type: model.SecretTypeCard}
// 	expected := model.Card{
// 		SecretMeta: model.SecretMeta{
// 			ID:   1,
// 			Name: "card-1",
// 			Type: model.SecretTypeCard,
// 		},
// 		Number:     "1234123412341234",
// 		Month:      12,
// 		Year:       25,
// 		HolderName: "No Name",
// 		CVC:        123,
// 	}

// 	card, err := s.secretSrv.GetCard(ctx, req)
// 	s.Require().NoError(err)
// 	s.Assert().Equal(expected, *card)
// }

// func (s *TestSuite) TestList() {
// 	ctx := context.Background()

// 	type args struct {
// 		req model.SecretListRequest
// 	}
// 	tests := []struct {
// 		name     string
// 		args     args
// 		expected []model.SecretMeta
// 	}{
// 		{
// 			name: "Test passwords",
// 			args: args{
// 				req: model.SecretListRequest{UserID: 1, Name: "pwd"},
// 			},
// 			expected: []model.SecretMeta{
// 				{
// 					ID:   1,
// 					Name: "pwd-1",
// 					Type: model.SecretTypePassword,
// 				},
// 				{
// 					ID:   2,
// 					Name: "pwd-2",
// 					Type: model.SecretTypePassword,
// 				},
// 			},
// 		},
// 		{
// 			name: "Test notes",
// 			args: args{
// 				req: model.SecretListRequest{UserID: 1, Name: "note"},
// 			},
// 			expected: []model.SecretMeta{
// 				{
// 					ID:   1,
// 					Name: "note-1",
// 					Type: model.SecretTypeNote,
// 				},
// 				{
// 					ID:   2,
// 					Name: "note-2",
// 					Type: model.SecretTypeNote,
// 				},
// 			},
// 		},
// 		{
// 			name: "Test cards",
// 			args: args{
// 				req: model.SecretListRequest{UserID: 1, Name: "card-1"},
// 			},
// 			expected: []model.SecretMeta{
// 				{
// 					ID:   1,
// 					Name: "card-1",
// 					Type: model.SecretTypeCard,
// 				},
// 			},
// 		},
// 		{
// 			name: "Test cards - empty",
// 			args: args{
// 				req: model.SecretListRequest{UserID: 3},
// 			},
// 			expected: []model.SecretMeta{},
// 		},
// 	}

// 	for _, tt := range tests {
// 		actual, err := s.secretSrv.ListSecretsMeta(ctx, &tt.args.req)
// 		s.Require().NoError(err)
// 		s.Assert().Equal(len(actual.Secrets), len(tt.expected))
// 		for _, v := range actual.Secrets {
// 			s.Assert().Contains(tt.expected, *v)
// 		}
// 	}
// }

// func (s *TestSuite) TestCreatePassword() {
// 	ctx := context.Background()

// 	req := model.UpdatePasswordRequest{
// 		UserID: 2,
// 		Data: &model.Password{
// 			SecretMeta: model.SecretMeta{Name: "pwd-1"},
// 			Login:      "login",
// 			Password:   "pwdd",
// 		},
// 	}
// 	expected := model.Password{
// 		SecretMeta: model.SecretMeta{
// 			ID:   3,
// 			Name: "pwd-1",
// 			Type: model.SecretTypePassword,
// 		},
// 		Login:    "login",
// 		Password: "pwdd",
// 	}

// 	actual, err := s.secretSrv.CreatePassword(ctx, req)
// 	s.Require().NoError(err)
// 	s.Assert().Equal(expected.Login, actual.Login)
// 	s.Assert().Equal(expected.Name, actual.Name)
// 	s.Assert().Equal(expected.Password, actual.Password)
// }

// func (s *TestSuite) TestUpdatePassword() {
// 	ctx := context.Background()

// 	req := model.UpdatePasswordRequest{
// 		UserID: 1,
// 		Data: &model.Password{
// 			SecretMeta: model.SecretMeta{
// 				ID:   1,
// 				Name: "new-name",
// 			},
// 			Login:    "new-login",
// 			Password: "new-pwdd",
// 		},
// 	}
// 	expected := model.Password{
// 		SecretMeta: model.SecretMeta{
// 			ID:   1,
// 			Name: "new-name",
// 			Type: model.SecretTypePassword,
// 		},
// 		Login:    "new-login",
// 		Password: "new-pwdd",
// 	}

// 	actual, err := s.secretSrv.UpdatePassword(ctx, req)
// 	s.Require().NoError(err)
// 	s.Assert().Equal(expected, *actual)
// }

// func (s *TestSuite) TestCreateNote() {
// 	ctx := context.Background()

// 	req := model.UpdateNoteRequest{
// 		UserID: 2,
// 		Data: &model.Note{
// 			SecretMeta: model.SecretMeta{
// 				Name: "note-name",
// 				Note: "note ................................",
// 			},
// 		},
// 	}
// 	expected := model.Note{
// 		SecretMeta: model.SecretMeta{
// 			ID:   3,
// 			Name: "note-name",
// 			Type: model.SecretTypeNote,
// 			Note: "note ................................",
// 		},
// 	}

// 	actual, err := s.secretSrv.CreateNote(ctx, req)
// 	s.Require().NoError(err)
// 	s.Assert().Equal(expected.Note, actual.Note)
// 	s.Assert().Equal(expected.Name, actual.Name)
// }

// func (s *TestSuite) TestUpdateNote() {
// 	ctx := context.Background()

// 	req := model.UpdateNoteRequest{
// 		UserID: 1,
// 		Data: &model.Note{
// 			SecretMeta: model.SecretMeta{
// 				ID:   1,
// 				Name: "new-name",
// 				Note: "new-note",
// 			},
// 		},
// 	}
// 	expected := model.Note{
// 		SecretMeta: model.SecretMeta{
// 			ID:   1,
// 			Name: "new-name",
// 			Type: model.SecretTypeNote,
// 			Note: "new-note",
// 		},
// 	}

// 	actual, err := s.secretSrv.UpdateNote(ctx, req)
// 	s.Require().NoError(err)
// 	s.Assert().Equal(expected, *actual)
// }

// func (s *TestSuite) TestCreateCard() {
// 	ctx := context.Background()

// 	req := model.UpdateCardRequest{
// 		UserID: 2,
// 		Card: model.Card{
// 			SecretMeta: model.SecretMeta{
// 				Name: "card-1",
// 				Type: model.SecretTypeCard,
// 			},
// 			Number:     "1234123412341234",
// 			Month:      12,
// 			Year:       25,
// 			HolderName: "No Name",
// 			CVC:        123,
// 		},
// 	}
// 	expected := model.Card{
// 		SecretMeta: model.SecretMeta{
// 			ID:   1,
// 			Name: "card-1",
// 			Type: model.SecretTypeCard,
// 		},
// 		Number:     "1234123412341234",
// 		Month:      12,
// 		Year:       25,
// 		HolderName: "No Name",
// 		CVC:        123,
// 	}

// 	actual, err := s.secretSrv.CreateCard(ctx, req)

// 	s.Require().NoError(err)
// 	s.Assert().Equal(expected.Name, actual.Name)
// 	s.Assert().Equal(expected.Type, actual.Type)
// 	s.Assert().Equal(expected.Number, actual.Number)
// 	s.Assert().Equal(expected.Month, actual.Month)
// 	s.Assert().Equal(expected.Year, actual.Year)
// 	s.Assert().Equal(expected.HolderName, actual.HolderName)
// 	s.Assert().Equal(expected.CVC, actual.CVC)
// }

// func (s *TestSuite) TestUpdateCard() {
// 	ctx := context.Background()

// 	req := model.UpdateCardRequest{
// 		UserID: 2,
// 		Card: model.Card{
// 			SecretMeta: model.SecretMeta{
// 				ID:   2,
// 				Name: "card-new",
// 				Type: model.SecretTypeCard,
// 			},
// 			Number:     "0234123412341234",
// 			Month:      10,
// 			Year:       22,
// 			HolderName: "No NewName",
// 			CVC:        321,
// 		},
// 	}
// 	expected := model.Card{
// 		SecretMeta: model.SecretMeta{
// 			ID:   2,
// 			Name: "card-new",
// 			Type: model.SecretTypeCard,
// 		},
// 		Number:     "0234123412341234",
// 		Month:      10,
// 		Year:       22,
// 		HolderName: "No NewName",
// 		CVC:        321,
// 	}

// 	actual, err := s.secretSrv.UpdateCard(ctx, req)

// 	s.Require().NoError(err)
// 	s.Assert().Equal(expected, *actual)
// }

// func (s *TestSuite) TestFile() {
// 	ctx := context.Background()

// 	crReq := model.CreateFileRequest{
// 		UserID: 1,
// 		Body:   []byte("test"),
// 		Name:   "test file",
// 	}

// 	fileMeta, err := s.secretSrv.CreateFile(ctx, crReq)
// 	s.Require().NoError(err)
// 	s.Assert().GreaterOrEqual(fileMeta.ID, int64(1))
// 	s.Assert().Equal(fileMeta.Name, crReq.Name)
// 	_, err = uuid.Parse(fileMeta.Path)
// 	s.Require().NoError(err)

// 	req := model.SecretRequest{ID: fileMeta.ID, UserID: 1, Type: model.SecretTypeFile}
// 	file, err := s.secretSrv.GetFile(ctx, req)
// 	s.Require().NoError(err)
// 	s.Assert().GreaterOrEqual(fileMeta.ID, file.ID)
// 	s.Assert().Equal(fileMeta.Name, file.Name)
// 	s.Assert().Equal(fileMeta.Path, file.Path)
// 	s.Assert().Equal(string(crReq.Body), string(file.Body))
// }
