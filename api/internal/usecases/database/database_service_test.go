package database

import (
	"errors"
	"testing"

	"github.com/Agu-GC/MELI-Challenge/api/internal/domain"
	repoMocks "github.com/Agu-GC/MELI-Challenge/api/internal/repositories/mocks"
	cipherMocks "github.com/Agu-GC/MELI-Challenge/api/internal/usecases/cipher/mocks"
	"github.com/Agu-GC/MELI-Challenge/api/pkg"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseService_CreateDatabase(t *testing.T) {
	mockDBRepo := new(repoMocks.MockDatabaseRepository)
	mockCipher := new(cipherMocks.MockCipherService)
	mockLogger := pkg.NewLogger("test")

	cipherKey := []byte("test-key-32-characters-long!")
	service := &databaseService{
		dbRepo:        mockDBRepo,
		cipherService: mockCipher,
		cipherKey:     cipherKey,
		logger:        mockLogger,
	}

	t.Run("Succed", func(t *testing.T) {
		expectedDB := &domain.Database{
			Host:              "localhost",
			Port:              5432,
			Username:          "user",
			EncryptedPassword: "encrypted-pass",
			Name:              "mydb",
		}
		dbInfo := &DBConnectionInfo{
			Host:     "localhost",
			Port:     5432,
			Username: "user",
			Password: "pass111",
			Name:     "mydb",
		}

		mockCipher.On("Encrypt", "pass111", cipherKey).Return("encrypted-pass", nil)
		mockDBRepo.On("Create", expectedDB).Return(nil)

		result, err := service.CreateDatabase(dbInfo)

		assert.NoError(t, err)
		assert.Equal(t, uint(0), result)

		mockDBRepo.AssertExpectations(t)
		mockCipher.AssertExpectations(t)
	})

	t.Run("Error saving", func(t *testing.T) {
		dbInfo := &DBConnectionInfo{
			Host:     "127.1.1.2",
			Port:     5432,
			Username: "user",
			Password: "pass333",
			Name:     "mydb",
		}
		expectedDB := &domain.Database{
			Host:              "127.1.1.2",
			Port:              5432,
			Username:          "user",
			EncryptedPassword: "encrypted-pass",
			Name:              "mydb",
		}
		mockCipher.On("Encrypt", "pass333", cipherKey).Return("encrypted-pass", nil)
		mockDBRepo.On("Create", expectedDB).Return(errors.New("saving error"))

		result, err := service.CreateDatabase(dbInfo)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "saving error")
		assert.Equal(t, uint(0), result)
		mockDBRepo.AssertExpectations(t)
		mockCipher.AssertExpectations(t)
	})

	t.Run("Error encripting", func(t *testing.T) {
		dbInfo := &DBConnectionInfo{
			Host:     "127.1.1.1",
			Port:     5432,
			Username: "user",
			Password: "pass222",
			Name:     "mydb",
		}
		mockCipher.On("Encrypt", "pass222", cipherKey).Return("", errors.New("error encrypting"))

		result, err := service.CreateDatabase(dbInfo)

		assert.Error(t, err)
		assert.Equal(t, "error encrypting", err.Error())
		assert.Equal(t, uint(0), result)
		mockCipher.AssertExpectations(t)
	})
}

