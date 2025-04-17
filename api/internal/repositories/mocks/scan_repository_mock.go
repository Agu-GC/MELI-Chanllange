package mocks

import (
	"github.com/Agu-GC/MELI-Challenge/api/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockScanRepository struct {
	mock.Mock
}

func (m *MockScanRepository) Create(scan *domain.Scan) error {
	args := m.Called(scan)
	return args.Error(0)
}

func (m *MockScanRepository) GetByID(id uint) (*domain.Scan, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Scan), args.Error(1)
}
