package services

import (
	"fmt"
	"os"

	"github.com/wapikit/wapi.go/manager"
	wapi "github.com/wapikit/wapi.go/pkg/client"
	"github.com/wapikit/wapi.go/pkg/components"
	"github.com/wapikit/wapi.go/pkg/messaging"
)

type WhatsAppService struct {
	client        *wapi.Client
	messageClient *messaging.MessagingClient
}

func NewWhatsAppService() (*WhatsAppService, error) {
	apiAccessToken := os.Getenv("WHATSAPP_API_ACCESS_TOKEN")
	businessAccountId := os.Getenv("WHATSAPP_BUSINESS_ACCOUNT_ID")
	phoneNumberId := os.Getenv("WHATSAPP_PHONE_NUMBER_ID")

	if apiAccessToken == "" || businessAccountId == "" || phoneNumberId == "" {
		return nil, fmt.Errorf("missing WhatsApp env vars: WHATSAPP_API_ACCESS_TOKEN, WHATSAPP_BUSINESS_ACCOUNT_ID, or WHATSAPP_PHONE_NUMBER_ID not set")
	}

	client := wapi.New(&wapi.ClientConfig{
		ApiAccessToken:    apiAccessToken,
		BusinessAccountId: businessAccountId,
	})

	messageClient := client.NewMessagingClient(phoneNumberId)

	return &WhatsAppService{
		client:        client,
		messageClient: messageClient,
	}, nil
}

func (s *WhatsAppService) SendTemplateMessage(templateName, phoneNumber string, language string, templateComponents []components.TemplateMessageComponent) (*manager.MessageSendResponse, error) {
	templateMsg := &components.TemplateMessage{
		Name: templateName,
		Language: components.TemplateMessageLanguage{
			Code: language,
		},
		Components: templateComponents,
	}

	res, err := s.messageClient.Message.Send(templateMsg, phoneNumber)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *WhatsAppService) SendVerificationCode(phoneNumber, code string, language string) error {
	_, err := s.SendTemplateMessage("verify_code", phoneNumber, language, []components.TemplateMessageComponent{
		components.TemplateMessageComponentBodyType{
			Type: components.TemplateMessageComponentTypeBody,
			Parameters: []components.TemplateMessageParameter{
				components.TemplateMessageBodyAndHeaderParameter{
					Type:          "text",
					ParameterName: &[]string{"code"}[0],
					Text:          &code,
				},
			},
		},
	})
	return err
}
