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

func (s *TestSuite) TestGetList() {
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
				req: model.SecretListRequest{UserID: 2, Type: model.SecretTypeCard},
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
