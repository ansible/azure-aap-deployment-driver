package handler

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	Config   oauth2.Config
	provider *oidc.Provider
	ctx      context.Context
}

func NewAuthenticator(ctx context.Context, ssoEndpoint string, redirectUrl string, clientId string, clientSecret string) (*Authenticator, error) {
	provider, err := oidc.NewProvider(ctx, ssoEndpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to instantiate new OIDC provider: %v", err)
	}

	oauth2Config := oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID}, // TODO Add scopes for needed data
	}
	return &Authenticator{
		Config:   oauth2Config,
		provider: provider,
		ctx:      ctx,
	}, nil
}

func (a *Authenticator) VerifyToken(token *oauth2.Token, clientId string) (*oidc.IDToken, error) {
	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("failed to Id token from oauth2 token")
	}

	// Verify the token
	verifier := a.provider.VerifierContext(a.ctx, &oidc.Config{ClientID: clientId})
	idToken, err := verifier.Verify(a.ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %v", err)
	}
	return idToken, err
}
