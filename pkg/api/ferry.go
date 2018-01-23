package api

import "github.com/twz123/oidc-token-ferry/pkg/oidc"

type TokenFerry struct {
	ClientID, ClientSecret string
	Creds                  *oidc.Credentials
}
