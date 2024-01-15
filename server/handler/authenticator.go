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
}

const ID_SCOPE string = "api.ansible_on_cloud"

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
		Scopes:       []string{oidc.ScopeOpenID, ID_SCOPE, "profile", "email"}, // TODO Add scopes for needed data
	}
	return &Authenticator{
		Config:   oauth2Config,
		provider: provider,
		ctx:      ctx,
	}, nil
}

func (a *Authenticator) VerifyToken(token *oauth2.Token, clientId string) (*oidc.IDToken, error) {
	// The returned token contains an access token and an ID token
	// The access token is a JWT and contains the data we need
	verifier := a.provider.VerifierContext(a.ctx, &oidc.Config{ClientID: ID_SCOPE})
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
	accessToken.Claims(&claims)
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
			log.Warnf("Not a string: %v", orgMap["id"])
		}
		session.OrganizationId = id
	}
	return session
}
