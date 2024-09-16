package oauth

import (
	"context"
	"encoding/json"
	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/services/config"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"io"
)

// Enterprise Only huehue

type IOAuth interface {
	// GetOAuthURL returns the URL to redirect the user to for OAuth.
	GetOAuthURL(url string) string
	ExchangeCode(code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (*model.Profile, error)
}

type OAuth struct {
	config    *oauth2.Config
	appConfig *config.Configuration
}

func New(config *config.Configuration) (*OAuth, error) {
	var scopes = []string{"openid", "profile", "email"}
	if config.OAuthClientID == "" || config.OAuthClientSecret == "" || config.OAuthCallbackURL == "" || config.OAuthAuthURL == "" || config.OAuthTokenURL == "" {
		return nil, errors.New("OAuth configuration is missing")
	}
	if config.OAuthScope != "" {
		scopes = []string{config.OAuthScope}
	}
	var cfg = &oauth2.Config{
		ClientID:     config.OAuthClientID,
		ClientSecret: config.OAuthClientSecret,
		RedirectURL:  config.OAuthCallbackURL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.OAuthAuthURL,
			TokenURL: config.OAuthTokenURL,
		},
	}
	return &OAuth{config: cfg, appConfig: config}, nil
}

func (o *OAuth) GetAuthoriseURL() string {
	url := o.config.AuthCodeURL("state", oauth2.AccessTypeOnline)
	return url
}

func (o *OAuth) ExchangeCode(code string) (*oauth2.Token, error) {
	token, err := o.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (o *OAuth) GetUserInfo(token *oauth2.Token) (*model.Profile, error) {
	if o.appConfig.OAuthProfileURL == "" {
		return nil, errors.New("Cannot get profile without Profile Path")
	}
	client := o.config.Client(context.Background(), token)
	profileResp, err := client.Get(o.appConfig.OAuthProfileURL)
	if err != nil {
		return nil, err
	}
	if profileResp.StatusCode >= 300 {
		return nil, errors.New("Error getting profile")
	}

	defer profileResp.Body.Close()
	parsedBody, err := io.ReadAll(profileResp.Body)
	if err != nil {
		return nil, errors.New("Error parsing profile")
	}

	var profile model.Profile
	unmarshalErr := json.Unmarshal(parsedBody, &profile)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return &profile, nil
}
