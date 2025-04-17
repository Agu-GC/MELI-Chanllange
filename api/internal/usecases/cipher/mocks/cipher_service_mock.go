package mocks

import "github.com/stretchr/testify/mock"

type MockCipherService struct {
	mock.Mock
}

func (m *MockCipherService) Encrypt(plaintext string, key []byte) (string, error) {
	args := m.Called(plaintext, key)
	return args.String(0), args.Error(1)
}

func (m *MockCipherService) Decrypt(encodedCiphertext string, key []byte) (string, error) {
	args := m.Called(encodedCiphertext, key)
	return args.String(0), args.Error(1)
}
