package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/api"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/config"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/registry"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/ui"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete [group]",
	Aliases: []string{"del"},
	Short:   "Delete group from worker",
	Long:    `Delete a secret group from the Cloudflare Worker.`,
	Example: `  coffer delete github
  coffer del test-group`,
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
			ui.Info("Group: %s", grp)
		}

		// Confirm deletion
		fmt.Printf("Delete group '%s' from worker? (y/N) ", grp)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			ui.Warning("Cancelled")
			return nil
		}

		// Delete from worker
		ui.Info("Deleting from worker...")
		client := api.NewClient(cfg.URL, cfg.Passphrase)
		if err := client.Delete(grp); err != nil {
			return err
		}

		ui.Success("Deleted group from worker: %s", grp)

		// Ask about local registry
		reg, _ := registry.New(grp)
		if _, err := os.Stat(reg.Path); err == nil {
			fmt.Printf("Delete local registry file too? (y/N) ")
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}

			response = strings.TrimSpace(strings.ToLower(response))
			if response == "y" || response == "yes" {
				if err := os.Remove(reg.Path); err != nil {
					return fmt.Errorf("failed to delete local registry: %w", err)
				}
				ui.Success("Deleted local registry: %s", reg.Path)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
