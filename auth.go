package graphServiceClient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/linganmin/zaplog"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const tenantID = "common" // 普通用户账号固定为common

type Auth struct {
	ClientID     string
	RedirectURI  string
	Scopes       []string
	ClientSecret string
}

func NewAuth(clientID string, clientSecret string, redirectURI string, scopes []string) *Auth {

	return &Auth{
		ClientID:     clientID,
		RedirectURI:  redirectURI,
		Scopes:       scopes,
		ClientSecret: clientSecret,
	}
}

func (c *Auth) GetTokens(ctx context.Context, authCode string) (*TokenResponse, error) {
	logger := zaplog.FromContext(ctx)
	tokenUrl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)
	data := url.Values{}
	data.Set("client_id", c.ClientID)
	data.Set("scope", strings.Join(c.Scopes, " "))
	data.Set("redirect_uri", c.RedirectURI)
	data.Set("grant_type", "authorization_code")
	data.Set("client_secret", c.ClientSecret)
	data.Set("code", authCode)

	req, err := http.NewRequest("POST", tokenUrl, strings.NewReader(data.Encode()))
	if err != nil {
		logger.Errorf("failed to get tokens %+v", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("failed to get tokens %+v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("failed to get tokens %+v", err)
		return nil, err
	}

	result := &TokenResponse{}
	if err := json.Unmarshal(body, result); err != nil {
		logger.Errorf("failed to get tokens %+v", err)
		return nil, err
	}

	return result, nil
}

func (c *Auth) RefreshAccessToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	logger := zaplog.FromContext(ctx)

	tokenUrl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)
	data := url.Values{}
	data.Set("client_id", c.ClientID)
	data.Set("scope", strings.Join(c.Scopes, " "))
	data.Set("redirect_uri", c.RedirectURI)
	data.Set("grant_type", "refresh_token")
	data.Set("client_secret", c.ClientSecret)
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", tokenUrl, strings.NewReader(data.Encode()))
	if err != nil {
		logger.Errorf("failed to refresh access token %+v", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("failed to refresh access token %+v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("failed to refresh access token %+v", err)
		return nil, err
	}

	result := &TokenResponse{}
	if err := json.Unmarshal(body, result); err != nil {
		logger.Errorf("failed to refresh access token %+v", err)
		return nil, err
	}

	return result, nil
}
