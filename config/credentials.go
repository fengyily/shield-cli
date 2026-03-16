package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// Credentials represents connector credentials
type Credentials struct {
	ConnectorName string `json:"connector_name"`
	Password      string `json:"password"`
}

// GetCredentialFilePath returns the path to the credential file
func GetCredentialFilePath() string {
	var dir string
	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("LOCALAPPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
		}
		dir = filepath.Join(dir, "ShieldCLI")
	default:
		home, _ := os.UserHomeDir()
		dir = filepath.Join(home, ".shield-cli")
	}
	return filepath.Join(dir, ".credential")
}

// EncryptAndSave encrypts and saves credentials using machine fingerprint
func (c *Credentials) EncryptAndSave(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	key, err := getDerivedKey()
	if err != nil {
		return fmt.Errorf("failed to get encryption key: %w", err)
	}

	plaintext, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to serialize credentials: %w", err)
	}

	encrypted, err := encryptAESGCM(plaintext, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt credentials: %w", err)
	}

	if err := os.WriteFile(path, encrypted, 0600); err != nil {
		return fmt.Errorf("failed to write credentials file: %w", err)
	}

	return nil
}

// LoadCredentials loads and decrypts credentials
func LoadCredentials(path string) (*Credentials, error) {
	encrypted, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	key, err := getDerivedKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get decryption key: %w", err)
	}

	plaintext, err := decryptAESGCM(encrypted, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt credentials (machine fingerprint may not match): %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(plaintext, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}

	return &creds, nil
}

// GetOrCreateCredentials loads existing credentials or generates new ones
func GetOrCreateCredentials() (*Credentials, error) {
	path := GetCredentialFilePath()

	// Try to load existing credentials
	if creds, err := LoadCredentials(path); err == nil {
		return creds, nil
	}

	// Generate new credentials using machine fingerprint as connector_name
	fingerprint, err := GetMachineFingerprint()
	if err != nil {
		return nil, fmt.Errorf("failed to get machine fingerprint: %w", err)
	}

	// Use first 12 chars of fingerprint as connector name
	name := fingerprint
	if len(name) > 12 {
		name = name[:12]
	}
	name = "shield-" + name

	return &Credentials{
		ConnectorName: name,
		Password:      "",
	}, nil
}

func getDerivedKey() ([]byte, error) {
	fingerprint, err := GetMachineFingerprint()
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256([]byte(fingerprint))
	return hash[:], nil
}

func encryptAESGCM(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decryptAESGCM(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce := ciphertext[:nonceSize]
	ciphertextBytes := ciphertext[nonceSize:]

	return gcm.Open(nil, nonce, ciphertextBytes, nil)
}
