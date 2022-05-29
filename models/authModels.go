package models

// AuthorizationCode is a value provided after initial successful
// authentication/authorization, it is used to get access/refresh tokens
type AuthorizationCode struct {
	Value string
}

// Tokens holds access and refresh tokens
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
