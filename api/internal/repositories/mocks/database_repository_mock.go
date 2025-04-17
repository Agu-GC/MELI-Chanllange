package mocks

import (
	"github.com/Agu-GC/MELI-Challenge/api/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockDatabaseRepository struct {
	mock.Mock
}

func (m *MockDatabaseRepository) Create(database *domain.Database) error {
	args := m.Called(database)
	return args.Error(0)
}

func (m *MockDatabaseRepository) GetByID(id uint) (*domain.Database, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Database), args.Error(1)
}

func (m *MockDatabaseRepository) GetWithLastScanInfo(id uint) (*domain.Database, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Database), args.Error(1)
}
