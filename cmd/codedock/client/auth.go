package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"codedock/internal/models"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TOTPCode string `json:"totpCode,omitempty"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (c *Client) Login(email, password string) (*AuthResponse, error) {
	return c.LoginWithTOTP(email, password, "")
}

func (c *Client) LoginWithTOTP(email, password, totpCode string) (*AuthResponse, error) {
	payload := AuthRequest{
		Email:    email,
		Password: password,
		TOTPCode: totpCode,
	}

	resp, err := c.sendRequest("POST", "/auth/signin", payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed: %s", string(respBody))
	}

	var result struct {
		Data *AuthResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Data == nil {
		return nil, fmt.Errorf("empty response data")
	}

	return result.Data, nil
}

func (c *Client) Logout() error {
	resp, err := c.sendRequest("POST", "/auth/logout", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("logout failed: %s", string(respBody))
	}

	return nil
}

func (c *Client) VerifyToken(token string) (*models.User, error) {
	c.Token = token
	return c.Me()
}

func (c *Client) ExchangeLocalCLISession(internalToken string) (*AuthResponse, error) {
	headers := map[string]string{
		"X-Internal-Token": internalToken,
	}
	resp, err := c.sendRequestWithHeaders("POST", "/system/local-cli-session", nil, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("local CLI session failed (%d): %s", resp.StatusCode, string(respBody))
	}

	var sessionResp struct {
		Token     string `json:"token"`
		ExpiresAt string `json:"expiresAt"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&sessionResp); err != nil {
		return nil, fmt.Errorf("failed to decode local session response: %w", err)
	}

	c.Token = sessionResp.Token
	user, err := c.Me()
	if err != nil {
		return nil, fmt.Errorf("retrieving user profile for session: %w", err)
	}

	return &AuthResponse{
		Token: sessionResp.Token,
		User:  user,
	}, nil
}

