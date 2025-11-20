package api

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code,omitempty"` // Optional: omit from request if empty
}

type RegisterResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func (c *Client) Register(req RegisterRequest) (*RegisterResponse, error) {
	resp, err := c.doRequest("POST", "/api/core/v1/users", req, false)
	if err != nil {
		return nil, err
	}

	var result RegisterResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	// Save tokens automatically after successful registration
	c.Config.AccessToken = result.AccessToken
	c.Config.RefreshToken = result.RefreshToken

	if err := c.Config.Save(); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) Login(req LoginRequest) (*LoginResponse, error) {
	resp, err := c.doRequest("POST", "/api/core/v1/users/token", req, false)
	if err != nil {
		return nil, err
	}

	var result LoginResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	c.Config.AccessToken = result.Access
	c.Config.RefreshToken = result.Refresh

	if err := c.Config.Save(); err != nil {
		return nil, err
	}

	return &result, nil
}
