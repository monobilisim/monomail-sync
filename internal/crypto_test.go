package internal

import (
	"os"
	"strings"
	"testing"
)

func setTestEncryptionKey(t *testing.T) {
	t.Helper()
	os.Setenv("MONOMAIL_SYNC_ENCRYPTION_KEY", "01234567890123456789012345678901")
	t.Cleanup(func() { os.Unsetenv("MONOMAIL_SYNC_ENCRYPTION_KEY") })
}

func TestEncryptDecrypt_Roundtrip(t *testing.T) {
	setTestEncryptionKey(t)

	tests := []string{
		"hello",
		"password123!@#",
		"こんにちは世界",
		"a",
		strings.Repeat("x", 10000),
	}

	for _, plaintext := range tests {
		encrypted, err := EncryptString(plaintext)
		if err != nil {
			t.Fatalf("EncryptString(%q) error: %v", plaintext, err)
		}

		if !strings.HasPrefix(encrypted, encryptionPrefix) {
			t.Errorf("encrypted value should start with %q", encryptionPrefix)
		}

		if encrypted == plaintext {
			t.Error("encrypted value should differ from plaintext")
		}

		decrypted, err := DecryptString(encrypted)
		if err != nil {
			t.Fatalf("DecryptString error: %v", err)
		}

		if decrypted != plaintext {
			t.Errorf("DecryptString = %q, want %q", decrypted, plaintext)
		}
	}
}

func TestEncryptDecrypt_EmptyString(t *testing.T) {
	setTestEncryptionKey(t)

	encrypted, err := EncryptString("")
	if err != nil {
		t.Fatalf("EncryptString('') error: %v", err)
	}
	if encrypted != "" {
		t.Errorf("EncryptString('') = %q, want ''", encrypted)
	}

	decrypted, err := DecryptString("")
	if err != nil {
		t.Fatalf("DecryptString('') error: %v", err)
	}
	if decrypted != "" {
		t.Errorf("DecryptString('') = %q, want ''", decrypted)
	}
}

func TestEncrypt_UniqueNonce(t *testing.T) {
	setTestEncryptionKey(t)

	enc1, _ := EncryptString("same-input")
	enc2, _ := EncryptString("same-input")

	if enc1 == enc2 {
		t.Error("two encryptions of same plaintext should produce different ciphertexts (random nonce)")
	}
}

func TestDecrypt_NotEncrypted(t *testing.T) {
	setTestEncryptionKey(t)

	_, err := DecryptString("plaintext-without-prefix")
	if err != errValueNotEncrypted {
		t.Errorf("expected errValueNotEncrypted, got: %v", err)
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	setTestEncryptionKey(t)

	_, err := DecryptString(encryptionPrefix + "not-valid-base64!!!")
	if err == nil {
		t.Error("expected error for invalid base64")
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	setTestEncryptionKey(t)

	encrypted, _ := EncryptString("secret")

	tampered := encrypted[:len(encrypted)-2] + "XX"
	_, err := DecryptString(tampered)
	if err == nil {
		t.Error("expected error for tampered ciphertext")
	}
}

func TestEncrypt_MissingKey(t *testing.T) {
	os.Unsetenv("MONOMAIL_SYNC_ENCRYPTION_KEY")

	_, err := EncryptString("test")
	if err != errEncryptionKeyMissing {
		t.Errorf("expected errEncryptionKeyMissing, got: %v", err)
	}
}

func TestGetEncryptionKey_Base64(t *testing.T) {
	os.Setenv("MONOMAIL_SYNC_ENCRYPTION_KEY", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=")
	t.Cleanup(func() { os.Unsetenv("MONOMAIL_SYNC_ENCRYPTION_KEY") })

	key, err := getEncryptionKey()
	if err != nil {
		t.Fatalf("getEncryptionKey error: %v", err)
	}
	if len(key) != 32 {
		t.Errorf("key length = %d, want 32", len(key))
	}
}

func TestGetEncryptionKey_InvalidLength(t *testing.T) {
	os.Setenv("MONOMAIL_SYNC_ENCRYPTION_KEY", "tooshort")
	t.Cleanup(func() { os.Unsetenv("MONOMAIL_SYNC_ENCRYPTION_KEY") })

	_, err := getEncryptionKey()
	if err == nil {
		t.Error("expected error for short non-base64 key")
	}
}

func TestEncryptDecrypt_WithBase64Key(t *testing.T) {
	os.Setenv("MONOMAIL_SYNC_ENCRYPTION_KEY", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=")
	t.Cleanup(func() { os.Unsetenv("MONOMAIL_SYNC_ENCRYPTION_KEY") })

	encrypted, err := EncryptString("base64-key-test")
	if err != nil {
		t.Fatalf("EncryptString error: %v", err)
	}

	decrypted, err := DecryptString(encrypted)
	if err != nil {
		t.Fatalf("DecryptString error: %v", err)
	}

	if decrypted != "base64-key-test" {
		t.Errorf("roundtrip failed: got %q", decrypted)
	}
}
