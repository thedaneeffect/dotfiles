package registry

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/config"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/ui"
)

// Registry manages a local secrets registry file
type Registry struct {
	Group string
	Path  string
}

// RegistryEntry represents a file tracked in the registry
type RegistryEntry struct {
	Path   string
	Exists bool
}

// New creates a new registry for the given group
func New(group string) (*Registry, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	registryPath := filepath.Join(homeDir, ".coffer")
	if group != config.DefaultGroup {
		registryPath = filepath.Join(homeDir, fmt.Sprintf(".coffer.%s", group))
	}

	return &Registry{
		Group: group,
		Path:  registryPath,
	}, nil
}

// Init creates the registry file if it doesn't exist
func (r *Registry) Init() error {
	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		f, err := os.Create(r.Path)
		if err != nil {
			return fmt.Errorf("failed to create registry: %w", err)
		}
		f.Close()
		ui.Success("Created coffer registry at %s", r.Path)
	}
	return nil
}

// Add adds paths to the registry (expanding directories recursively)
func (r *Registry) Add(paths []string) error {
	if err := r.Init(); err != nil {
		return err
	}

	// Load existing entries
	existing, err := r.All()
	if err != nil {
		return err
	}

	existingSet := make(map[string]bool)
	for _, p := range existing {
		existingSet[p] = true
	}

	homeDir, _ := os.UserHomeDir()
	var newPaths []string

	for _, path := range paths {
		// Expand ~ to home directory
		if strings.HasPrefix(path, "~/") {
			path = filepath.Join(homeDir, path[2:])
		} else if path == "~" {
			path = homeDir
		}

		// Make absolute if relative
		if !filepath.IsAbs(path) {
			absPath, err := filepath.Abs(path)
			if err != nil {
				ui.Error("Failed to resolve path: %s", path)
				continue
			}
			path = absPath
		}

		// Check if path exists
		info, err := os.Stat(path)
		if err != nil {
			ui.Error("Path does not exist: %s", path)
			continue
		}

		// If directory, recursively add all files
		if info.IsDir() {
			ui.Info("Expanding directory: %s", path)
			fileCount := 0

			err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					relPath, _ := filepath.Rel(homeDir, filePath)
					if !existingSet[relPath] {
						newPaths = append(newPaths, relPath)
						existingSet[relPath] = true
						ui.Success("  Added: %s", relPath)
						fileCount++
					} else {
						ui.Warning("  Already tracked: %s", relPath)
					}
				}
				return nil
			})

			if err != nil {
				ui.Error("Failed to walk directory %s: %v", path, err)
			} else {
				ui.Success("Added %d new files from directory", fileCount)
			}
		} else {
			// Single file
			relPath, _ := filepath.Rel(homeDir, path)
			if existingSet[relPath] {
				ui.Warning("Already tracked: %s", relPath)
			} else {
				newPaths = append(newPaths, relPath)
				ui.Success("Added: %s", relPath)
			}
		}
	}

	// Append new paths to registry
	if len(newPaths) > 0 {
		f, err := os.OpenFile(r.Path, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open registry: %w", err)
		}
		defer f.Close()

		for _, p := range newPaths {
			if _, err := f.WriteString(p + "\n"); err != nil {
				return fmt.Errorf("failed to write to registry: %w", err)
			}
		}
	}

	return nil
}

// Remove removes paths from the registry
func (r *Registry) Remove(paths []string) error {
	existing, err := r.All()
	if err != nil {
		return err
	}

	homeDir, _ := os.UserHomeDir()

	// Build set of paths to remove (normalized)
	toRemove := make(map[string]bool)
	for _, path := range paths {
		// Normalize path
		path = strings.TrimPrefix(path, homeDir+"/")
		path = strings.TrimPrefix(path, "./")
		toRemove[path] = true
	}

	// Filter out removed paths
	var kept []string
	var removed []string
	for _, path := range existing {
		if toRemove[path] {
			removed = append(removed, path)
		} else {
			kept = append(kept, path)
		}
	}

	// Write back filtered list
	f, err := os.Create(r.Path)
	if err != nil {
		return fmt.Errorf("failed to write registry: %w", err)
	}
	defer f.Close()

	for _, p := range kept {
		if _, err := f.WriteString(p + "\n"); err != nil {
			return fmt.Errorf("failed to write to registry: %w", err)
		}
	}

	// Report results
	for _, p := range removed {
		ui.Success("Removed: %s", p)
	}

	// Report paths that weren't tracked
	for path := range toRemove {
		found := false
		for _, r := range removed {
			if r == path {
				found = true
				break
			}
		}
		if !found {
			ui.Warning("Not tracked: %s", path)
		}
	}

	return nil
}

// List returns all registry entries with existence status
func (r *Registry) List() ([]RegistryEntry, error) {
	paths, err := r.All()
	if err != nil {
		return nil, err
	}

	homeDir, _ := os.UserHomeDir()
	entries := make([]RegistryEntry, 0, len(paths))

	for _, path := range paths {
		fullPath := filepath.Join(homeDir, path)
		_, err := os.Stat(fullPath)
		exists := err == nil

		entries = append(entries, RegistryEntry{
			Path:   path,
			Exists: exists,
		})
	}

	return entries, nil
}

// All returns all paths in the registry (including non-existent ones)
func (r *Registry) All() ([]string, error) {
	// Return empty list if registry doesn't exist
	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		return []string{}, nil
	}

	f, err := os.Open(r.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open registry: %w", err)
	}
	defer f.Close()

	var paths []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		paths = append(paths, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read registry: %w", err)
	}

	return paths, nil
}

// Exists checks if a path is in the registry
func (r *Registry) Exists(path string) (bool, error) {
	paths, err := r.All()
	if err != nil {
		return false, err
	}

	homeDir, _ := os.UserHomeDir()
	path = strings.TrimPrefix(path, homeDir+"/")
	path = strings.TrimPrefix(path, "./")

	for _, p := range paths {
		if p == path {
			return true, nil
		}
	}

	return false, nil
}

// IsEmpty returns true if the registry has no entries
func (r *Registry) IsEmpty() (bool, error) {
	paths, err := r.All()
	if err != nil {
		return false, err
	}
	return len(paths) == 0, nil
}
