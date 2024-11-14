package service

import (
	"context"
	"keeper/internal/core/model"
)

func (s *TestSuite) TestGetPassword() {
	ctx := context.Background()

	type args struct {
		req model.SecretRequest
	}
	tests := []struct {
		name     string
		args     args
		expected model.Password
	}{
		{
			name: "Test success - 1",
			args: args{
				req: model.SecretRequest{ID: 1, UserID: 1, Type: model.SecretTypePassword},
			},
			expected: model.Password{
				SecretMeta: model.SecretMeta{
					ID:   1,
					Name: "pwd-1",
					Tags: nil,
					Type: model.SecretTypePassword,
				},
				Login:    "stub-login",
				Password: "stub-password",
			},
		},
		{
			name: "Test success - 2",
			args: args{
				req: model.SecretRequest{ID: 2, UserID: 1, Type: model.SecretTypePassword},
			},
			expected: model.Password{
				SecretMeta: model.SecretMeta{
					ID:   2,
					Name: "pwd-2",
					Tags: nil,
					Type: model.SecretTypePassword,
				},
				Login:    "stub-login",
				Password: "stub-password",
			},
		},
	}

	for _, tt := range tests {
		actual, err := s.secretSrv.GetPassword(ctx, tt.args.req)
		s.Require().NoError(err)
		s.Assert().Equal(tt.expected, *actual)
	}
}

func (s *TestSuite) TestGetNote() {
	ctx := context.Background()

	type args struct {
		req model.SecretRequest
	}
	tests := []struct {
		name     string
		args     args
		expected model.Note
	}{
		{
			name: "Test success - 1",
			args: args{
				req: model.SecretRequest{ID: 1, UserID: 1, Type: model.SecretTypeNote},
			},
			expected: model.Note{
				SecretMeta: model.SecretMeta{
					ID:   1,
					Name: "note-1",
					Tags: nil,
					Type: model.SecretTypeNote,
				},
				Note: "stub-note",
			},
		},
		{
			name: "Test success - 2",
			args: args{
				req: model.SecretRequest{ID: 2, UserID: 1, Type: model.SecretTypeNote},
			},
			expected: model.Note{
				SecretMeta: model.SecretMeta{
					ID:   2,
					Name: "note-2",
					Tags: nil,
					Type: model.SecretTypeNote,
				},
				Note: "stub-note",
			},
		},
	}

	for _, tt := range tests {
		actual, err := s.secretSrv.GetNote(ctx, tt.args.req)
		s.Require().NoError(err)
		s.Assert().Equal(tt.expected, *actual)
	}
}

func (s *TestSuite) TestGetCard() {
	ctx := context.Background()

	req := model.SecretRequest{ID: 1, UserID: 1, Type: model.SecretTypeCard}
	expected := model.Card{
		SecretMeta: model.SecretMeta{
			ID:   1,
			Name: "card-1",
			Tags: nil,
			Type: model.SecretTypeCard,
		},
		CardData: model.CardData{
			Number:     "1234123412341234",
			Month:      12,
			Year:       25,
			HolderName: "No Name",
			CVC:        123,
		},
	}

	card, err := s.secretSrv.GetCard(ctx, req)
	s.Require().NoError(err)
	s.Assert().Equal(expected, *card)
}

