package mocks

import (
	"github.com/Agu-GC/MELI-Challenge/api/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockClassificationRepository struct {
	mock.Mock
}

func (m *MockClassificationRepository) Create(classification *domain.Classification) error {
	args := m.Called(classification)
	return args.Error(0)
}

func (m *MockClassificationRepository) GetAll() ([]*domain.Classification, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Classification), args.Error(1)
}
