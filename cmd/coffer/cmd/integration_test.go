package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/api"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/config"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/crypto"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/registry"
)

// TestWorkerIntegration tests push/pull/delete/groups with actual Cloudflare Worker
// Requires COFFER_URL and COFFER_PASSPHRASE environment variables
// Skips if not configured
func TestWorkerIntegration(t *testing.T) {
	// Check if worker credentials are configured
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.URL == "" || cfg.Passphrase == "" {
		t.Skip("Skipping integration test: COFFER_URL or COFFER_PASSPHRASE not configured")
	}

	// Use a dedicated test group
	testGroup := "integration-test-group"

	// Create temp directory for test
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, "home")
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatalf("Failed to create temp home: %v", err)
	}

	// Override HOME for test
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", oldHome)

	// Create test files
	testFiles := map[string]string{
		"test1.txt": "integration test content 1",
		"test2.txt": "integration test content 2",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(homeDir, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Create registry and add files
	reg, err := registry.New(testGroup)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	var filePaths []string
	for path := range testFiles {
		filePaths = append(filePaths, filepath.Join(homeDir, path))
	}

	if err := reg.Add(filePaths); err != nil {
		t.Fatalf("Failed to add files to registry: %v", err)
	}

	// Test 1: Push to worker
	t.Run("push to worker", func(t *testing.T) {
		files, err := reg.All()
		if err != nil {
			t.Fatalf("Failed to get tracked files: %v", err)
		}

		// Create tarball
		tarData, err := crypto.CreateTarball(files, homeDir)
		if err != nil {
			t.Fatalf("Failed to create tarball: %v", err)
		}

		// Encrypt
		encrypted, err := crypto.Encrypt(tarData, cfg.Passphrase)
		if err != nil {
			t.Fatalf("Failed to encrypt: %v", err)
		}

		// Base64 encode
		base64Data := crypto.EncodeBase64(encrypted)

		// Push to worker
		client := api.NewClient(cfg.URL, cfg.Passphrase)
		metadata := api.Metadata{
			Files: files,
			Size:  "test",
		}

		if err := client.Push(testGroup, base64Data, metadata); err != nil {
			t.Fatalf("Failed to push to worker: %v", err)
		}

		t.Log("Successfully pushed to worker")
	})

	// Test 2: List groups
	t.Run("list groups", func(t *testing.T) {
		client := api.NewClient(cfg.URL, cfg.Passphrase)
		groups, err := client.ListGroups()
		if err != nil {
			t.Fatalf("Failed to list groups: %v", err)
		}

		// Check that our test group is in the list
		found := false
		for _, g := range groups {
			if g == testGroup {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Test group '%s' not found in groups list: %v", testGroup, groups)
		}

		t.Logf("Found %d groups (including test group)", len(groups))
	})

	// Test 3: Get metadata
	t.Run("get metadata", func(t *testing.T) {
		client := api.NewClient(cfg.URL, cfg.Passphrase)
		metadata, err := client.GetMetadata()
		if err != nil {
			t.Fatalf("Failed to get metadata: %v", err)
		}

		groupMeta, exists := metadata[testGroup]
		if !exists {
			t.Fatalf("Test group metadata not found")
		}

		if len(groupMeta.Files) != 2 {
			t.Errorf("Expected 2 files in metadata, got %d", len(groupMeta.Files))
		}

		t.Logf("Metadata: %d files, size: %s", len(groupMeta.Files), groupMeta.Size)
	})

	// Test 4: Pull from worker
	t.Run("pull from worker", func(t *testing.T) {
		// Delete local test files first
		for path := range testFiles {
			if err := os.Remove(filepath.Join(homeDir, path)); err != nil {
				t.Fatalf("Failed to remove local file: %v", err)
			}
		}

		// Pull from worker
		client := api.NewClient(cfg.URL, cfg.Passphrase)
		base64Data, err := client.Pull(testGroup)
		if err != nil {
			t.Fatalf("Failed to pull from worker: %v", err)
		}

		// Decode
		encrypted, err := crypto.DecodeBase64(base64Data)
		if err != nil {
			t.Fatalf("Failed to decode base64: %v", err)
		}

		// Decrypt
		tarData, err := crypto.Decrypt(encrypted, cfg.Passphrase)
		if err != nil {
			t.Fatalf("Failed to decrypt: %v", err)
		}

		// Extract
		extractedFiles, err := crypto.ExtractTarball(tarData, homeDir)
		if err != nil {
			t.Fatalf("Failed to extract tarball: %v", err)
		}

		if len(extractedFiles) != 2 {
			t.Errorf("Expected 2 extracted files, got %d", len(extractedFiles))
		}

		// Verify content
		for path, expectedContent := range testFiles {
			fullPath := filepath.Join(homeDir, path)
			content, err := os.ReadFile(fullPath)
			if err != nil {
				t.Errorf("Failed to read restored file %s: %v", path, err)
				continue
			}

			if string(content) != expectedContent {
				t.Errorf("File %s content mismatch.\nExpected: %q\nGot: %q",
					path, expectedContent, string(content))
			}
		}

		t.Log("Successfully pulled and verified files")
	})

	// Test 5: Delete from worker (cleanup)
	t.Run("delete from worker", func(t *testing.T) {
		client := api.NewClient(cfg.URL, cfg.Passphrase)
		if err := client.Delete(testGroup); err != nil {
			t.Fatalf("Failed to delete test group: %v", err)
		}

		// Verify deletion by trying to pull (should fail)
		_, err := client.Pull(testGroup)
		if err == nil {
			t.Error("Expected error when pulling deleted group, got nil")
		}

		t.Log("Successfully deleted test group")
	})
}

// TestWorkerRoundtrip tests the complete roundtrip workflow
func TestWorkerRoundtrip(t *testing.T) {
	cfg, err := config.Load()
	if err != nil || cfg.URL == "" || cfg.Passphrase == "" {
		t.Skip("Skipping roundtrip test: COFFER_URL or COFFER_PASSPHRASE not configured")
	}

	testGroup := "roundtrip-test"

	// Setup
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, "home")
	os.MkdirAll(homeDir, 0755)

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", oldHome)

	// Create test file with specific content
	testContent := "This is a roundtrip test with special chars: !@#$%^&*()"
	testFile := "roundtrip.txt"
	if err := os.WriteFile(filepath.Join(homeDir, testFile), []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add to registry
	reg, _ := registry.New(testGroup)
	reg.Add([]string{filepath.Join(homeDir, testFile)})

	// Get files and create tarball
	files, _ := reg.All()
	tarData, _ := crypto.CreateTarball(files, homeDir)
	encrypted, _ := crypto.Encrypt(tarData, cfg.Passphrase)
	base64Data := crypto.EncodeBase64(encrypted)

	// Push
	client := api.NewClient(cfg.URL, cfg.Passphrase)
	metadata := api.Metadata{Files: files, Size: "test"}
	if err := client.Push(testGroup, base64Data, metadata); err != nil {
		t.Fatalf("Push failed: %v", err)
	}

	// Delete local file
	os.Remove(filepath.Join(homeDir, testFile))

	// Pull
	base64Data, err = client.Pull(testGroup)
	if err != nil {
		t.Fatalf("Pull failed: %v", err)
	}

	// Decrypt and extract
	encrypted, _ = crypto.DecodeBase64(base64Data)
	tarData, _ = crypto.Decrypt(encrypted, cfg.Passphrase)
	crypto.ExtractTarball(tarData, homeDir)

	// Verify content matches exactly
	content, _ := os.ReadFile(filepath.Join(homeDir, testFile))
	if string(content) != testContent {
		t.Errorf("Roundtrip content mismatch.\nExpected: %q\nGot: %q", testContent, string(content))
	}

	// Cleanup
	client.Delete(testGroup)

	t.Log("Roundtrip test passed: content matches exactly")
}