func (s *TestSuite) TestList() {
	ctx := context.Background()

	type args struct {
		req model.SecretListRequest
	}
	tests := []struct {
		name     string
		args     args
		expected []model.SecretMeta
	}{
		{
			name: "Test passwords",
			args: args{
				req: model.SecretListRequest{UserID: 1, Type: model.SecretTypePassword},
			},
			expected: []model.SecretMeta{
				{
					ID:   1,
					Name: "pwd-1",
					Type: model.SecretTypePassword,
				},
				{
					ID:   2,
					Name: "pwd-2",
					Type: model.SecretTypePassword,
				},
			},
		},
		{
			name: "Test notes",
			args: args{
				req: model.SecretListRequest{UserID: 1, Type: model.SecretTypeNote},
			},
			expected: []model.SecretMeta{
				{
					ID:   1,
					Name: "note-1",
					Type: model.SecretTypeNote,
				},
				{
					ID:   2,
					Name: "note-2",
					Type: model.SecretTypeNote,
				},
			},
		},
		{
			name: "Test cards",
			args: args{
				req: model.SecretListRequest{UserID: 1, Type: model.SecretTypeCard},
			},
			expected: []model.SecretMeta{
				{
					ID:   1,
					Name: "card-1",
					Type: model.SecretTypeCard,
				},
			},
		},
		{
			name: "Test cards - empty",
			args: args{
				req: model.SecretListRequest{UserID: 3, Type: model.SecretTypeCard},
			},
			expected: []model.SecretMeta{},
		},
	}

	for _, tt := range tests {
		actual, err := s.secretSrv.ListSecretsMeta(ctx, &tt.args.req)
		s.Require().NoError(err)
		s.Assert().Equal(len(actual.Secrets), len(tt.expected))
		for _, v := range actual.Secrets {
			s.Assert().Contains(tt.expected, *v)
		}
	}
}

func (s *TestSuite) TestCreatePassword() {
	ctx := context.Background()

	req := model.UpdatePasswordRequest{
		UserID: 2,
		Data: &model.Password{
			SecretMeta: model.SecretMeta{Name: "pwd-1"},
			Login:      "login",
			Password:   "pwdd",
		},
	}
	expected := model.Password{
		SecretMeta: model.SecretMeta{
			ID:   3,
			Name: "pwd-1",
			Type: model.SecretTypePassword,
		},
		Login:    "login",
		Password: "pwdd",
	}

	actual, err := s.secretSrv.CreatePassword(ctx, req)
	s.Require().NoError(err)
	s.Assert().Equal(expected.Login, actual.Login)
	s.Assert().Equal(expected.Name, actual.Name)
	s.Assert().Equal(expected.Password, actual.Password)
}

func (s *TestSuite) TestUpdatePassword() {
	ctx := context.Background()

	req := model.UpdatePasswordRequest{
		UserID: 1,
		Data: &model.Password{
			SecretMeta: model.SecretMeta{
				ID:   1,
				Name: "new-name",
			},
			Login:    "new-login",
			Password: "new-pwdd",
		},
	}
	expected := model.Password{
		SecretMeta: model.SecretMeta{
			ID:   1,
			Name: "new-name",
			Type: model.SecretTypePassword,
		},
		Login:    "new-login",
		Password: "new-pwdd",
	}

	actual, err := s.secretSrv.UpdatePassword(ctx, req)
	s.Require().NoError(err)
	s.Assert().Equal(expected, *actual)
}

func (s *TestSuite) TestCreateNote() {
	ctx := context.Background()

	req := model.UpdateNoteRequest{
		UserID: 2,
		Data: &model.Note{
			SecretMeta: model.SecretMeta{Name: "note-name"},
			Note:       "note ................................",
		},
	}
	expected := model.Note{
		SecretMeta: model.SecretMeta{
			ID:   3,
			Name: "note-name",
			Type: model.SecretTypeNote,
		},
		Note: "note ................................",
	}

	actual, err := s.secretSrv.CreateNote(ctx, req)
	s.Require().NoError(err)
	s.Assert().Equal(expected.Note, actual.Note)
	s.Assert().Equal(expected.Name, actual.Name)
}

func (s *TestSuite) TestUpdateNote() {
	ctx := context.Background()

	req := model.UpdateNoteRequest{
		UserID: 1,
		Data: &model.Note{
			SecretMeta: model.SecretMeta{
				ID:   1,
				Name: "new-name",
			},
			Note: "new-note",
		},
	}
	expected := model.Note{
		SecretMeta: model.SecretMeta{
			ID:   1,
			Name: "new-name",
			Type: model.SecretTypeNote,
		},
		Note: "new-note",
	}

	actual, err := s.secretSrv.UpdateNote(ctx, req)
	s.Require().NoError(err)
	s.Assert().Equal(expected, *actual)
}

