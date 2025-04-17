package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"github.com/Agu-GC/MELI-Challenge/api/pkg"
)

type CipherServiceInterface interface {
	Encrypt(plaintext string, key []byte) (string, error)
	Decrypt(plaintext string, key []byte) (string, error)
}

type cipherService struct {
	logger pkg.Logger
}

func NewCipherService(logger pkg.Logger) CipherServiceInterface {
	return &cipherService{logger: logger}
}

func (cs *cipherService) Encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		cs.logger.Error("Error creating the cipher", map[string]any{"error": err})
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		cs.logger.Error("Error creating the GCM", map[string]any{"error": err})
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		cs.logger.Error("Error ReadFull", map[string]any{"error": err})
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (cs *cipherService) Decrypt(encodedCiphertext string, key []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		cs.logger.Error("Error decoding the string", map[string]any{"error": err})
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		cs.logger.Error("Error creating the cipher", map[string]any{"error": err})
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		cs.logger.Error("Error creating the GCM", map[string]any{"error": err})
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		cs.logger.Error("ciphertext too short", map[string]any{"error": err})
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		cs.logger.Error("error last step decripting", map[string]any{"error": err})
		return "", err
	}

	return string(plaintext), nil
}
