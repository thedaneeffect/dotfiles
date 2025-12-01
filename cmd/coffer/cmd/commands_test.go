package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/registry"
)

// TestLocalCommandsIntegration tests add, list, status, remove commands
func TestLocalCommandsIntegration(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, "home")
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatalf("Failed to create temp home: %v", err)
	}

	// Override HOME for registry
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", oldHome)

	// Create test files to track
	testFiles := map[string]string{
		"test1.txt":     "content 1",
		"dir/test2.txt": "content 2",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(homeDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", path, err)
		}
	}

	// Test group
	testGroup := "test-group"

	// Test 1: Add files to registry
	t.Run("add files", func(t *testing.T) {
		reg, err := registry.New(testGroup)
		if err != nil {
			t.Fatalf("Failed to create registry: %v", err)
		}

		paths := []string{
			filepath.Join(homeDir, "test1.txt"),
			filepath.Join(homeDir, "dir/test2.txt"),
		}

		if err := reg.Add(paths); err != nil {
			t.Fatalf("Failed to add files: %v", err)
		}

		// Verify registry contains files
		all, err := reg.All()
		if err != nil {
			t.Fatalf("Failed to get all files: %v", err)
		}

		if len(all) != 2 {
			t.Errorf("Expected 2 files in registry, got %d", len(all))
		}
	})

	// Test 2: List files
	t.Run("list files", func(t *testing.T) {
		reg, err := registry.New(testGroup)
		if err != nil {
			t.Fatalf("Failed to create registry: %v", err)
		}

		entries, err := reg.List()
		if err != nil {
			t.Fatalf("Failed to list files: %v", err)
		}

		if len(entries) != 2 {
			t.Errorf("Expected 2 entries, got %d", len(entries))
		}

		// All files should exist
		for _, entry := range entries {
			if !entry.Exists {
				t.Errorf("File %s should exist", entry.Path)
			}
		}
	})

	// Test 3: Status (all files exist)
	t.Run("status - all exist", func(t *testing.T) {
		reg, err := registry.New(testGroup)
		if err != nil {
			t.Fatalf("Failed to create registry: %v", err)
		}

		entries, err := reg.List()
		if err != nil {
			t.Fatalf("Failed to list files: %v", err)
		}

		var existing, missing int
		for _, entry := range entries {
			if entry.Exists {
				existing++
			} else {
				missing++
			}
		}

		if existing != 2 {
			t.Errorf("Expected 2 existing files, got %d", existing)
		}
		if missing != 0 {
			t.Errorf("Expected 0 missing files, got %d", missing)
		}
	})

	// Test 4: Status after deleting a file
	t.Run("status - one missing", func(t *testing.T) {
		// Delete one test file
		if err := os.Remove(filepath.Join(homeDir, "test1.txt")); err != nil {
			t.Fatalf("Failed to delete test file: %v", err)
		}

		reg, err := registry.New(testGroup)
		if err != nil {
			t.Fatalf("Failed to create registry: %v", err)
		}

		entries, err := reg.List()
		if err != nil {
			t.Fatalf("Failed to list files: %v", err)
		}

		var existing, missing int
		for _, entry := range entries {
			if entry.Exists {
				existing++
			} else {
				missing++
			}
		}

		if existing != 1 {
			t.Errorf("Expected 1 existing file, got %d", existing)
		}
		if missing != 1 {
			t.Errorf("Expected 1 missing file, got %d", missing)
		}
	})

	// Test 5: Remove files from registry
	t.Run("remove files", func(t *testing.T) {
		reg, err := registry.New(testGroup)
		if err != nil {
			t.Fatalf("Failed to create registry: %v", err)
		}

		// Remove one file
		if err := reg.Remove([]string{"test1.txt"}); err != nil {
			t.Fatalf("Failed to remove file: %v", err)
		}

		// Verify only one file remains
		all, err := reg.All()
		if err != nil {
			t.Fatalf("Failed to get all files: %v", err)
		}

		if len(all) != 1 {
			t.Errorf("Expected 1 file after removal, got %d", len(all))
		}

		if all[0] != "dir/test2.txt" {
			t.Errorf("Expected dir/test2.txt to remain, got %s", all[0])
		}
	})

	// Test 6: Add directory (recursive expansion)
	t.Run("add directory recursively", func(t *testing.T) {
		// Create new registry for clean test
		newGroup := "dir-test"
		reg, err := registry.New(newGroup)
		if err != nil {
			t.Fatalf("Failed to create registry: %v", err)
		}

		// Create a directory with multiple files
		dirPath := filepath.Join(homeDir, "multi")
		files := []string{
			filepath.Join(dirPath, "file1.txt"),
			filepath.Join(dirPath, "file2.txt"),
			filepath.Join(dirPath, "subdir/file3.txt"),
		}

		for _, file := range files {
			dir := filepath.Dir(file)
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatalf("Failed to create dir: %v", err)
			}
			if err := os.WriteFile(file, []byte("content"), 0644); err != nil {
				t.Fatalf("Failed to create file: %v", err)
			}
		}

		// Add directory
		if err := reg.Add([]string{dirPath}); err != nil {
			t.Fatalf("Failed to add directory: %v", err)
		}

		// Verify all files were added
		all, err := reg.All()
		if err != nil {
			t.Fatalf("Failed to get all files: %v", err)
		}

		if len(all) != 3 {
			t.Errorf("Expected 3 files from directory, got %d", len(all))
		}
	})
}

// TestRegistryDuplicates tests that adding duplicate paths doesn't duplicate entries
func TestRegistryDuplicates(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, "home")
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatalf("Failed to create temp home: %v", err)
	}

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", oldHome)

	// Create test file
	testFile := filepath.Join(homeDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	reg, err := registry.New("dup-test")
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Add same file twice
	if err := reg.Add([]string{testFile}); err != nil {
		t.Fatalf("Failed to add file first time: %v", err)
	}
	if err := reg.Add([]string{testFile}); err != nil {
		t.Fatalf("Failed to add file second time: %v", err)
	}

	// Should only have one entry
	all, err := reg.All()
	if err != nil {
		t.Fatalf("Failed to get all files: %v", err)
	}

	if len(all) != 1 {
		t.Errorf("Expected 1 file (no duplicates), got %d", len(all))
	}
}