func (s *TestSuite) TestCreateCard() {
	ctx := context.Background()

	req := model.UpdateCardRequest{
		UserID: 2,
		Card: model.Card{
			SecretMeta: model.SecretMeta{
				Name: "card-1",
				Type: model.SecretTypeCard,
			},
			CardData: model.CardData{
				Number:     "1234123412341234",
				Month:      12,
				Year:       25,
				HolderName: "No Name",
				CVC:        123,
			},
		},
	}
	expected := model.Card{
		SecretMeta: model.SecretMeta{
			ID:   1,
			Name: "card-1",
			Tags: nil,
			Type: model.SecretTypeCard,
		},
		CardData: model.CardData{
			Number:     "1234123412341234",
			Month:      12,
			Year:       25,
			HolderName: "No Name",
			CVC:        123,
		},
	}

	actual, err := s.secretSrv.CreateCard(ctx, req)

	s.Require().NoError(err)
	s.Assert().Equal(expected.Name, actual.Name)
	s.Assert().Equal(expected.Type, actual.Type)
	s.Assert().Equal(expected.Number, actual.Number)
	s.Assert().Equal(expected.Month, actual.Month)
	s.Assert().Equal(expected.Year, actual.Year)
	s.Assert().Equal(expected.HolderName, actual.HolderName)
	s.Assert().Equal(expected.CVC, actual.CVC)
}

func (s *TestSuite) TestUpdateCard() {
	ctx := context.Background()

	req := model.UpdateCardRequest{
		UserID: 2,
		Card: model.Card{
			SecretMeta: model.SecretMeta{
				ID:   2,
				Name: "card-new",
				Type: model.SecretTypeCard,
			},
			CardData: model.CardData{
				Number:     "0234123412341234",
				Month:      10,
				Year:       22,
				HolderName: "No NewName",
				CVC:        321,
			},
		},
	}
	expected := model.Card{
		SecretMeta: model.SecretMeta{
			ID:   2,
			Name: "card-new",
			Type: model.SecretTypeCard,
		},
		CardData: model.CardData{
			Number:     "0234123412341234",
			Month:      10,
			Year:       22,
			HolderName: "No NewName",
			CVC:        321,
		},
	}

	actual, err := s.secretSrv.UpdateCard(ctx, req)

	s.Require().NoError(err)
	s.Assert().Equal(expected, *actual)
}

func (s *TestSuite) TestFile() {
	ctx := context.Background()

	crReq := model.CreateFileRequest{
		UserID: 1,
		File: model.File{
			Body: []byte("test"),
			FileMeta: model.FileMeta{
				Path:       "test.txt",
				SecretMeta: model.SecretMeta{Name: "test file"},
			},
		},
	}

	fileMeta, err := s.secretSrv.CreateFile(ctx, crReq)
	s.Require().NoError(err)
	s.Assert().GreaterOrEqual(fileMeta.ID, int64(1))
	s.Assert().Equal(fileMeta.Name, crReq.File.Name)
	s.Assert().Equal(fileMeta.Path, crReq.File.Path)

	req := model.SecretRequest{ID: fileMeta.ID, UserID: 1, Type: model.SecretTypeFile}
	file, err := s.secretSrv.GetFile(ctx, req)
	s.Require().NoError(err)
	s.Assert().GreaterOrEqual(fileMeta.ID, file.ID)
	s.Assert().Equal(fileMeta.Name, file.Name)
	s.Assert().Equal(fileMeta.Path, file.Path)
	s.Assert().Equal(string(crReq.File.Body), string(file.Body))
}
