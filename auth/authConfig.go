package auth

import (
	"fmt"
	"net/url"
)

var DefaultConfig = AuthorizationConfig{
	RedirectPort: "8080",
	RedirectPath: "/myapp",
	Scope:        "https://management.azure.com/.default",
	OpenCMD:      "open",
}

type AuthorizationConfig struct {
	RedirectPort string
	RedirectPath string
	Scope        string
	ClientID     string
	OpenCMD      string
	ClientSecret string
}

// RedirectURL )
func (c AuthorizationConfig) RedirectURL() string {
	host := "localhost"
	if c.RedirectPort != "" {
		host = fmt.Sprintf("%s:%s", host, c.RedirectPort)
	}
	uri := url.URL{
		Host:   host,
		Scheme: "http",
		Path:   c.RedirectPath,
	}

	return uri.String()
}
