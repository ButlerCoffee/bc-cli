package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/config"
	"github.com/hassek/bc-cli/templates"
	"github.com/hassek/bc-cli/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var subscriptionsCmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "Manage your Butler Coffee subscriptions",
	Long:  `View your active subscriptions and browse available tiers interactively.`,
	RunE:  runSubscriptions,
}

func init() {
	rootCmd.AddCommand(subscriptionsCmd)
}

func runSubscriptions(cmd *cobra.Command, args []string) error {
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

	// Create display items for the prompt
	type promptItem struct {
		Name        string
		Description string
		Price       string
		Tier        string
	}

	items := make([]promptItem, len(available)+1)
	for i, sub := range available {
		item := promptItem{
			Name:        sub.Name,
			Description: sub.Description,
			Price:       fmt.Sprintf("%s %s/%s", sub.Currency, sub.Price, sub.BillingPeriod),
			Tier:        sub.Tier,
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
		Active:   "▸ {{ .Name | cyan }}",
		Inactive: "  {{ .Name }}",
		Selected: "{{ .Name | green }}",
		Details: `
--------- Subscription Details ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Price:" | faint }}	{{ .Price }}
{{ "Description:" | faint }}	{{ .Description }}`,
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

	displaySubscriptionDetails(selectedSub, api.Subscription{}, cfg.IsAuthenticated())

	// Ask if user wants to subscribe (if authenticated)
	if cfg.IsAuthenticated() {
		fmt.Println()
		confirmed, err := promptConfirm(fmt.Sprintf("Would you like to subscribe to %s now", selectedSub.Name))
		if err == nil && confirmed {
			// User wants to subscribe - start order configuration flow
			return createOrderAndSubscribe(cfg, client, selectedSub)
		}
	} else if !cfg.IsAuthenticated() {
		fmt.Println("\nPlease login first to subscribe:")
		fmt.Println("  bc-cli login")
	}

	return nil
}

func displaySubscriptionDetails(sub api.AvailableSubscription, activeSub api.Subscription, isAuthenticated bool) {
	type activeSubData struct {
		ID        string
		Status    string
		StartedAt string
		ExpiresAt string
	}

	var activeData activeSubData
	if activeSub.ID != "" {
		activeData.ID = activeSub.ID
		activeData.Status = activeSub.Status
		if activeSub.StartedAt != nil {
			activeData.StartedAt = utils.FormatTimestamp(*activeSub.StartedAt)
		}
		if activeSub.ExpiresAt != nil {
			activeData.ExpiresAt = utils.FormatTimestamp(*activeSub.ExpiresAt)
		}
	}

	if err := templates.RenderToStdout(templates.SubscriptionDetailsTemplate, struct {
		Name          string
		Currency      string
		Price         string
		BillingPeriod string
		Description   string
		Features      []string
		ActiveSub     activeSubData
	}{
		Name:          sub.Name,
		Currency:      sub.Currency,
		Price:         sub.Price,
		BillingPeriod: sub.BillingPeriod,
		Description:   sub.Description,
		Features:      sub.Features,
		ActiveSub:     activeData,
	}); err != nil {
		fmt.Printf("Error rendering template: %v\n", err)
	}
}

func createOrderAndSubscribe(cfg *config.Config, client *api.Client, tier api.AvailableSubscription) error {
	if !cfg.IsAuthenticated() {
		return fmt.Errorf("you must be logged in to subscribe. Please run 'bc-cli login' first")
	}

	if err := templates.RenderToStdout(templates.OrderConfigIntroTemplate, struct {
		MinQuantity int
		MaxQuantity int
	}{
		MinQuantity: cfg.MinQuantityKg,
		MaxQuantity: cfg.MaxQuantityKg,
	}); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	totalQuantity, err := promptQuantityInt("Total quantity per month (kg)", cfg.MinQuantityKg, cfg.MaxQuantityKg, cfg.MinQuantityKg)
	if err != nil {
		return err
	}

	fmt.Printf("\n✓ Total: %d kg per month\n", totalQuantity)

	// Step 2: Ask if they want to split or keep uniform
	if err := templates.RenderToStdout(templates.OrderSplitIntroTemplate, nil); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	wantsSplit, err := promptConfirm("Would you like different grind methods?")
	if err != nil {
		return err
	}

	var lineItems []api.OrderLineItem

	if !wantsSplit {
		// Simple flow - all coffee the same way
		lineItems, err = configureUniformOrder(totalQuantity)
		if err != nil {
			return err
		}
	} else {
		// Complex flow - split into multiple preferences
		lineItems, err = configureLineItems(totalQuantity)
		if err != nil {
			return err
		}
	}

	// Step 3: Show summary and confirm
	if err := showOrderSummary(tier, totalQuantity, lineItems); err != nil {
		return err
	}

	confirmed, err := promptConfirm("Looks good! Proceed to checkout?")
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("\nOrder cancelled.")
		return nil
	}

	// Step 4: Create order via API
	fmt.Print("\nCreating order... ")
	order, err := client.CreateOrder(api.CreateOrderRequest{
		Tier:            tier.Tier,
		TotalQuantityKg: totalQuantity,
		LineItems:       lineItems,
	})
	if err != nil {
		fmt.Println("✗")
		return fmt.Errorf("failed to create order: %w", err)
	}
	fmt.Println("✓")

	// Step 5: Create checkout session
	fmt.Print("Opening checkout in your browser... ")
	checkout, err := client.CreateCheckoutSession(order.ID)
	if err != nil {
		fmt.Println("✗")
		return fmt.Errorf("failed to create checkout session: %w", err)
	}
	fmt.Println("✓")

	// Step 6: Open browser
	if err := openBrowser(checkout.CheckoutURL); err != nil {
		fmt.Printf("\nCouldn't open browser automatically. Please visit:\n%s\n", checkout.CheckoutURL)
	}

	fmt.Printf("\nOrder created successfully!\n")
	fmt.Printf("Order ID: %s\n\n", order.ID)

	// Step 7: Wait for payment completion
	fmt.Println("Waiting for payment confirmation...")
	fmt.Println("(You have 5 minutes to complete the payment)")

	subscription, completed := waitForSubscriptionActivation(client, order.ID, 5*60) // 5 minutes

	if completed && subscription != nil {
		// Payment successful!
		if err := templates.RenderToStdout(templates.SuccessArtTemplate, nil); err != nil {
			fmt.Printf("Error rendering template: %v\n", err)
		}
		if err := templates.RenderToStdout(templates.SuccessMessageTemplate, struct {
			TotalQuantity int
			TierName      string
		}{
			TotalQuantity: totalQuantity,
			TierName:      tier.Name,
		}); err != nil {
			fmt.Printf("Error rendering template: %v\n", err)
		}
	} else {
		// Timeout or user didn't complete payment
		fmt.Println("Complete your payment to activate your subscription.")
		fmt.Println("Your order will be processed once payment is received.")
	}

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

// Helper functions for order configuration

func promptQuantityInt(label string, min, max, defaultVal int) (int, error) {
	validate := func(input string) error {
		var val int
		_, err := fmt.Sscanf(input, "%d", &val)
		if err != nil {
			return fmt.Errorf("please enter a valid whole number")
		}
		if val < min || val > max {
			return fmt.Errorf("quantity must be between %d and %d kg", min, max)
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Default:  fmt.Sprintf("%d", defaultVal),
	}

	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	var quantity int
	if _, err := fmt.Sscanf(result, "%d", &quantity); err != nil {
		return 0, fmt.Errorf("invalid quantity: %w", err)
	}
	return quantity, nil
}

func promptConfirm(label string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
		Default:   "y",
	}

	result, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return false, nil
		}
		return false, err
	}

	// If empty (user just pressed Enter), default to yes
	result = strings.TrimSpace(result)
	if result == "" {
		return true, nil
	}

	return strings.ToLower(result) == "y" || strings.ToLower(result) == "yes", nil
}

