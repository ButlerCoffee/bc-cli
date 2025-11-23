package api

type SubscriptionPaymentLinkRequest struct {
	Tier string `json:"tier"`
}

type SubscriptionPaymentLinkResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data struct {
		Tier        string `json:"tier"`
		PaymentLink string `json:"payment_link"`
		Message     string `json:"message,omitempty"`
	} `json:"data"`
}

type SubscriptionPaymentLink struct {
	Tier        string `json:"tier"`
	PaymentLink string `json:"payment_link"`
	Message     string `json:"message,omitempty"`
}

type Subscription struct {
	ID                string  `json:"id"`
	Tier              string  `json:"tier"`
	Status            string  `json:"status"`
	StripePaymentLink string  `json:"stripe_payment_link"`
	StartedAt         *string `json:"started_at"`
	ExpiresAt         *string `json:"expires_at"`
	CreatedOn         string  `json:"created_on"`
}

type AvailableSubscription struct {
	Tier          string   `json:"tier"`
	Name          string   `json:"name"`
	Price         string   `json:"price"`
	Currency      string   `json:"currency"`
	BillingPeriod string   `json:"billing_period"`
	Description   string   `json:"description"`
	Features      []string `json:"features"`
}

func (c *Client) GetSubscriptionPaymentLink(tier string) (*SubscriptionPaymentLink, error) {
	req := SubscriptionPaymentLinkRequest{Tier: tier}
	resp, err := c.doRequest("POST", "/api/core/v1/subscriptions/payment-link", req, true)
	if err != nil {
		return nil, err
	}

	var result SubscriptionPaymentLinkResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &SubscriptionPaymentLink{
		Tier:        result.Data.Tier,
		PaymentLink: result.Data.PaymentLink,
		Message:     result.Data.Message,
	}, nil
}

type ListSubscriptionsResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data []Subscription `json:"data"`
}

func (c *Client) ListSubscriptions() ([]Subscription, error) {
	resp, err := c.doRequest("GET", "/api/core/v1/subscriptions", nil, true)
	if err != nil {
		return nil, err
	}

	var result ListSubscriptionsResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

type AvailableSubscriptionsResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data []AvailableSubscription `json:"data"`
}

func (c *Client) GetAvailableSubscriptions() ([]AvailableSubscription, error) {
	resp, err := c.doRequest("GET", "/api/core/v1/subscriptions/available", nil, false)
	if err != nil {
		return nil, err
	}

	var result AvailableSubscriptionsResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
