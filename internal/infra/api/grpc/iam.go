package grpc

import (
	"context"
	"errors"
	coreErrors "keeper/internal/core/errors"
	"keeper/internal/core/model"
	"keeper/internal/logger"
	pb "keeper/internal/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *KeeperServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.Token, error) {
	log := logger.Log.With(zap.Any("Login", in.Login))
	log.Info("Login request")

	token, err := s.iam.Login(
		ctx,
		&model.UserRequest{
			Login:    in.Login,
			Password: in.Password,
		},
	)
	if err != nil && errors.Is(err, coreErrors.ErrNotFound404) {
		return nil, status.Errorf(codes.NotFound, "Not found")
	}
	if err != nil {
		log.Error("Login error", zap.Error(err))
		return nil, status.Errorf(codes.Unknown, "Login error: %s", err)
	}

	return &pb.Token{Token: token}, nil
}

func (s *KeeperServer) SignUp(ctx context.Context, in *pb.SignUpRequest) (*pb.Token, error) {
	log := logger.Log.With(zap.Any("Login", in.Login))
	log.Info("Register request")

	token, err := s.iam.Register(
		ctx,
		&model.UserRequest{
			Login:    in.Login,
			Password: in.Password,
		},
	)
	if err != nil && errors.Is(err, coreErrors.ErrNotFound404) {
		return nil, status.Errorf(codes.NotFound, "Not found")
	}
	if err != nil {
		log.Error("Register error", zap.Error(err))
		return nil, status.Errorf(codes.Unknown, "Register error: %s", err)
	}

	return &pb.Token{Token: token}, nil
}
