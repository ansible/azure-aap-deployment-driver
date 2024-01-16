package handler

import (
	"context"
	"fmt"
	"server/model"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	Config   oauth2.Config
	provider *oidc.Provider
	ctx      context.Context
	Audience string
}

type AuthenticatorConfig struct {
	Context      context.Context
	SsoEndpoint  string
	RedirecUrl   string
	ClientId     string
	ClientSecret string
	Scopes       []string
	Audience     string // Should contain expected audience based on scope list
}

func NewAuthenticator(authConfig AuthenticatorConfig) (*Authenticator, error) {
	provider, err := oidc.NewProvider(authConfig.Context, authConfig.SsoEndpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to instantiate new OIDC provider: %v", err)
	}

	oauth2Config := oauth2.Config{
		ClientID:     authConfig.ClientId,
		ClientSecret: authConfig.ClientSecret,
		RedirectURL:  authConfig.RedirecUrl,
		Endpoint:     provider.Endpoint(),
		Scopes:       authConfig.Scopes,
	}
	return &Authenticator{
		Config:   oauth2Config,
		provider: provider,
		ctx:      authConfig.Context,
		Audience: authConfig.Audience,
	}, nil
}

func (a *Authenticator) VerifyToken(token *oauth2.Token) (*oidc.IDToken, error) {
	// The returned token contains an access token and an ID token
	// The access token is a JWT and contains the data we need
	verifier := a.provider.VerifierContext(a.ctx, &oidc.Config{ClientID: a.Audience})
	accessToken, err := verifier.Verify(a.ctx, token.AccessToken)
	if err != nil {
		log.Errorf("Unable to verify SSO access token for user details: %v", err)
		return nil, err
	}
	return accessToken, nil
}

func (a *Authenticator) ExtractUserInfo(accessToken *oidc.IDToken) *model.SsoSession {
	session := &model.SsoSession{}
	claims := jwt.MapClaims{}
	err := accessToken.Claims(&claims)
	if err != nil {
		log.Errorf("Unable to extract claims (user info) from access token: %v", err)
		return session // Empty, will still support storing login, just without user data
	}
	email, ok := claims["email"].(string)
	if !ok {
		log.Warnf("Email not found in access token claims, value: %v", claims["email"])
	} else {
		session.Email = email
	}
	name, ok := claims["name"].(string)
	if !ok {
		log.Warnf("Name not found in access token claims, value: %v", claims["name"])
	} else {
		session.Name = name
	}
	orgMap, ok := claims["organization"].(map[string]interface{})
	if !ok {
		log.Warnf("Organization map not found or not parseable in access token claims, value: %v", claims["organization"])
	} else {
		id, ok := orgMap["id"].(string)
		if !ok {
			log.Warnf("Organization ID is not a string: %v", orgMap["id"])
		}
		session.OrganizationId = id
	}
	return session
}
