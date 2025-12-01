package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/api"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/config"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/registry"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/ui"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list [group]",
	Aliases: []string{"ls"},
	Short:   "List tracked files",
	Long:    `List all files tracked in the coffer registry.`,
	Example: `  coffer list
  coffer ls github`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Determine group (from args or flag)
		grp := group
		if len(args) > 0 {
			grp = args[0]
		}

		// Show group being used if not default
		if grp != config.DefaultGroup {
			ui.Info("Group: %s", grp)
		}

		// Create registry
		reg, err := registry.New(grp)
		if err != nil {
			return fmt.Errorf("failed to create registry: %w", err)
		}

		// Check if empty
		isEmpty, err := reg.IsEmpty()
		if err != nil {
			return fmt.Errorf("failed to check registry: %w", err)
		}

		if isEmpty {
			ui.Warning("No files tracked locally for group: %s", grp)
			ui.Warning("Use: coffer add <path> -g %s", grp)
			return nil
		}

		// List entries
		ui.Success("Locally tracked files:")
		entries, err := reg.List()
		if err != nil {
			return fmt.Errorf("failed to list files: %w", err)
		}

		for _, entry := range entries {
			if entry.Exists {
				ui.Plain("  ✓ %s", entry.Path)
			} else {
				ui.Warning("  ✗ %s (missing)", entry.Path)
			}
		}

		fmt.Println()

		// Show remote metadata if worker is configured
		cfg, err := config.Load()
		if err == nil && cfg.URL != "" && cfg.Passphrase != "" {
			ui.Success("Remote files (in worker):")

			client := api.NewClient(cfg.URL, cfg.Passphrase)
			metadata, err := client.GetMetadata()
			if err != nil {
				ui.Warning("  (failed to fetch metadata: %v)", err)
				return nil
			}

			groupMeta, exists := metadata[grp]
			if !exists {
				ui.Warning("  (no files in worker for this group)")
			} else {
				if len(groupMeta.Files) == 0 {
					ui.Warning("  (no files)")
				} else {
					for _, file := range groupMeta.Files {
						ui.Plain("  ✓ %s", file)
					}
					fmt.Println()
					ui.Info("  Size: %s | Uploaded: %s", groupMeta.Size, groupMeta.Uploaded.Format("2006-01-02 15:04:05"))
				}
			}
		} else {
			ui.Warning("Remote files: (COFFER_URL/COFFER_PASSPHRASE not configured)")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