func configureUniformOrder(totalQuantity int) ([]api.OrderLineItem, error) {
	if err := templates.RenderToStdout(templates.UniformOrderIntroTemplate, struct{ TotalQuantity int }{TotalQuantity: totalQuantity}); err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}

	// Prompt for grind type
	grindType, err := selectGrindType()
	if err != nil {
		return nil, err
	}

	// Show confirmation based on choice
	fmt.Println()
	if grindType == "whole_bean" {
		fmt.Println("✓ You'll grind these beans yourself")
	} else {
		fmt.Println("✓ We'll grind these beans for you!")
	}

	// ALWAYS prompt for brewing method
	brewingMethod, err := selectBrewingMethod(grindType)
	if err != nil {
		return nil, err
	}

	// Confirmation message
	fmt.Printf("\n✓ Perfect! All %d kg will be ", totalQuantity)
	if grindType == "whole_bean" {
		fmt.Printf("whole beans, roasted for %s.\n", brewingMethodDisplay(brewingMethod))
	} else {
		fmt.Printf("ground for %s.\n", brewingMethodDisplay(brewingMethod))
	}
	fmt.Println()

	// Create single line item with full quantity
	lineItems := []api.OrderLineItem{
		{
			QuantityKg:    totalQuantity,
			GrindType:     grindType,
			BrewingMethod: brewingMethod,
		},
	}

	return lineItems, nil
}

