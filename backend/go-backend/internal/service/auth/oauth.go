package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"resume/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/microsoft"
)

type OAuthService struct {
	config *config.AuthConfig
}

func NewOAuthService(cfg *config.AuthConfig) *OAuthService {
	return &OAuthService{config: cfg}
}

func (o *OAuthService) GetGoogleConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     o.config.Google.ClientID,
		ClientSecret: o.config.Google.ClientSecret,
		RedirectURL:  o.config.RedirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
}

func (o *OAuthService) GetGitHubConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     o.config.GitHub.ClientID,
		ClientSecret: o.config.GitHub.ClientSecret,
		RedirectURL:  o.config.RedirectURL,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}

func (o *OAuthService) GetMicrosoftConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     o.config.Microsoft.ClientID,
		ClientSecret: o.config.Microsoft.ClientSecret,
		RedirectURL:  o.config.RedirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     microsoft.AzureADEndpoint("common"),
	}
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type GitHubUserInfo struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type MicrosoftUserInfo struct {
	ID                string `json:"id"`
	Mail              string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"`
	DisplayName       string `json:"displayName"`
}

func (o *OAuthService) GetGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))

	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
