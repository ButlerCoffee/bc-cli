package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to your Butler Coffee account",
	Long:  `Authenticate with your Butler Coffee account using your username and password.`,
	RunE:  runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.IsAuthenticated() {
		fmt.Println("You are already logged in.")
		fmt.Print("Do you want to login with a different account? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" && response != "yes" {
			return nil
		}
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read username: %w", err)
	}
	username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	password := string(passwordBytes)
	fmt.Println()

	client := api.NewClient(cfg)

	fmt.Println("\nAuthenticating...")
	_, err = client.Login(api.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	fmt.Println("\n✓ Successfully logged in!")
	fmt.Printf("Welcome back, %s!\n", username)

	return nil
}