func configureLineItems(totalQuantity int) ([]api.OrderLineItem, error) {
	var lineItems []api.OrderLineItem
	remaining := totalQuantity
	preferenceNum := 1

	// Introduction
	if err := templates.RenderToStdout(templates.SplitOrderIntroTemplate, struct{ TotalQuantity int }{TotalQuantity: totalQuantity}); err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}

	for remaining > 0 {
		// Show preference header with remaining amount
		lowRemaining := float64(remaining) < float64(totalQuantity)*0.3
		if err := templates.RenderToStdout(templates.PreferenceHeaderTemplate, struct {
			PreferenceNum int
			TotalQuantity int
			Remaining     int
			LowRemaining  bool
		}{
			PreferenceNum: preferenceNum,
			TotalQuantity: totalQuantity,
			Remaining:     remaining,
			LowRemaining:  lowRemaining,
		}); err != nil {
			return nil, fmt.Errorf("failed to render template: %w", err)
		}

		// Prompt for quantity with smart defaults
		maxQty := remaining
		defaultQty := remaining
		if remaining > 2 {
			defaultQty = min(2, remaining)
		}

		quantity, err := promptQuantityInt("  How much for this preference? (kg)", 1, maxQty, defaultQty)
		if err != nil {
			return nil, err
		}

		// Show allocation confirmation
		if quantity >= remaining {
			fmt.Printf("\n✓ Allocating %d kg (this will complete your order!)\n\n", quantity)
		} else {
			fmt.Printf("\n✓ Allocating %d kg\n\n", quantity)
		}

		// Prompt for grind type with explanation
		fmt.Printf("  How would you like these %d kg prepared?\n\n", quantity)
		grindType, err := selectGrindType()
		if err != nil {
			return nil, err
		}

		// Show grind type confirmation
		fmt.Println()
		if grindType == "ground" {
			fmt.Println("✓ We'll grind these beans for you!")
		} else {
			fmt.Println("✓ You'll grind these beans yourself")
		}

		// Prompt for brewing method (ALWAYS, regardless of grind type)
		brewingMethod, err := selectBrewingMethod(grindType)
		if err != nil {
			return nil, err
		}

		lineItems = append(lineItems, api.OrderLineItem{
			QuantityKg:    quantity,
			GrindType:     grindType,
			BrewingMethod: brewingMethod,
		})

		// Updated confirmation message
		fmt.Printf("\n✓ Added: %d kg ", quantity)
		if grindType == "whole_bean" {
			fmt.Printf("whole beans for %s", brewingMethodDisplay(brewingMethod))
		} else {
			fmt.Printf("ground for %s", brewingMethodDisplay(brewingMethod))
		}
		fmt.Println()

		remaining -= quantity

		// Show progress bar
		showProgressBar(totalQuantity-remaining, totalQuantity)

		preferenceNum++

		// Check if we're done
		if remaining <= 0 {
			break
		}
	}

	// Success message
	fmt.Println("\n" + strings.Repeat("─", 60) + "\n")
	fmt.Printf("🎉 Perfect! You've allocated all %d kg!\n\n", totalQuantity)

	return lineItems, nil
}

func selectGrindType() (string, error) {
	items := []struct {
		Value   string
		Display string
	}{
		{"whole_bean", "Whole Bean (I'll grind it myself)"},
		{"ground", "Ground (We'll grind it for you)"},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "  ▸ {{ .Display | cyan }}",
		Inactive: "    {{ .Display }}",
		Selected: "✓ {{ .Display }}",
	}

	prompt := promptui.Select{
		Label:     "  Grind type",
		Items:     items,
		Templates: templates,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return items[idx].Value, nil
}

func selectBrewingMethod(grindType string) (string, error) {
	// Show helpful message first
	fmt.Println("  What is your preferred brewing method?")
	fmt.Println("  This helps us understand the best profiles to ensure the best tasting experience!")

	items := []struct {
		Value       string
		Display     string
		Description string
	}{
		{"espresso", "Espresso", "very fine grind"},
		{"moka", "Moka Pot", "fine-medium grind"},
		{"v60", "V60 Pour Over", "medium grind"},
		{"french_press", "French Press", "coarse grind"},
		{"pour_over", "Pour Over", "medium grind"},
		{"drip", "Drip Coffee", "medium grind"},
		{"cold_brew", "Cold Brew", "extra coarse grind"},
	}

	// Update labels based on grind type
	var label string
	var templates *promptui.SelectTemplates

	if grindType == "whole_bean" {
		label = "  Select your brewing method"
		// For whole beans, show brewing methods without grind descriptions
		templates = &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "  ▸ {{ .Display | cyan }}",
			Inactive: "    {{ .Display }}",
			Selected: "✓ {{ .Display }}",
		}
	} else {
		// For ground coffee, show grind descriptions
		label = "  Select your brewing method"
		templates = &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "  ▸ {{ .Display | cyan }} ({{ .Description | faint }})",
			Inactive: "    {{ .Display }} ({{ .Description | faint }})",
			Selected: "✓ {{ .Display }}",
		}
	}

	prompt := promptui.Select{
		Label:     label,
		Items:     items,
		Templates: templates,
		Size:      8,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return items[idx].Value, nil
}

