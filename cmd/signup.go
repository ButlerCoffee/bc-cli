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

var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Create a new Butler Coffee account",
	Long:  `Create a new Butler Coffee account by providing your username, email, password, and invitation code.`,
	RunE:  runSignup,
}

func init() {
	rootCmd.AddCommand(signupCmd)
}

func runSignup(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to Butler Coffee! Let's create your account.")
	fmt.Println()

	fmt.Print("Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read username: %w", err)
	}
	username = strings.TrimSpace(username)

	fmt.Print("Email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read email: %w", err)
	}
	email = strings.TrimSpace(email)

	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	password := string(passwordBytes)
	fmt.Println()

	fmt.Print("Confirm Password: ")
	confirmPasswordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password confirmation: %w", err)
	}
	confirmPassword := string(confirmPasswordBytes)
	fmt.Println()

	if password != confirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	fmt.Print("Invitation Code (default is empty, press Enter to skip): ")
	code, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read invitation code: %w", err)
	}
	code = strings.TrimSpace(code)

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client := api.NewClient(cfg)

	fmt.Println("\nCreating account...")
	resp, err := client.Register(api.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
		Code:     code,
	})
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	fmt.Println("\n✓ Account created successfully!")
	fmt.Printf("User ID: %s\n", resp.ID)
	fmt.Println("\nYou are now logged in and ready to use Butler Coffee CLI!")

	return nil
}
