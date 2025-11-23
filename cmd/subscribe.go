package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/config"
	"github.com/spf13/cobra"
)

var subscriptionsCmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "Manage your Butler Coffee subscriptions",
	Long:  `Manage your Butler Coffee subscriptions - list, browse available tiers, and subscribe to new plans.`,
	RunE:  runSubscriptions,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List your active subscriptions",
	Long:  `View all your Butler Coffee subscriptions and their current status.`,
	RunE:  runSubscriptions,
}

var subscribeCmd = &cobra.Command{
	Use:   "subscribe [tier]",
	Short: "Subscribe to a Butler Coffee subscription tier",
	Long: `Subscribe to a Butler Coffee subscription tier:
  - butler: Butler Coffee - High-quality curated selections
  - collection: Collection Coffee - Premium rare and exclusive selections
  - premium: Premium Coffee - The absolute finest coffees in the world

Example:
  bc-cli subscriptions subscribe butler
  bc-cli subscriptions subscribe collection
  bc-cli subscriptions subscribe premium`,
	Args: cobra.ExactArgs(1),
	RunE: runSubscribe,
}

func init() {
	rootCmd.AddCommand(subscriptionsCmd)
	subscriptionsCmd.AddCommand(listCmd)
	subscriptionsCmd.AddCommand(browseCmd)
	subscriptionsCmd.AddCommand(subscribeCmd)
}

func runSubscribe(cmd *cobra.Command, args []string) error {
	tier := args[0]

	// Validate tier
	validTiers := map[string]bool{
		"butler":     true,
		"collection": true,
		"premium":    true,
	}

	if !validTiers[tier] {
		return fmt.Errorf("invalid tier: %s. Valid options are: butler, collection, premium", tier)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if !cfg.IsAuthenticated() {
		return fmt.Errorf("you must be logged in to subscribe. Please run 'bc-cli login' first")
	}

	client := api.NewClient(cfg)

	fmt.Printf("Getting payment link for %s tier...\n", tier)
	resp, err := client.GetSubscriptionPaymentLink(tier)
	if err != nil {
		return fmt.Errorf("failed to get payment link: %w", err)
	}

	fmt.Println("\n✓ Payment link ready!")
	fmt.Printf("Opening browser to complete your subscription...\n\n")
	fmt.Printf("Payment URL: %s\n\n", resp.PaymentLink)

	// Open browser
	if err := openBrowser(resp.PaymentLink); err != nil {
		fmt.Printf("⚠ Could not open browser automatically: %v\n", err)
		fmt.Printf("Please open the URL above in your browser to complete payment.\n")
	} else {
		fmt.Println("✓ Browser opened! Complete your payment to activate your subscription.")
	}

	return nil
}

func runSubscriptions(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if !cfg.IsAuthenticated() {
		return fmt.Errorf("you must be logged in to view subscriptions. Please run 'bc-cli login' first")
	}

	client := api.NewClient(cfg)

	subscriptions, err := client.ListSubscriptions()
	if err != nil {
		return fmt.Errorf("failed to list subscriptions: %w", err)
	}

	if len(subscriptions) == 0 {
		fmt.Println("You don't have any subscriptions yet.")
		fmt.Println("\nBrowse available subscription tiers:")
		fmt.Println("  bc-cli subscriptions browse")
		fmt.Println("\nOr subscribe directly to a tier:")
		fmt.Println("  bc-cli subscriptions subscribe butler")
		fmt.Println("  bc-cli subscriptions subscribe collection")
		fmt.Println("  bc-cli subscriptions subscribe premium")
		return nil
	}

	fmt.Println("=== Your Butler Coffee Subscriptions ===")
	fmt.Println()

	hasActive := false
	for _, sub := range subscriptions {
		isActive := sub.Status == "active"
		if isActive {
			hasActive = true
		}

		// Display with visual indicators
		statusIndicator := ""
		if isActive {
			statusIndicator = " ✓"
		}

		fmt.Printf("┌─ %s%s\n", strings.ToUpper(sub.Tier), statusIndicator)
		fmt.Printf("│  Status: %s\n", strings.ToUpper(sub.Status))

		if sub.StartedAt != nil {
			fmt.Printf("│  Started: %s\n", *sub.StartedAt)
		}
		if sub.ExpiresAt != nil {
			fmt.Printf("│  Expires: %s\n", *sub.ExpiresAt)
		}

		fmt.Printf("└─ ID: %s\n", sub.ID)
		fmt.Println()
	}

	if hasActive {
		fmt.Println("✓ = Active subscription")
		fmt.Println()
	}

	fmt.Println("To browse all available tiers, run: bc-cli subscriptions browse")

	return nil
}

// openBrowser opens the specified URL in the default browser
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
