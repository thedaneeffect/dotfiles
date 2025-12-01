package crypto

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// OpenSSL magic bytes
	opensslSaltHeader = "Salted__"
	saltLength        = 8
	pbkdf2Iterations  = 10000 // Match openssl default
	aesKeySize        = 32    // AES-256
	aesIVSize         = 16
)

// CreateTarball creates a gzipped tarball from the given files
func CreateTarball(files []string, baseDir string) ([]byte, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	for _, file := range files {
		fullPath := filepath.Join(baseDir, file)

		// Get file info
		info, err := os.Stat(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to stat %s: %w", file, err)
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return nil, fmt.Errorf("failed to create header for %s: %w", file, err)
		}
		header.Name = file // Use relative path

		// Write header
		if err := tw.WriteHeader(header);err != nil {
			return nil, fmt.Errorf("failed to write header for %s: %w", file, err)
		}

		// Write file contents
		f, err := os.Open(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open %s: %w", file, err)
		}

		if _, err := io.Copy(tw, f); err != nil {
			f.Close()
			return nil, fmt.Errorf("failed to copy %s: %w", file, err)
		}
		f.Close()
	}

	// Close writers
	if err := tw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close tar writer: %w", err)
	}
	if err := gw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	return buf.Bytes(), nil
}

// ExtractTarball extracts a gzipped tarball to the given directory
func ExtractTarball(data []byte, destDir string) ([]string, error) {
	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	var extracted []string

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tar entry: %w", err)
		}

		// Only handle files
		if header.Typeflag != tar.TypeReg {
			continue
		}

		// Create full path
		targetPath := filepath.Join(destDir, header.Name)

		// Create directory if needed
		targetDir := filepath.Dir(targetPath)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", targetDir, err)
		}

		// Create file
		outFile, err := os.Create(targetPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create %s: %w", targetPath, err)
		}

		// Copy contents
		if _, err := io.Copy(outFile, tr); err != nil {
			outFile.Close()
			return nil, fmt.Errorf("failed to extract %s: %w", targetPath, err)
		}
		outFile.Close()

		// Set permissions
		if err := os.Chmod(targetPath, os.FileMode(header.Mode)); err != nil {
			return nil, fmt.Errorf("failed to set permissions on %s: %w", targetPath, err)
		}

		extracted = append(extracted, header.Name)
	}

	return extracted, nil
}

// Encrypt encrypts data using AES-256-CBC with PBKDF2 key derivation
// Compatible with: openssl enc -aes-256-cbc -salt -pbkdf2
func Encrypt(data []byte, passphrase string) ([]byte, error) {
	// Generate random salt
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key and IV using PBKDF2 (OpenSSL derives both in one call)
	keyIV := pbkdf2.Key([]byte(passphrase), salt, pbkdf2Iterations, aesKeySize+aesIVSize, sha256.New)
	key := keyIV[:aesKeySize]
	iv := keyIV[aesKeySize:]

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Pad data to block size using PKCS#7
	padded := pkcs7Pad(data, aes.BlockSize)

	// Encrypt using CBC mode
	ciphertext := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)

	// Prepend OpenSSL header: "Salted__" + salt
	result := make([]byte, 0, len(opensslSaltHeader)+saltLength+len(ciphertext))
	result = append(result, []byte(opensslSaltHeader)...)
	result = append(result, salt...)
	result = append(result, ciphertext...)

	return result, nil
}

// Decrypt decrypts data encrypted with Encrypt or openssl
func Decrypt(data []byte, passphrase string) ([]byte, error) {
	// Check minimum length
	headerLen := len(opensslSaltHeader) + saltLength
	if len(data) < headerLen {
		return nil, fmt.Errorf("data too short to be valid encrypted data")
	}

	// Verify OpenSSL header
	if string(data[:len(opensslSaltHeader)]) != opensslSaltHeader {
		return nil, fmt.Errorf("invalid encrypted data: missing 'Salted__' header")
	}

	// Extract salt
	salt := data[len(opensslSaltHeader) : len(opensslSaltHeader)+saltLength]
	ciphertext := data[headerLen:]

	// Derive key and IV using PBKDF2 (OpenSSL derives both in one call)
	keyIV := pbkdf2.Key([]byte(passphrase), salt, pbkdf2Iterations, aesKeySize+aesIVSize, sha256.New)
	key := keyIV[:aesKeySize]
	iv := keyIV[aesKeySize:]

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Decrypt using CBC mode
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of block size")
	}

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove PKCS#7 padding
	unpadded, err := pkcs7Unpad(plaintext, aes.BlockSize)
	if err != nil {
		return nil, fmt.Errorf("failed to unpad: %w (wrong passphrase?)", err)
	}

	return unpadded, nil
}

// EncodeBase64 encodes data to base64
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeBase64 decodes base64 data
func DecodeBase64(data string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}
	return decoded, nil
}

// pkcs7Pad adds PKCS#7 padding to data
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7Unpad removes PKCS#7 padding from data
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data is empty")
	}

	padding := int(data[len(data)-1])
	if padding > blockSize || padding == 0 {
		return nil, fmt.Errorf("invalid padding")
	}

	// Verify padding
	for i := 0; i < padding; i++ {
		if data[len(data)-1-i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding bytes")
		}
	}

	return data[:len(data)-padding], nil
}
