package crypto

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestEncryptDecryptRoundtrip(t *testing.T) {
	testCases := []struct {
		name       string
		data       []byte
		passphrase string
	}{
		{
			name:       "simple text",
			data:       []byte("Hello, World!"),
			passphrase: "test-passphrase",
		},
		{
			name:       "empty data",
			data:       []byte(""),
			passphrase: "test-passphrase",
		},
		{
			name:       "binary data",
			data:       []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD},
			passphrase: "test-passphrase",
		},
		{
			name:       "long passphrase",
			data:       []byte("test data"),
			passphrase: "this-is-a-very-long-passphrase-for-testing-purposes-123456789",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := Encrypt(tc.data, tc.passphrase)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			// Decrypt
			decrypted, err := Decrypt(encrypted, tc.passphrase)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			// Compare
			if !bytes.Equal(tc.data, decrypted) {
				t.Errorf("Decrypted data does not match original.\nOriginal:  %q\nDecrypted: %q", tc.data, decrypted)
			}
		})
	}
}

func TestDecryptWithWrongPassphrase(t *testing.T) {
	data := []byte("secret message")
	passphrase := "correct-passphrase"

	encrypted, err := Encrypt(data, passphrase)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Try to decrypt with wrong passphrase
	_, err = Decrypt(encrypted, "wrong-passphrase")
	if err == nil {
		t.Error("Expected error when decrypting with wrong passphrase, got nil")
	}
}

func TestOpenSSLCompatibility(t *testing.T) {
	// Skip if openssl is not available
	if _, err := exec.LookPath("openssl"); err != nil {
		t.Skip("openssl not available, skipping compatibility test")
	}

	testData := "This is a test message for OpenSSL compatibility!"
	passphrase := "test-passphrase-123"

	// Test 1: Encrypt with Go, decrypt with OpenSSL
	t.Run("Go encrypt, OpenSSL decrypt", func(t *testing.T) {
		// Encrypt with Go
		encrypted, err := Encrypt([]byte(testData), passphrase)
		if err != nil {
			t.Fatalf("Encrypt failed: %v", err)
		}

		// Write encrypted data to temp file
		tmpEncrypted := filepath.Join(t.TempDir(), "encrypted.bin")
		if err := os.WriteFile(tmpEncrypted, encrypted, 0644); err != nil {
			t.Fatalf("Failed to write encrypted file: %v", err)
		}

		// Decrypt with OpenSSL
		cmd := exec.Command("openssl", "enc", "-aes-256-cbc", "-d", "-pbkdf2",
			"-in", tmpEncrypted, "-pass", "pass:"+passphrase)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("OpenSSL decrypt failed: %v\nOutput: %s", err, output)
		}

		// Compare
		if string(output) != testData {
			t.Errorf("OpenSSL decrypted data does not match.\nExpected: %q\nGot:      %q", testData, string(output))
		}
	})

	// Test 2: Encrypt with OpenSSL, decrypt with Go
	t.Run("OpenSSL encrypt, Go decrypt", func(t *testing.T) {
		// Create temp file with test data
		tmpPlain := filepath.Join(t.TempDir(), "plain.txt")
		if err := os.WriteFile(tmpPlain, []byte(testData), 0644); err != nil {
			t.Fatalf("Failed to write plain file: %v", err)
		}

		// Encrypt with OpenSSL
		tmpEncrypted := filepath.Join(t.TempDir(), "encrypted.bin")
		cmd := exec.Command("openssl", "enc", "-aes-256-cbc", "-salt", "-pbkdf2",
			"-in", tmpPlain, "-out", tmpEncrypted, "-pass", "pass:"+passphrase)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("OpenSSL encrypt failed: %v\nOutput: %s", err, output)
		}

		// Read encrypted data
		encrypted, err := os.ReadFile(tmpEncrypted)
		if err != nil {
			t.Fatalf("Failed to read encrypted file: %v", err)
		}

		// Decrypt with Go
		decrypted, err := Decrypt(encrypted, passphrase)
		if err != nil {
			t.Fatalf("Decrypt failed: %v", err)
		}

		// Compare
		if string(decrypted) != testData {
			t.Errorf("Go decrypted data does not match.\nExpected: %q\nGot:      %q", testData, string(decrypted))
		}
	})
}

func TestBase64Encoding(t *testing.T) {
	testData := []byte("test data")

	encoded := EncodeBase64(testData)
	decoded, err := DecodeBase64(encoded)
	if err != nil {
		t.Fatalf("DecodeBase64 failed: %v", err)
	}

	if !bytes.Equal(testData, decoded) {
		t.Errorf("Decoded data does not match.\nOriginal: %q\nDecoded:  %q", testData, decoded)
	}
}

func TestTarballRoundtrip(t *testing.T) {
	// Create temp directory with test files
	tmpDir := t.TempDir()
	baseDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		t.Fatalf("Failed to create base dir: %v", err)
	}

	// Create test files
	files := map[string]string{
		"file1.txt":     "content of file 1",
		"dir/file2.txt": "content of file 2",
	}

	for path, content := range files {
		fullPath := filepath.Join(baseDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", path, err)
		}
	}

	// Create tarball
	var filePaths []string
	for path := range files {
		filePaths = append(filePaths, path)
	}

	tarData, err := CreateTarball(filePaths, baseDir)
	if err != nil {
		t.Fatalf("CreateTarball failed: %v", err)
	}

	// Extract to new directory
	extractDir := filepath.Join(tmpDir, "extract")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		t.Fatalf("Failed to create extract dir: %v", err)
	}

	extracted, err := ExtractTarball(tarData, extractDir)
	if err != nil {
		t.Fatalf("ExtractTarball failed: %v", err)
	}

	// Verify extracted files
	if len(extracted) != len(files) {
		t.Errorf("Expected %d extracted files, got %d", len(files), len(extracted))
	}

	for path, expectedContent := range files {
		fullPath := filepath.Join(extractDir, path)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			t.Errorf("Failed to read extracted file %s: %v", path, err)
			continue
		}

		if string(content) != expectedContent {
			t.Errorf("File %s content mismatch.\nExpected: %q\nGot:      %q", path, expectedContent, string(content))
		}
	}
}
