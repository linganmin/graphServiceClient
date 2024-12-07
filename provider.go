package graphServiceClient

import (
	"context"
	"fmt"
	abs "github.com/microsoft/kiota-abstractions-go"
)

type TokenAuthProvider struct {
	AccessToken string
}

func (p *TokenAuthProvider) AuthenticateRequest(context context.Context, request *abs.RequestInformation, _ map[string]interface{}) error {
	request.Headers.Add("Authorization", fmt.Sprintf("Bearer %s", p.AccessToken))
	return nil
}
