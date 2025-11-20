package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bc-cli",
	Short: "Butler Coffee CLI tool",
	Long: `Butler Coffee CLI tool - A command line interface for managing your coffee operations.

Complete documentation is available at https://github.com/hassek/bc-cli`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Butler Coffee CLI!")
		fmt.Println("Use 'bc-cli --help' to see available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "Print version information")
}
