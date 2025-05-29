package encryption_test

import (
	"context"
	"testing"

	"github.com/android-sms-gateway/twilio-fallback/internal/encryption"
	"go.uber.org/fx"
)

func TestNewEncryptor(t *testing.T) {
	tests := []struct {
		name    string
		key     []byte
		wantErr bool
	}{
		{
			name:    "key length less than 32 bytes",
			key:     []byte("shortkey"),
			wantErr: true,
		},
		{
			name:    "key length equal to 32 bytes",
			key:     []byte("thisisaverysecure32bytekey1234567890"),
			wantErr: false,
		},
		{
			name:    "key length greater than 32 bytes",
			key:     []byte("thisisaverysecure32bytekey12345678901234567890"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := encryption.NewEncryptor(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEncryptor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncryptDecrypt(t *testing.T) {
	// Test valid encryption/decryption
	t.Run("Valid encryption/decryption", func(t *testing.T) {
		key := "thisisaverysecure32bytekey1234567890"
		text := "sensitive data"

		encryptor, err := encryption.NewEncryptor([]byte(key))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		encrypted, err := encryptor.Encrypt(text)
		if err != nil {
			t.Fatalf("Encryption failed: %v", err)
		}

		decrypted, err := encryptor.Decrypt(encrypted)
		if err != nil {
			t.Fatalf("Decryption failed: %v", err)
		}

		if decrypted != text {
			t.Fatalf("Decrypted text doesn't match original. Expected: %s, Got: %s", text, decrypted)
		}
	})

	// Test invalid ciphertext
	t.Run("Invalid ciphertext", func(t *testing.T) {
		key := "thisisaverysecure32bytekey1234567890"

		encryptor, err := encryption.NewEncryptor([]byte(key))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		_, err = encryptor.Decrypt("invalidciphertext")
		if err == nil {
			t.Fatal("Expected error when decrypting invalid ciphertext")
		}
	})
}

func TestFxModule(t *testing.T) {
	app := fx.New(
		fx.Provide(
			func() encryption.Config {
				return encryption.Config{
					Key: "thisisaverysecure32bytekey1234567890",
				}
			},
		),
		encryption.Module,
	)

	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("Failed to start app: %v", err)
	}

	if err := app.Stop(context.Background()); err != nil {
		t.Fatalf("Failed to stop app: %v", err)
	}
}
