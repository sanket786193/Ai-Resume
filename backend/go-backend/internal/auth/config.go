package auth

import "resume/internal/config"

// Ensure config.AuthConfig is compatible with auth.AuthConfig.
type configAdapter struct {
	*config.Config
}

// AuthConfigAdapter wraps app config for auth service.
func AuthConfigAdapter(c *config.Config) AuthConfig {
	return &configAdapter{Config: c}
}

func (c *configAdapter) GetJWTSecret() string {
	return c.Auth.JWTSecret
}

func (c *configAdapter) GetJWTExpiryHours() int {
	return c.Auth.JWTExpiryHours
}

func (c *configAdapter) GetRefreshExpiryDays() int {
	return c.Auth.RefreshExpiryDays
}
