package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

const encryptionPrefix = "enc:"

var errEncryptionKeyMissing = errors.New("encryption key missing")
var errValueNotEncrypted = errors.New("value is not encrypted")

func getEncryptionKey() ([]byte, error) {
	key := os.Getenv("MONOMAIL_SYNC_ENCRYPTION_KEY")
	if key == "" {
		return nil, errEncryptionKeyMissing
	}

	if len(key) == 32 {
		return []byte(key), nil
	}

	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, errors.New("invalid encryption key format")
	}
	if len(decoded) != 32 {
		return nil, errors.New("invalid encryption key length")
	}

	return decoded, nil
}

func EncryptString(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return encryptionPrefix + encoded, nil
}

func DecryptString(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	if len(ciphertext) < len(encryptionPrefix) || ciphertext[:len(encryptionPrefix)] != encryptionPrefix {
		return "", errValueNotEncrypted
	}

	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	encoded := ciphertext[len(encryptionPrefix):]
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(raw) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, data := raw[:nonceSize], raw[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
