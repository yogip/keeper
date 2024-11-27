package grpc

import (
	"context"
	"fmt"
	"log"

	"keeper/internal/core/model"
	pb "keeper/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	ctx    context.Context
	conn   *grpc.ClientConn
	client pb.KeeperClient
}

func NewClient(ctx context.Context, serverHost string) *Client {
	log.Println("Init grpc client", "serverHost", serverHost)
	conn, err := grpc.NewClient(serverHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewKeeperClient(conn)

	return &Client{
		ctx:    ctx,
		conn:   conn,
		client: client,
	}
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Login(req model.UserRequest) (model.Token, error) {
	log.Println("Call Login method", "Login", req.Login)
	token, err := c.client.Login(c.ctx, &pb.LoginRequest{Login: req.Login, Password: req.Password})
	if err != nil {
		return "", fmt.Errorf("grpc - login error: %w", err)
	}

	return model.Token(token.Token), nil
}

func (c *Client) SignUp(req model.UserRequest) (model.Token, error) {
	token, err := c.client.SignUp(c.ctx, &pb.SignUpRequest{Login: req.Login, Password: req.Password})
	if err != nil {
		return "", fmt.Errorf("grpc - signup update error: %w", err)
	}

	return model.Token(token.Token), nil
}

func (c *Client) GetPassword(metricID int64) (*model.Password, error) {
	pwd, err := c.client.GetPassword(c.ctx, &pb.PasswordRequest{Id: metricID})
	if err != nil {
		return nil, fmt.Errorf("grpc - get secret error: %w", err)
	}
	resp := &model.Password{
		SecretMeta: model.SecretMeta{ID: pwd.Id, Name: pwd.Name},
		Password:   pwd.Password,
	}
	return resp, nil
}
