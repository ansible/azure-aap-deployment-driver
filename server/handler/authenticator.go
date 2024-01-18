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
	RedirectUrl  string
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
		RedirectURL:  authConfig.RedirectUrl,
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
	email, emailFound := claims["email"].(string)
	name, nameFound := claims["name"].(string)
	var orgFound bool
	var orgId string
	orgMap, ok := claims["organization"].(map[string]interface{})
	if ok {
		orgId, orgFound = orgMap["id"].(string)
	}
	if !emailFound || !nameFound || !orgFound {
		log.Warnf("Email, name and/or organization not found in access token claims: %v", claims)
	} else {
		session.Name = name
		session.Email = email
		session.OrganizationId = orgId
	}
	return session
}