func TestDatabaseService_GetDatabaseByID(t *testing.T) {
	mockDBRepo := new(repoMocks.MockDatabaseRepository)
	mockCipher := new(cipherMocks.MockCipherService)
	mockLogger := pkg.NewLogger("test")

	cipherKey := []byte("test-key-32-characters-long!")
	service := &databaseService{
		dbRepo:        mockDBRepo,
		cipherService: mockCipher,
		cipherKey:     cipherKey,
		logger:        mockLogger,
	}

	t.Run("Succed", func(t *testing.T) {
		expectedDB := &domain.Database{
			Host:              "localhost",
			Port:              5432,
			Username:          "user",
			EncryptedPassword: "encrypted-pass",
			Dialect:           "postgres",
			Name:              "mydb",
		}

		mockDBRepo.On("GetByID", uint(1)).Return(expectedDB, nil)
		mockCipher.On("Decrypt", "encrypted-pass", cipherKey).Return("decrypted-pass", nil)

		result, err := service.GetDatabaseByID(1)

		assert.NoError(t, err)
		assert.Equal(t, &DBConnectionInfo{
			Host:     "localhost",
			Port:     5432,
			Username: "user",
			Password: "decrypted-pass",
			Dialect:  "postgres",
			Name:     "mydb",
		}, result)

		mockDBRepo.AssertExpectations(t)
		mockCipher.AssertExpectations(t)
	})

	t.Run("Database not found", func(t *testing.T) {
		mockDBRepo.On("GetByID", uint(2)).Return(
			&domain.Database{},
			errors.New("record not found"),
		)

		result, err := service.GetDatabaseByID(2)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockDBRepo.AssertExpectations(t)
	})

	t.Run("Decrypt error", func(t *testing.T) {
		expectedDB := &domain.Database{
			EncryptedPassword: "bad-encrypted",
		}

		mockDBRepo.On("GetByID", uint(3)).Return(expectedDB, nil)
		mockCipher.On("Decrypt", "bad-encrypted", cipherKey).Return(
			"",
			errors.New("decryption error"),
		)

		result, err := service.GetDatabaseByID(3)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "decryption error")
		assert.Nil(t, result)
		mockDBRepo.AssertExpectations(t)
		mockCipher.AssertExpectations(t)
	})
}

func TestDatabaseService_GetDatabaseScanResult(t *testing.T) {
	mockDBRepo := new(repoMocks.MockDatabaseRepository)
	mockCipher := new(cipherMocks.MockCipherService)
	mockLogger := pkg.NewLogger("test")

	cipherKey := []byte("test-key-32-characters-long!")
	service := &databaseService{
		dbRepo:        mockDBRepo,
		cipherService: mockCipher,
		cipherKey:     cipherKey,
		logger:        mockLogger,
	}

	t.Run("Succed", func(t *testing.T) {
		expectedScan := DBScanResult{
			DatabaseID:   uint(123),
			DatabaseName: "mydb",
			Host:         "localhost",
			Port:         5432,
			Tables: []TableScanned{
				{TableName: "user", Columns: []ColumnScanned{{ColumnName: "user_email", DataType: "EMAIL"}}},
			},
		}
		expectedDB := &domain.Database{
			Base:              domain.Base{ID: uint(123)},
			Host:              "localhost",
			Port:              5432,
			Username:          "user",
			EncryptedPassword: "encrypted-pass",
			Dialect:           "postgres",
			Name:              "mydb",
			Scans: []domain.Scan{
				{
					DatabaseID: uint(123),
					Status:     "completed",
					Tables: []domain.Table{
						{
							ScanID:    11,
							TableName: "user",
							Columns: []domain.Column{
								{
									ColumnName:     "user_email",
									Classification: domain.Classification{Name: "EMAIL"},
								},
							},
						},
					},
				},
			},
		}
		mockDBRepo.On("GetWithLastScanInfo", uint(1)).Return(expectedDB, nil)

		result, err := service.GetDatabaseScanResult(1)

		assert.Nil(t, err)
		assert.Equal(t, &expectedScan, result)
		assert.Nil(t, err)
		mockDBRepo.AssertExpectations(t)
	})

	t.Run("no scans found", func(t *testing.T) {
		expectedDB := &domain.Database{
			Base:              domain.Base{ID: uint(123)},
			Host:              "localhost",
			Port:              5432,
			Username:          "user",
			EncryptedPassword: "encrypted-pass",
			Dialect:           "postgres",
			Name:              "mydb",
			Scans:             []domain.Scan{},
		}

		mockDBRepo.On("GetWithLastScanInfo", uint(2)).Return(expectedDB, nil)

		result, err := service.GetDatabaseScanResult(2)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "no scans found", err.Error())
		mockDBRepo.AssertExpectations(t)
	})

	t.Run("database not found", func(t *testing.T) {
		mockDBRepo.On("GetWithLastScanInfo", uint(3)).Return(&domain.Database{}, errors.New("record not found"))

		result, err := service.GetDatabaseScanResult(3)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "record not found", err.Error())
		mockDBRepo.AssertExpectations(t)
	})
}
