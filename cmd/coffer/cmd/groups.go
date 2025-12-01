package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/api"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/config"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/ui"
)

// groupsCmd represents the groups command
var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List all groups in worker",
	Long:  `List all secret groups stored in the Cloudflare Worker.`,
	Example: `  coffer groups`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load and validate config
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if err := cfg.Validate(); err != nil {
			return err
		}

		// Fetch groups from worker
		ui.Info("Fetching groups from worker...")
		client := api.NewClient(cfg.URL, cfg.Passphrase)
		groups, err := client.ListGroups()
		if err != nil {
			return err
		}

		if len(groups) == 0 {
			ui.Warning("No groups found in worker")
			ui.Warning("Push a group with: coffer push [GROUP]")
			return nil
		}

		ui.Success("Available groups in worker:")
		for _, group := range groups {
			ui.Plain("  - %s", group)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(groupsCmd)
}
