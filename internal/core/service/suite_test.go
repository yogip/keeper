package service

import (
	"context"
	"database/sql"
	"fmt"
	"keeper/internal/core/config"
	"keeper/internal/core/encryption"
	"keeper/internal/helpers/container"
	"keeper/internal/infra/repo"
	"keeper/internal/infra/s3"
	"keeper/internal/logger"
	"keeper/migrations"
	"os"
	"testing"
	"time"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	psqlContainer  *container.PostgreSQLContainer
	minIOContainer *container.MinIOContainer
	encrypter      *encryption.EncryptionService
	secretSrv      *SecretService
	db             *sql.DB
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	//---------------- common part for all tests
	privKeyData := `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`

	ctx, ctxCancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer ctxCancel()

	// Prepare test contaner
	logger.Log.Debug("Creating PostgreSQL container")
	psqlContainer, err := container.NewPostgreSQLContainer(ctx)
	s.Require().NoError(err)

	logger.Log.Debug("Creating MinIO container")
	s3Container, err := container.NewMinIOContainer(ctx)
	s.Require().NoError(err)

	logger.Log.Debug("Run DB migration")
	err = migrations.RunMigration(ctx, psqlContainer.GetDSN())
	s.Require().NoError(err)

	db, err := sql.Open("pgx", psqlContainer.GetDSN())
	s.Require().NoError(err)

	// Load fixtures
	logger.Log.Debug("Load fixtures")
	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("../../../fixtures/storage"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())

	// Prepare repo and services
	repoSecret := repo.NewSecretRepo(db)

	masterKey, err := encryption.NewPrivateKey()
	s.Require().NoError(err)

	encPrivateKey, err := masterKey.Encrypt([]byte(privKeyData))
	s.Require().NoError(err)

	f, err := os.OpenFile("/tmp/encryption_key_v_1.pem", os.O_CREATE|os.O_WRONLY, 0777)
	s.Require().NoError(err)

	f.Write(encPrivateKey)
	f.Close()

	s3Client, err := s3.NewS3Client(
		ctx,
		&config.S3Config{
			BucketName:      s3Container.Config.BucketName,
			Endpoint:        s3Container.GetDSN(),
			AccessKeyID:     s3Container.Config.User,
			SecretAccessKey: s3Container.Config.Password,
		},
	)
	s.Require().NoError(err)

	s.db = db
	s.psqlContainer = psqlContainer
	s.minIOContainer = s3Container
	s.encrypter = encryption.NewEncryptionService("/tmp", masterKey)
	s.secretSrv = NewSecretService(repoSecret, s3Client, s.encrypter, 1)
	logger.Log.Debug("Setup done")
}

// Helper to construct fixtures secret part.
func (s *TestSuite) encript(plaitext string) {
	ch, key, _ := s.encrypter.Encrypt([]byte(plaitext), 1)
	fmt.Println("-------------")
	fmt.Println(ch)
	fmt.Println("-")
	fmt.Println(key)
	fmt.Println("-------------")
}

func (s *TestSuite) TearDownSuite() {
	logger.Log.Debug("Start to tear down")
	ctx, ctxCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer ctxCancel()

	s.db.Close()

	err := s.psqlContainer.Terminate(ctx)
	s.Require().NoError(err)

	err = s.minIOContainer.Terminate(ctx)
	s.Require().NoError(err)

	os.Remove("/tmp/encryption_key_v_1.pem")
}
