package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/config"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/registry"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/ui"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:     "status [group]",
	Aliases: []string{"st"},
	Short:   "Show status of tracked files",
	Long:    `Show status of tracked files, including which files exist and which are missing.`,
	Example: `  coffer status
  coffer st github`,
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
			ui.Warning("No files tracked for group: %s", grp)
			return nil
		}

		// Get entries
		entries, err := reg.List()
		if err != nil {
			return fmt.Errorf("failed to list files: %w", err)
		}

		// Count existing vs missing
		var existing, missing int
		for _, entry := range entries {
			if entry.Exists {
				existing++
			} else {
				missing++
				ui.Error("Missing: %s", entry.Path)
			}
		}

		// Print summary
		fmt.Println()
		ui.Success("%d files exist", existing)
		if missing > 0 {
			ui.Error("%d files missing", missing)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
