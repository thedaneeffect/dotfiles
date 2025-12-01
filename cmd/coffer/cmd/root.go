package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/thedaneeffect/dotfiles/cmd/coffer/internal/config"
)

var (
	// Global flags
	group string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "coffer",
	Short:   "Manage encrypted dotfiles with Cloudflare Workers",
	Long:    `Encrypted file management for dotfiles using Cloudflare Workers KV storage.`,
	Version: "1.0.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global --group/-g flag available to all commands
	rootCmd.PersistentFlags().StringVarP(&group, "group", "g", config.DefaultGroup, "coffer group name")
}


