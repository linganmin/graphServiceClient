package graphServiceClient

import (
	"context"
	"fmt"
	"github.com/linganmin/zaplog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestClient_GetTokens(t *testing.T) {
	type fields struct {
		clientID     string
		redirectURI  string
		scopes       []string
		clientSecret string
	}
	type args struct {
		ctx      context.Context
		authCode string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TokenResponse
		wantErr bool
	}{
		{name: testing.CoverMode(), fields: fields{
			clientID:     "",
			redirectURI:  "http://localhost",
			scopes:       []string{"Mail.Send", "offline_access"},
			clientSecret: "",
		}, args: args{
			ctx:      context.Background(),
			authCode: "M.C533_BL2.2.U.570641d2-1ca3-1622-ea36-21671651f114",
		}, want: nil, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				ClientID:     tt.fields.clientID,
				RedirectURI:  tt.fields.redirectURI,
				Scopes:       tt.fields.scopes,
				ClientSecret: tt.fields.clientSecret,
			}
			got, err := c.GetTokens(tt.args.ctx, tt.args.authCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTokens() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_RefreshAccessToken(t *testing.T) {
	type fields struct {
		clientID     string
		redirectURI  string
		scopes       []string
		clientSecret string
	}
	type args struct {
		ctx          context.Context
		refreshToken string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TokenResponse
		wantErr bool
	}{
		{name: testing.CoverMode(), fields: fields{
			clientID:     "",
			redirectURI:  "http://localhost",
			scopes:       []string{"Mail.Send", "offline_access"},
			clientSecret: "",
		}, args: args{
			ctx:          context.Background(),
			refreshToken: "",
		}, want: nil, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				ClientID:     tt.fields.clientID,
				RedirectURI:  tt.fields.redirectURI,
				Scopes:       tt.fields.scopes,
				ClientSecret: tt.fields.clientSecret,
			}
			got, err := c.RefreshAccessToken(tt.args.ctx, tt.args.refreshToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RefreshAccessToken() got = %v, want %v", got, tt.want)
			}

			//_SendMail(got.AccessToken)
			_SendMailWithAttachments(got.AccessToken)

		})
	}
}

func _SendMailWithAttachments(accessToken string) {
	// 创建邮件请求
	message := models.NewMessage()
	subject := "邮件主题"
	message.SetSubject(&subject)

	body := models.NewItemBody()
	contentType := models.TEXT_BODYTYPE
	body.SetContentType(&contentType)
	content := "这是邮件的内容"
	body.SetContent(&content)
	message.SetBody(body)

	recipient := models.NewRecipient()
	emailAddress := models.NewEmailAddress()
	address := "saboran@163.com"
	emailAddress.SetAddress(&address)
	recipient.SetEmailAddress(emailAddress)

	toRecipients := []models.Recipientable{
		recipient,
	}
	message.SetToRecipients(toRecipients)

	requestBody := users.NewItemSendMailPostRequestBody()
	requestBody.SetMessage(message)
	saveToSentItems := true
	requestBody.SetSaveToSentItems(&saveToSentItems)

	// 附件

	fileBytes, err := os.ReadFile("examples/attachments.xlsx")
	if err != nil {
		zaplog.FromContext(context.Background()).Errorf("读取文件失败:%+v", err)
		return
	}
	// 创建附件
	attachment := models.NewFileAttachment()
	filename := "attachment.xlsx"
	attachment.SetName(&filename)
	attachment.SetContentBytes(fileBytes)
	// 添加附件
	message.SetAttachments([]models.Attachmentable{attachment})

	// 获取客户端
	client := BuildGraphClient(context.Background(), accessToken)
	// 发送邮件
	err = client.Me().SendMail().Post(context.Background(), requestBody, nil)
	if err != nil {
		zaplog.FromContext(context.Background()).Errorf("发送邮件失败:%+v", err)
		return
	}
	fmt.Println("邮件发送成功")

}

func _SendMail(accessToken string) {
	// 创建邮件请求
	message := models.NewMessage()
	subject := "邮件主题"
	message.SetSubject(&subject)

	body := models.NewItemBody()
	contentType := models.TEXT_BODYTYPE
	body.SetContentType(&contentType)
	content := "这是邮件的内容"
	body.SetContent(&content)
	message.SetBody(body)

	recipient := models.NewRecipient()
	emailAddress := models.NewEmailAddress()
	address := "saboran@163.com"
	emailAddress.SetAddress(&address)
	recipient.SetEmailAddress(emailAddress)

	toRecipients := []models.Recipientable{
		recipient,
	}
	message.SetToRecipients(toRecipients)

	requestBody := users.NewItemSendMailPostRequestBody()
	requestBody.SetMessage(message)
	saveToSentItems := true
	requestBody.SetSaveToSentItems(&saveToSentItems)

	// 获取客户端
	client := BuildGraphClient(context.Background(), accessToken)
	// 发送邮件
	err := client.Me().SendMail().Post(context.Background(), requestBody, nil)
	if err != nil {
		log.Fatalf("发送邮件失败: %v", err)
		return
	}
	fmt.Println("邮件发送成功")

}
