package cmd

import (
	"fmt"
	"strings"

	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/config"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse available subscription tiers",
	Long:  `Browse all available Butler Coffee subscription tiers and learn more about each one.`,
	RunE:  runBrowse,
}

func init() {
	// This will be added to subscriptionsCmd in subscribe.go
	// We need to register it there to avoid circular imports
}

func runBrowse(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client := api.NewClient(cfg)

	// Get available subscriptions
	available, err := client.GetAvailableSubscriptions()
	if err != nil {
		return fmt.Errorf("failed to get available subscriptions: %w", err)
	}

	if len(available) == 0 {
		fmt.Println("No subscription tiers available at this time.")
		return nil
	}

	// Get user's active subscriptions if authenticated
	var activeSubscriptions []api.Subscription
	if cfg.IsAuthenticated() {
		activeSubscriptions, err = client.ListSubscriptions()
		if err != nil {
			// Don't fail if we can't get user subscriptions, just continue
			fmt.Printf("Note: Could not fetch your active subscriptions: %v\n\n", err)
		}
	}

	// Create a map of active subscription tiers
	activeTiers := make(map[string]api.Subscription)
	for _, sub := range activeSubscriptions {
		activeTiers[sub.Tier] = sub
	}

	// Create display items for the prompt
	type promptItem struct {
		Name        string
		Description string
		Price       string
		Tier        string
		IsActive    bool
		Status      string
	}

	items := make([]promptItem, len(available)+1)
	for i, sub := range available {
		activeSub, isActive := activeTiers[sub.Tier]

		item := promptItem{
			Name:        sub.Name,
			Description: sub.Description,
			Price:       fmt.Sprintf("%s %s/%s", sub.Currency, sub.Price, sub.BillingPeriod),
			Tier:        sub.Tier,
			IsActive:    isActive && activeSub.Status == "active",
			Status:      "",
		}

		if isActive {
			item.Status = activeSub.Status
		}

		items[i] = item
	}

	// Add exit option
	items[len(available)] = promptItem{
		Name:        "← Exit",
		Description: "Return to main menu",
		Price:       "",
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ .Name | cyan }}{{ if .IsActive }} ✓{{ end }}",
		Inactive: "  {{ .Name }}{{ if .IsActive }} ✓{{ end }}",
		Selected: "{{ .Name | green }}",
		Details: `
--------- Subscription Details ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Price:" | faint }}	{{ .Price }}
{{ "Description:" | faint }}	{{ .Description }}{{if .IsActive}}
{{ "Status:" | faint }}	{{ .Status | green }}{{end}}`,
	}

	prompt := promptui.Select{
		Label:     "Select a subscription tier to learn more",
		Items:     items,
		Templates: templates,
		Size:      10,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		// User cancelled (Ctrl+C) or error
		fmt.Println("\nExiting...")
		return nil
	}

	// Check if user selected exit
	if idx == len(available) {
		return nil
	}

	// Display detailed information about selected subscription
	selectedSub := available[idx]
	activeSub := activeTiers[selectedSub.Tier]

	displaySubscriptionDetails(selectedSub, activeSub, cfg.IsAuthenticated())

	// Ask if user wants to subscribe (if not already active and authenticated)
	if cfg.IsAuthenticated() && (activeSub.ID == "" || activeSub.Status != "active") {
		fmt.Println()
		confirmPrompt := promptui.Prompt{
			Label:     fmt.Sprintf("Would you like to subscribe to %s now", selectedSub.Name),
			IsConfirm: true,
		}

		result, err := confirmPrompt.Run()
		if err == nil && (result == "y" || result == "yes") {
			// User wants to subscribe
			return runSubscribe(nil, []string{selectedSub.Tier})
		}
	}

	return nil
}

func displaySubscriptionDetails(sub api.AvailableSubscription, activeSub api.Subscription, isAuthenticated bool) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("%s\n", sub.Name)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("\nPrice: %s %s/%s\n", sub.Currency, sub.Price, sub.BillingPeriod)
	fmt.Printf("Description: %s\n", sub.Description)

	// Show active status if user has this subscription
	if activeSub.ID != "" {
		fmt.Printf("\nStatus: %s", strings.ToUpper(activeSub.Status))
		if activeSub.Status == "active" {
			fmt.Print(" ✓")
		}
		fmt.Println()

		if activeSub.StartedAt != nil {
			fmt.Printf("Started: %s\n", *activeSub.StartedAt)
		}
		if activeSub.ExpiresAt != nil {
			fmt.Printf("Expires: %s\n", *activeSub.ExpiresAt)
		}
	}

	fmt.Println("\nFeatures:")
	for _, feature := range sub.Features {
		fmt.Printf("  • %s\n", feature)
	}

	// Show subscription option if not already active
	if activeSub.ID == "" || activeSub.Status != "active" {
		fmt.Printf("\nTo subscribe to %s, run:\n", sub.Name)
		fmt.Printf("  bc-cli subscriptions subscribe %s\n", sub.Tier)
	}

	fmt.Println()
}
