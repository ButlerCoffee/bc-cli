package templates

const ActiveSubscriptionsTemplate = `=== Your Active Subscriptions ===

{{range .Subscriptions}}{{if eq .Status "active"}}┌─ {{.Tier | upper}} ✓
│  Status: {{.Status | upper}}
{{if .StartedAt}}│  Started: {{.StartedAt}}
{{end}}{{if .ExpiresAt}}│  Expires: {{.ExpiresAt}}
{{end}}└─ ID: {{.ID}}

{{end}}{{end}}{{if .HasActive}}✓ = Active subscription

{{end}}{{repeat "=" 60}}

`

const SubscriptionDetailsTemplate = `
{{repeat "=" 60}}
{{.Name}}
{{repeat "=" 60}}

Price: {{.Currency}} {{.Price}}/{{.BillingPeriod}}
Description: {{.Description}}
{{if .ActiveSub.ID}}
Status: {{.ActiveSub.Status | upper}}{{if eq .ActiveSub.Status "active"}} ✓{{end}}
{{if .ActiveSub.StartedAt}}Started: {{.ActiveSub.StartedAt}}
{{end}}{{if .ActiveSub.ExpiresAt}}Expires: {{.ActiveSub.ExpiresAt}}
{{end}}{{end}}
Features:
{{range .Features}}  • {{.}}
{{end}}
`

const OrderConfigIntroTemplate = `
{{repeat "─" 60}}
Let's configure your coffee order!
{{repeat "─" 60}}

How much coffee would you like per month?
You can order anywhere from 1 kg to 50 kg.
`

const OrderSplitIntroTemplate = `{{repeat "─" 60}}

Would you like your coffee prepared different ways?
For example, you could get:
  • 2 kg whole bean + 3 kg ground for espresso
  • 2 kg ground for moka + 2 kg ground for v60 + 1 kg whole bean

Or keep it simple with everything the same way.`

const UniformOrderIntroTemplate = `{{repeat "─" 60}}

Great! Let's prepare all {{.TotalQuantity}} kg the same way.

`

const SplitOrderIntroTemplate = `{{repeat "─" 60}}

Great! Now let's split your {{.TotalQuantity}} kg into different
grinding preferences. You can have:
  • Whole beans (you grind at home)
  • Pre-ground for specific brewing methods
We'll help you allocate all {{.TotalQuantity}} kg across your preferences.`

const PreferenceHeaderTemplate = `{{repeat "─" 60}}
┌─ Preference #{{.PreferenceNum}} ─────────────────────────────────────┐
│ Allocating from: {{.TotalQuantity}} kg total{{printf "%-24s" ""}}│
│ Remaining: {{.Remaining}} kg{{if .LowRemaining}} ⚠️  (almost done!){{end}}{{printf "%-18s" ""}}│
└──────────────────────────────────────────────────────┘
`

const ProgressBarTemplate = `
┌──────────────────────────────────────────────────────┐
│ Progress: {{progressBar .Current .Total 30}} {{.Current}}/{{.Total}} kg{{if ge .Current .Total}} ✓{{end}}{{printf "%-5s" ""}}│
└──────────────────────────────────────────────────────┘`

const OrderSummaryTemplate = `Your Order Summary:
┌────────────────────────────────────────────────────────┐
│ Tier: {{printf "%-48s" .TierName}} │
│ Total: {{.TotalQuantity}} kg/month{{printf "%-38s" ""}} │
│ Price: {{.Currency}} {{printf "%.2f" .TotalPrice}}/{{.BillingPeriod}}{{printf "%-36s" ""}} │
│                                                         │
│ How your coffee will be prepared:                     │
{{range $i, $item := .LineItems}}│ {{printf "%-54s" (printf "   %d. %s" (add $i 1) $item)}} │
{{end}}└────────────────────────────────────────────────────────┘
`

const CheckoutHeaderTemplate = `
{{repeat "─" 60}}
Opening checkout...
`

const SuccessMessageTemplate = `
🎉 Congratulations! Your subscription is now active!

📦 Your first shipment of {{.TotalQuantity}} kg of fresh {{.TierName}} coffee
   will be shipped within the next 7 days.

☕ Get ready for an amazing coffee experience!
`

const SuccessArtTemplate = `
MMMMMMMMMMMMMWXOdc;;;cOWMMMMMMMMMMMMMMMM
MMMMMMMMMMMXxc,...''..'xWMMMMMMMMMMMMMMM
MMMMMMMMMMXc.......,,'.'xNX0OKWMMMMMMMMM
MMMMMMMMMMNo.......;cc:''::,,;kWMMMMMMMM
MMMMMMMMMMMXl..';;:cc:,'',;,,oKMMMMMMMMM
MMMMMMMMMMMW0;.,,'.''';:cdxdlxNMMMMMMMMM
MMMMMMMMMWKo;...';clodxxdxxoc:dKX0O0NMMM
MMMMMMMMMWd....:okO000Oc':dloxdxxxdl0MMM
MWWMMMMMMMNOxxlokO00000OxkxccooxxddkXMMM
XolKWMMMMMMMWKc;oxO00KK0KKOdc:dO00kd0WMM
d..,oOKNNNXOo,...';coddlcdxl,oNMNOdokXWM
c.....'::;'..',......,:...'..cXXo:llccdK
l......;,.....;:;,.....':dxl,oNK:.....:d
k'.....;'.......',;;,..;0N0Kxd0x:'..',;k
Nd......;..........;o:..xXO0l''.,cooodxK
MNd'....,,.........;l;..dNKd....'xXNNWMM
MMWO:....'''.....',,'...okl....,xNMMMMMM
MMMMNkc'........''.....'cooolokXWMMMMMMM
MMMMMMWKko:,.......,;;cx0WMMMMMMMMMMMMMM
MMMMMMMMMMNKOxdddxk0KXWMMMMMMMMMMMMMMMMM
`
