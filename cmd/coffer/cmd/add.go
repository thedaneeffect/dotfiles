package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/config"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/registry"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/ui"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <path> [<path>...]",
	Short: "Add files or directories to coffer registry",
	Long: `Add files or directories to the coffer registry.

Directories are recursively expanded to include all files.`,
	Example: `  coffer add ~/.ssh/id_rsa
  coffer add ~/.ssh              # Adds all files in directory
  coffer add ~/.env -g work      # Add to 'work' group`,
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

		// Add paths
		if err := reg.Add(args); err != nil {
			return fmt.Errorf("failed to add paths: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
