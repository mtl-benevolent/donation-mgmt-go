package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

func parseKey(keyHex string) ([]byte, error) {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return key, fmt.Errorf("could not decode hex key: %w", err)
	}

	if len(key) != 32 {
		return []byte{}, fmt.Errorf("Encryption Key must be 32 bytes long")
	}

	return key, nil
}

func EncryptString(str string, keyHex string) (string, error) {
	key, err := parseKey(keyHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("could not create cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("could not create GCM: %w", err)
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("could not generate nonce: %w", err)
	}

	// Encrypt the data. The nonce is prepended to the ciphertext.
	plaintext := []byte(string(str))
	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)
	base64Ciphertext := base64.StdEncoding.EncodeToString(ciphertext)

	return base64Ciphertext, nil
}

func DecryptString(encrypted string, keyHex string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("could not decode base64 ciphertext: %w", err)
	}

	key, err := parseKey(keyHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("could not create AES cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("could not create GCM: %w", err)
	}

	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return "", fmt.Errorf("could not decrypt data: %w", err)
	}

	return string(plaintext), nil
}
