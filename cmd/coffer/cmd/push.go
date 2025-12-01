package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/api"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/config"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/crypto"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/registry"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/ui"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push [group]",
	Short: "Upload files to worker",
	Long:  `Encrypt and upload tracked files to Cloudflare Worker.`,
	Example: `  coffer push
  coffer push github`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Determine group (from args or flag)
		grp := group
		if len(args) > 0 {
			grp = args[0]
		}

		// Load and validate config
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if err := cfg.Validate(); err != nil {
			return err
		}

		// Show group being used if not default
		if grp != config.DefaultGroup {
			ui.Info("Using group: %s", grp)
		}

		// Create registry
		reg, err := registry.New(grp)
		if err != nil {
			return fmt.Errorf("failed to create registry: %w", err)
		}

		// Check if empty
		isEmpty, err := reg.IsEmpty()
		if err != nil {
			return err
		}
		if isEmpty {
			ui.Error("No secrets to push for group: %s", grp)
			ui.Warning("Add secrets with: coffer add <path> -g %s", grp)
			return fmt.Errorf("no secrets to push")
		}

		// Get all tracked files
		files, err := reg.All()
		if err != nil {
			return fmt.Errorf("failed to get tracked files: %w", err)
		}

		// Verify all files exist
		homeDir, _ := os.UserHomeDir()
		var missing int
		for _, file := range files {
			fullPath := fmt.Sprintf("%s/%s", homeDir, file)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				ui.Error("Missing file: %s", file)
				missing++
			}
		}
		if missing > 0 {
			return fmt.Errorf("cannot push: %d files missing", missing)
		}

		// Create tarball
		ui.Info("Creating tarball...")
		tarData, err := crypto.CreateTarball(files, homeDir)
		if err != nil {
			return fmt.Errorf("failed to create tarball: %w", err)
		}

		// Calculate size for display
		tarSizeKB := float64(len(tarData)) / 1024.0
		var sizeStr string
		if tarSizeKB < 1024 {
			sizeStr = fmt.Sprintf("%.1fK", tarSizeKB)
		} else {
			sizeStr = fmt.Sprintf("%.1fM", tarSizeKB/1024.0)
		}

		// Encrypt
		ui.Info("Encrypting (%s)...", sizeStr)
		encrypted, err := crypto.Encrypt(tarData, cfg.Passphrase)
		if err != nil {
			return fmt.Errorf("failed to encrypt: %w", err)
		}

		// Base64 encode
		ui.Info("Encoding and uploading...")
		base64Data := crypto.EncodeBase64(encrypted)

		// Upload to worker
		client := api.NewClient(cfg.URL, cfg.Passphrase)
		metadata := api.Metadata{
			Files: files,
			Size:  sizeStr,
		}

		if err := client.Push(grp, base64Data, metadata); err != nil {
			return fmt.Errorf("failed to push to worker: %w", err)
		}

		// Success
		ui.Success("Secrets pushed to worker")
		ui.Success("  Group: %s", grp)
		ui.Success("  Files: %d", len(files))
		ui.Success("  Size: %s", sizeStr)
		fmt.Println()
		ui.Success("Files pushed:")
		for _, file := range files {
			ui.Plain("  - %s", file)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
