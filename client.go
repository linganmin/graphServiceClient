package graphServiceClient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/linganmin/zaplog"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const tenantID = "common" // 普通用户账号固定为common

type Client struct {
	ClientID     string
	RedirectURI  string
	Scopes       []string
	ClientSecret string
}

func New(clientID string, clientSecret string, redirectURI string, scopes []string) *Client {

	return &Client{
		ClientID:     clientID,
		RedirectURI:  redirectURI,
		Scopes:       scopes,
		ClientSecret: clientSecret,
	}
}

func BuildGraphClient(ctx context.Context, accessToken string) *msgraphsdk.GraphServiceClient {
	logger := zaplog.FromContext(ctx)

	authProvider := &TokenAuthProvider{AccessToken: accessToken}

	// 创建请求适配器
	adapter, err := msgraphsdk.NewGraphRequestAdapter(authProvider)
	if err != nil {
		logger.Errorf("failed to create request adapter %+v", err)
		return nil
	}

	client := msgraphsdk.NewGraphServiceClient(adapter)
	return client
}

func (c *Client) GetTokens(ctx context.Context, authCode string) (*TokenResponse, error) {
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

func (c *Client) RefreshAccessToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
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
