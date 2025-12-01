package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/config"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/registry"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/ui"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove <path> [<path>...]",
	Aliases: []string{"rm"},
	Short:   "Remove files from coffer registry",
	Long:    `Remove files from the coffer registry.`,
	Example: `  coffer remove ~/.ssh/id_rsa
  secrets rm ~/.env -g work`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Show group being used if not default
		if group != config.DefaultGroup {
			ui.Info("Using group: %s", group)
		}

		// Create registry
		reg, err := registry.New(group)
		if err != nil {
			return fmt.Errorf("failed to create registry: %w", err)
		}

		// Remove paths
		if err := reg.Remove(args); err != nil {
			return fmt.Errorf("failed to remove paths: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
