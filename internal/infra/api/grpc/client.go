package grpc

import (
	"context"
	"fmt"
	"log"

	"keeper/internal/core/model"
	pb "keeper/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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

func (c *Client) setToken(token string) {
	log.Println("Call SetToken method")
	c.ctx = metadata.NewOutgoingContext(c.ctx, metadata.Pairs("token", string(token)))
}

func (c *Client) Login(req model.UserRequest) error {
	log.Println("Call Login method", "Login", req.Login)
	token, err := c.client.Login(c.ctx, &pb.LoginRequest{Login: req.Login, Password: req.Password})
	if err != nil {
		return fmt.Errorf("grpc - login error: %w", err)
	}
	c.setToken(token.Token)
	return nil
}

func (c *Client) SignUp(req model.UserRequest) error {
	token, err := c.client.SignUp(c.ctx, &pb.SignUpRequest{Login: req.Login, Password: req.Password})
	if err != nil {
		return fmt.Errorf("grpc - signup update error: %w", err)
	}

	c.setToken(token.Token)
	return nil
}

func (c *Client) ListSecrets(secretName string) (*model.SecretList, error) {
	l, err := c.client.ListSecrets(c.ctx, &pb.ListRequest{Name: secretName})
	if err != nil {
		return nil, fmt.Errorf("grpc - get secret list error: %w", err)
	}
	var resp model.SecretList
	for _, s := range l.Secrets {
		t := pbTypeToSecretType(s.Type)
		resp.Secrets = append(resp.Secrets, &model.SecretMeta{ID: s.Id, Name: s.Name, Type: t, Note: s.Note})
	}
	return &resp, nil
}

func (c *Client) GetSecret(secretID int64) (*model.Secret, error) {
	r, err := c.client.GetSecret(c.ctx, &pb.SecretRequest{Id: secretID})
	if err != nil {
		return nil, fmt.Errorf("grpc - get secret (id: %d) error: %w", secretID, err)
	}
	t := pbTypeToSecretType(r.Type)
	secret := model.NewSecret(r.Id, r.Name, t, r.Payload, r.Note)
	return secret, nil
}

func (c *Client) CreateSecret(secretType model.SecretType, name string, note string, payload []byte) (*model.Secret, error) {
	r, err := c.client.CreateSecret(
		c.ctx,
		&pb.SecretCreateRequest{
			Type:    secretTypeToPbType(secretType),
			Name:    name,
			Note:    note,
			Payload: payload,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("grpc - create secret error: %w", err)
	}
	t := pbTypeToSecretType(r.Type)
	secret := model.NewSecret(r.Id, r.Name, t, r.Payload, r.Note)
	return secret, nil
}

func (c *Client) CreateFileSecret(name, fileName, note string, payload []byte) (int64, error) {
	r, err := c.client.CreateFile(
		c.ctx,
		&pb.SecretFileCreateRequest{
			Name:     name,
			FileName: fileName,
			Payload:  payload,
			Note:     note,
		},
	)
	if err != nil {
		return 0, fmt.Errorf("grpc - create secret error: %w", err)
	}
	return r.Id, nil
}

func (c *Client) UpdateSecret(id int64, secretType model.SecretType, name string, note string, payload []byte) (*model.Secret, error) {
	r, err := c.client.UpdateSecret(
		c.ctx,
		&pb.SecretUpdateRequest{
			Id:      id,
			Type:    secretTypeToPbType(secretType),
			Name:    name,
			Note:    note,
			Payload: payload,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("grpc - update secret error: %w", err)
	}
	t := pbTypeToSecretType(r.Type)
	secret := model.NewSecret(r.Id, r.Name, t, r.Payload, r.Note)
	return secret, nil
}
