package uuid

import (
	"github.com/aria3ppp/rag-server/internal/vectorstore/usecase"

	"github.com/google/uuid"
)

type uuidIDGenerator struct{}

var _ usecase.IDGenerator = (*uuidIDGenerator)(nil)

func NewIDGenerator() *uuidIDGenerator {
	return &uuidIDGenerator{}
}

func (*uuidIDGenerator) NewID() (string, error) {
	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return randomUUID.String(), nil
}
