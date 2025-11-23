package api

import "fmt"

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code,omitempty"` // Optional: omit from request if empty
}

type RegisterResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data struct {
		ID           string `json:"id"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	} `json:"data"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data struct {
		AccessToken            string `json:"access_token"`
		RefreshToken           string `json:"refresh_token"`
		ExpiresAt              string `json:"expires_at"`
		RefreshTokenExpiresAt  string `json:"refresh_token_expires_at"`
		UserID                 string `json:"user_id"`
	} `json:"data"`
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
	c.Config.AccessToken = result.Data.AccessToken
	c.Config.RefreshToken = result.Data.RefreshToken

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

	c.Config.AccessToken = result.Data.AccessToken
	c.Config.RefreshToken = result.Data.RefreshToken
	c.Config.ExpiresAt = result.Data.ExpiresAt
	c.Config.RefreshTokenExpiresAt = result.Data.RefreshTokenExpiresAt

	if err := c.Config.Save(); err != nil {
		return nil, fmt.Errorf("failed to save config: %w", err)
	}

	return &result, nil
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data struct {
		AccessToken            string `json:"access_token"`
		RefreshToken           string `json:"refresh_token"`
		ExpiresAt              string `json:"expires_at"`
		RefreshTokenExpiresAt  string `json:"refresh_token_expires_at"`
	} `json:"data"`
}

func (c *Client) RefreshToken() error {
	if c.Config.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	req := RefreshTokenRequest{
		RefreshToken: c.Config.RefreshToken,
	}

	resp, err := c.doRequest("POST", "/api/core/v1/users/token/refresh", req, false)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}

	var result RefreshTokenResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return fmt.Errorf("failed to parse refresh token response: %w", err)
	}

	c.Config.AccessToken = result.Data.AccessToken
	c.Config.RefreshToken = result.Data.RefreshToken
	c.Config.ExpiresAt = result.Data.ExpiresAt
	c.Config.RefreshTokenExpiresAt = result.Data.RefreshTokenExpiresAt

	if err := c.Config.Save(); err != nil {
		return fmt.Errorf("failed to save config after refresh: %w", err)
	}

	return nil
}
