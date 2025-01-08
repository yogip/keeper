package service

import (
	"context"
	"errors"
	"fmt"
	"keeper/internal/core/config"
	"keeper/internal/core/model"
	"keeper/internal/infra/repo"
	"keeper/internal/logger"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID    int64
	UserLogin string
}

type IAM struct {
	cfg      *config.IAMConfig
	userRepo *repo.UserRepo
}

func NewIAMService(userRepo *repo.UserRepo, cfg *config.IAMConfig) *IAM {
	return &IAM{userRepo: userRepo, cfg: cfg}
}

// check login and password pair and generate token if password hash match
func (iam *IAM) Login(ctx context.Context, user *model.UserRequest) (string, error) {
	u, err := iam.userRepo.GetUser(ctx, user.Login)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", errors.New("")
	}

	err = bcrypt.CompareHashAndPassword(*u.PasswordHash, []byte(user.Password))
	if err != nil {
		return "", err
	}

	token, err := iam.buildToken(u)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Create new user and auth them (return token)
func (iam *IAM) Register(ctx context.Context, user *model.UserRequest) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hasing password error: %w", err)
	}

	u, err := iam.userRepo.CreateUser(ctx, user.Login, hashedPassword)
	if err != nil {
		return "", err
	}

	token, err := iam.buildToken(u)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (iam *IAM) buildToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		TokenClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(iam.cfg.TokenTTL)),
			},
			UserID:    user.ID,
			UserLogin: user.Login,
		},
	)

	tokenString, err := token.SignedString([]byte(iam.cfg.SecretKey))
	if err != nil {
		return "", fmt.Errorf("build token error: %w", err)
	}

	logger.Log.Info(fmt.Sprintf("buildToken token: %s", tokenString))
	return tokenString, nil
}

func (iam *IAM) ParseToken(rawToken string) (*model.User, error) {
	claims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(rawToken, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(iam.cfg.SecretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("parse token error: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("tokne is not valid")
	}

	return &model.User{ID: claims.UserID, Login: claims.UserLogin}, nil
}