func showProgressBar(current, total int) {
	if err := templates.RenderToStdout(templates.ProgressBarTemplate, struct {
		Current int
		Total   int
	}{
		Current: current,
		Total:   total,
	}); err != nil {
		fmt.Printf("Error rendering progress bar: %v\n", err)
	}
}

func showOrderSummary(tier api.AvailableSubscription, totalQuantity int, lineItems []api.OrderLineItem) error {
	// Calculate price based on quantity
	// tier.Price is the price per 1kg
	pricePerKg, _ := strconv.ParseFloat(tier.Price, 64)
	totalPrice := pricePerKg * float64(totalQuantity)

	// Format line items for display
	formattedItems := make([]string, len(lineItems))
	for i, item := range lineItems {
		if item.GrindType == "whole_bean" {
			formattedItems[i] = fmt.Sprintf("%d kg → Whole beans for %s",
				item.QuantityKg,
				brewingMethodDisplay(item.BrewingMethod))
		} else {
			grindDesc := getGrindDescription(item.BrewingMethod)
			formattedItems[i] = fmt.Sprintf("%d kg → Ground for %s (%s)",
				item.QuantityKg,
				brewingMethodDisplay(item.BrewingMethod),
				grindDesc)
		}
	}

	return templates.RenderToStdout(templates.OrderSummaryTemplate, struct {
		TierName      string
		TotalQuantity int
		Currency      string
		TotalPrice    float64
		BillingPeriod string
		LineItems     []string
	}{
		TierName:      tier.Name,
		TotalQuantity: totalQuantity,
		Currency:      tier.Currency,
		TotalPrice:    totalPrice,
		BillingPeriod: tier.BillingPeriod,
		LineItems:     formattedItems,
	})
}

func brewingMethodDisplay(method string) string {
	displays := map[string]string{
		"espresso":     "Espresso",
		"moka":         "Moka Pot",
		"v60":          "V60 Pour Over",
		"french_press": "French Press",
		"pour_over":    "Pour Over",
		"drip":         "Drip Coffee",
		"cold_brew":    "Cold Brew",
	}
	if display, ok := displays[method]; ok {
		return display
	}
	return method
}

func getGrindDescription(method string) string {
	descriptions := map[string]string{
		"espresso":     "very fine",
		"moka":         "fine-medium",
		"v60":          "medium",
		"french_press": "coarse",
		"pour_over":    "medium",
		"drip":         "medium",
		"cold_brew":    "extra coarse",
	}
	if desc, ok := descriptions[method]; ok {
		return desc
	}
	return ""
}

// waitForSubscriptionActivation polls the API for subscription activation
func waitForSubscriptionActivation(client *api.Client, orderID string, timeoutSeconds int) (*api.Subscription, bool) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	timeout := time.After(time.Duration(timeoutSeconds) * time.Second)
	dots := 0

	for {
		select {
		case <-ticker.C:
			// Poll for order status
			order, err := client.GetOrder(orderID)
			if err == nil && order.Status == "paid" {
				// Order is paid, fetch subscription
				subscriptions, err := client.ListSubscriptions()
				if err == nil && len(subscriptions) > 0 {
					// Find the active subscription
					for _, sub := range subscriptions {
						if sub.Status == "active" {
							return &sub, true
						}
					}
				}
			}

			// Show progress dots
			dots++
			if dots > 3 {
				dots = 1
			}
			fmt.Printf("\rChecking payment status%s   ", strings.Repeat(".", dots))

		case <-timeout:
			fmt.Println("\r" + strings.Repeat(" ", 50)) // Clear the line
			return nil, false
		}
	}
}
