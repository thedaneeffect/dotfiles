package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/api"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/config"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/crypto"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/ui"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull [group]",
	Short: "Download files from worker",
	Long:  `Download and decrypt secrets from Cloudflare Worker.`,
	Example: `  coffer pull
  coffer pull github`,
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

		// Download from worker
		ui.Info("Downloading from worker...")
		client := api.NewClient(cfg.URL, cfg.Passphrase)
		base64Data, err := client.Pull(grp)
		if err != nil {
			return err
		}

		// Base64 decode
		ui.Info("Decoding and decrypting secrets...")
		encrypted, err := crypto.DecodeBase64(base64Data)
		if err != nil {
			return fmt.Errorf("failed to decode base64: %w", err)
		}

		// Decrypt
		tarData, err := crypto.Decrypt(encrypted, cfg.Passphrase)
		if err != nil {
			return fmt.Errorf("failed to decrypt: %w (wrong passphrase?)", err)
		}

		// Calculate size for display
		tarSizeKB := float64(len(tarData)) / 1024.0
		var sizeStr string
		if tarSizeKB < 1024 {
			sizeStr = fmt.Sprintf("%.1fK", tarSizeKB)
		} else {
			sizeStr = fmt.Sprintf("%.1fM", tarSizeKB/1024.0)
		}

		// Extract to home directory
		homeDir, _ := os.UserHomeDir()
		extractedFiles, err := crypto.ExtractTarball(tarData, homeDir)
		if err != nil {
			return fmt.Errorf("failed to extract tarball: %w", err)
		}

		// Success
		ui.Success("Secrets pulled from worker")
		ui.Success("  Group: %s", grp)
		ui.Success("  Files: %d", len(extractedFiles))
		ui.Success("  Size: %s", sizeStr)
		fmt.Println()
		ui.Success("Files restored:")
		for _, file := range extractedFiles {
			ui.Plain("  - %s", file)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
