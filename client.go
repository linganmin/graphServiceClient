package graphServiceClient

import (
	"context"
	"github.com/linganmin/zaplog"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

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
