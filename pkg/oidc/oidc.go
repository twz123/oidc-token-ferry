package oidc

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	oidc "github.com/coreos/go-oidc"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type Config struct {
	IssuerURL, ClientID, ClientSecret string
}

type OIDCFlow struct {
	context      context.Context
	oidcProvider *oidc.Provider
	oidcVerifier *oidc.IDTokenVerifier
	oauth2Config *oauth2.Config
}

func NewOpenIDConnectFlow(config *Config) (*OIDCFlow, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, config.IssuerURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create OpenID Connect provider %s", config.IssuerURL)
	}

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: config.ClientID})

	return &OIDCFlow{
		context:      ctx,
		oidcProvider: provider,
		oidcVerifier: verifier,
		oauth2Config: oauth2Config,
	}, nil
}

type OIDCChallenge struct {
	flow  *OIDCFlow
	state string
}

func (flow *OIDCFlow) InitiateChallenge() (*OIDCChallenge, error) {
	state, err := generateRandomState()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate random state")
	}

	return &OIDCChallenge{
		flow:  flow,
		state: state,
	}, nil
}

func (challenge *OIDCChallenge) RedirectURL() string {
	return challenge.flow.oauth2Config.AuthCodeURL(challenge.state)
}

func (challenge *OIDCChallenge) ExchangeCode(code string) (*Credentials, error) {

	oauth2Token, err := challenge.flow.oauth2Config.Exchange(challenge.flow.context, code)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exchange code")
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, errors.Wrap(err, "no ID Token in server response")
	}

	// Parse and verify ID Token payload.
	idToken, err := challenge.flow.oidcVerifier.Verify(challenge.flow.context, rawIDToken)
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify ID Token")
	}

	// Extract custom claims
	var claims struct {
		EMail         string `json:"email"`
		EMailVerified bool   `json:"email_verified"`
	}

	if err := idToken.Claims(&claims); err != nil {
		return nil, errors.Wrap(err, "failed to parse claims")
	}

	return &Credentials{
		IDToken:      rawIDToken,
		RefreshToken: oauth2Token.RefreshToken,
		EMail:        claims.EMail,
	}, nil
}

type Credentials struct {
	IDToken      string
	RefreshToken string
	EMail        string
}

func generateRandomState() (string, error) {
	bytes := make([]byte, 18)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}
