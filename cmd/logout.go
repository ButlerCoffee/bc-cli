package cmd

import (
	"fmt"

	"github.com/hassek/bc-cli/config"
	"github.com/hassek/bc-cli/templates"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from your Butler Coffee account",
	Long:  `Clear your authentication tokens and logout from your Butler Coffee account.`,
	RunE:  runLogout,
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

func runLogout(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if !cfg.IsAuthenticated() {
		if err := templates.RenderToStdout(templates.NotLoggedInTemplate, nil); err != nil {
			return fmt.Errorf("failed to render template: %w", err)
		}
		return nil
	}

	// Clear authentication tokens
	cfg.AccessToken = ""
	cfg.RefreshToken = ""
	cfg.ExpiresAt = ""
	cfg.RefreshTokenExpiresAt = ""

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	if err := templates.RenderToStdout(templates.LogoutSuccessTemplate, nil); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return nil
}
