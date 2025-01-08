package container

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type (
	MinIOContainer struct {
		testcontainers.Container
		//add Config
		Config MinIOContainerConfig
	}
	//also add options pattern method
	MinIOContainerOption func(c *MinIOContainerConfig)

	MinIOContainerConfig struct {
		ImageTag   string
		User       string
		Password   string
		MappedPort string
		BucketName string
		Host       string
	}
)

func (c MinIOContainer) GetDSN() string {
	return fmt.Sprintf("%s:%s", c.Config.Host, c.Config.MappedPort)
}

func NewMinIOContainer(ctx context.Context, opts ...MinIOContainerOption) (*MinIOContainer, error) {
	const (
		image         = "bitnami/minio"
		containerPort = "9000"
	)

	config := MinIOContainerConfig{
		ImageTag:   "2024.11.7",
		User:       "admin",
		Password:   "password",
		BucketName: "bucket",
	}
	//handle possible options
	for _, opt := range opts {
		opt(&config)
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Env: map[string]string{
				"MINIO_ROOT_USER":       config.User,
				"MINIO_ROOT_PASSWORD":   config.Password,
				"MINIO_USE_SSL":         "0",
				"MINIO_DEFAULT_BUCKETS": config.BucketName,
			},
			ExposedPorts: []string{
				containerPort,
			},
			Image:      fmt.Sprintf("%s:%s", image, config.ImageTag),
			WaitingFor: wait.ForListeningPort(nat.Port(containerPort)),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("getting request provider: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting host for: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(containerPort))
	if err != nil {
		return nil, fmt.Errorf("getting mapped port for (%s): %w", containerPort, err)
	}
	config.MappedPort = mappedPort.Port()
	config.Host = host

	fmt.Println("Host:", config.Host, config.MappedPort)

	return &MinIOContainer{
		Container: container,
		Config:    config,
	}, nil
}
