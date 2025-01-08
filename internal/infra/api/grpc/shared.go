package grpc

import (
	"keeper/internal/core/model"
	pb "keeper/internal/proto"
)

func secretTypeToPbType(t model.SecretType) pb.SecretType {
	var st pb.SecretType
	switch t {
	case model.SecretTypePassword:
		st = pb.SecretType_PASSWORD
	case model.SecretTypeCard:
		st = pb.SecretType_CARD
	case model.SecretTypeFile:
		st = pb.SecretType_FILE
	case model.SecretTypeNote:
		st = pb.SecretType_NOTE
	}
	return st
}

func pbTypeToSecretType(t pb.SecretType) model.SecretType {
	var st model.SecretType
	switch t {
	case pb.SecretType_PASSWORD:
		st = model.SecretTypePassword
	case pb.SecretType_CARD:
		st = model.SecretTypeCard
	case pb.SecretType_FILE:
		st = model.SecretTypeFile
	case pb.SecretType_NOTE:
		st = model.SecretTypeNote
	}
	return st
}
