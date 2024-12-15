package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"slices"

	"keeper/internal/core/config"
	coreErrors "keeper/internal/core/errors"
	"keeper/internal/core/model"
	"keeper/internal/core/service"
	"keeper/internal/logger"
	pb "keeper/internal/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// KeeperServer реализует gRPC сервер метрик
type KeeperServer struct {
	pb.UnimplementedKeeperServer

	cfg           *config.Config
	secretService *service.SecretService
	iam           *service.IAM

	srv *grpc.Server
}

func iamInterceptor(iam *service.IAM) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger.Log.Debug(fmt.Sprintf("------- Got request from: %s", info.FullMethod))
		if slices.Contains([]string{"/proto.Keeper/Login", "/proto.Keeper/SignUp"}, info.FullMethod) {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Log.Error("failed to get metadata")
			return nil, status.Error(codes.Unauthenticated, "Access denied")
		}

		token := md.Get("token")
		if len(token) == 0 {
			logger.Log.Error("got empty token")
			return nil, status.Error(codes.Unauthenticated, "Access denied")
		}

		user, err := iam.ParseToken(token[0])
		if err != nil {
			logger.Log.Error("failed to parse token", zap.Error(err))
			return nil, status.Error(codes.Unauthenticated, "Access denied")
		}

		logger.Log.Debug(fmt.Sprintf("Got request from User %d, %s", user.ID, user.Login))
		ctx = context.WithValue(ctx, model.UserCtxKey, user)

		return handler(ctx, req)
	}
}

func NewKeeperServer(
	cfg *config.Config,
	iamService *service.IAM,
	secretService *service.SecretService,
) *KeeperServer {
	s := grpc.NewServer(grpc.UnaryInterceptor(iamInterceptor(iamService)))

	m := KeeperServer{
		cfg:           cfg,
		secretService: secretService,
		iam:           iamService,
		srv:           s,
	}
	pb.RegisterKeeperServer(s, &m)

	return &m
}

// Run KeeperServer server. It blocks until the server is stopped.
func (s *KeeperServer) Run(address string) error {
	logger.Log.Info("Run gRPC server", zap.String("Addres", address))
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen gRPC adress at %s. error: %w", address, err)
	}

	return s.srv.Serve(listen)
}

// Shutdown KeeperServer server. It blocks until the server is stopped. Under the hood calls http.Server.Shutdown.
func (s *KeeperServer) Shutdown(ctx context.Context) error {
	logger.Log.Info("Starting gracefull shutdown of gRPC server")
	s.srv.GracefulStop()
	logger.Log.Info("gRPC server is down")
	return nil
}

func (s *KeeperServer) ListSecrets(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	user, ok := ctx.Value(model.UserCtxKey).(*model.User)
	if !ok {
		logger.Log.Error("failed to get user from context")
		return nil, status.Errorf(codes.Unauthenticated, "Access denied")
	}

	log := logger.Log.With(
		zap.Any("request", in),
		zap.Int64("user_id", user.ID),
		zap.String("login", user.Login),
	)
	log.Info("ListSecrets request")

	l, err := s.secretService.ListSecretsMeta(
		ctx,
		&model.SecretListRequest{
			UserID: user.ID,
			Name:   in.Name,
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "ListSecrets error: %s", err)
	}

	secrets := make([]*pb.SecretMeta, 0, len(l.Secrets))
	for _, secret := range l.Secrets {
		st := secretTypeToPbType(secret.Type)
		secrets = append(secrets, &pb.SecretMeta{Id: secret.ID, Name: secret.Name, Type: st})
	}
	response := pb.ListResponse{Secrets: secrets}
	return &response, nil
}

func (s *KeeperServer) GetSecret(ctx context.Context, in *pb.SecretRequest) (*pb.Secret, error) {
	user, ok := ctx.Value(model.UserCtxKey).(*model.User)
	if !ok {
		logger.Log.Error("failed to get user from context")
		return nil, status.Errorf(codes.Unauthenticated, "Access denied")
	}

	log := logger.Log.With(
		zap.Any("request", in),
		zap.Int64("user_id", user.ID),
		zap.String("login", user.Login),
	)
	log.Info("GetSecret request")

	secret, err := s.secretService.GetSecret(
		ctx,
		model.SecretRequest{
			UserID: user.ID,
			ID:     in.Id,
			Type:   pbTypeToSecretType(in.Type),
		},
	)
	if err != nil && errors.Is(err, coreErrors.ErrNotFound404) {
		return nil, status.Errorf(codes.NotFound, "Not found")
	}
	if err != nil {
		log.Error("GetSecret error", zap.Error(err))
		return nil, status.Errorf(codes.Unknown, "Reading Secret error: %s", err)
	}

	response := pb.Secret{
		Id:      secret.ID,
		Name:    secret.Name,
		Type:    secretTypeToPbType(secret.Type),
		Payload: secret.Payload,
	}
	return &response, nil
}

func (s *KeeperServer) CreateSecret(ctx context.Context, in *pb.SecretCreateRequest) (*pb.Secret, error) {
	user, ok := ctx.Value(model.UserCtxKey).(*model.User)
	if !ok {
		logger.Log.Error("failed to get user from context")
		return nil, status.Errorf(codes.Unauthenticated, "Access denied")
	}

	log := logger.Log.With(
		zap.Any("request", in),
		zap.Int64("user_id", user.ID),
		zap.String("login", user.Login),
	)
	log.Info("CreateSecret request")

	secret, err := s.secretService.CreateSecret(
		ctx,
		model.SecretCreateRequest{
			UserID:  user.ID,
			Name:    in.Name,
			Type:    pbTypeToSecretType(in.Type),
			Payload: in.Payload,
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "CreateSecret error: %s", err)
	}

	response := pb.Secret{
		Id:      secret.ID,
		Name:    secret.Name,
		Type:    secretTypeToPbType(secret.Type),
		Payload: secret.Payload,
	}
	return &response, nil
}

func (s *KeeperServer) UpdateSecret(ctx context.Context, in *pb.SecretUpdateRequest) (*pb.Secret, error) {
	user, ok := ctx.Value(model.UserCtxKey).(*model.User)
	if !ok {
		logger.Log.Error("failed to get user from context")
		return nil, status.Errorf(codes.Unauthenticated, "Access denied")
	}

	log := logger.Log.With(
		zap.Any("request", in),
		zap.Int64("user_id", user.ID),
		zap.String("login", user.Login),
	)
	log.Info("UpdateSecret request")

	secret, err := s.secretService.UpdateSecret(
		ctx,
		model.SecretUpdateRequest{
			ID:      in.Id,
			UserID:  user.ID,
			Name:    in.Name,
			Type:    pbTypeToSecretType(in.Type),
			Payload: in.Payload,
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "UpdateSecret error: %s", err)
	}

	response := pb.Secret{
		Id:      secret.ID,
		Name:    secret.Name,
		Type:    secretTypeToPbType(secret.Type),
		Payload: secret.Payload,
	}
	return &response, nil
}
