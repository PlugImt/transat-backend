package services

import (
	"fmt"
	"os"

	"github.com/plugimt/transat-backend/utils"
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
		// si besoin de gérer les réponses, on peut rajouter un truc pour faire un callback, mais flemme
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
			Code:   language,
			Policy: "deterministic",
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
	res, err := s.SendTemplateMessage("verify_code", phoneNumber, language, []components.TemplateMessageComponent{
		components.TemplateMessageComponentBodyType{
			Type: components.TemplateMessageComponentTypeBody,
			Parameters: []components.TemplateMessageParameter{
				components.TemplateMessageBodyAndHeaderParameter{
					Type: "text",
					Text: &code,
				},
			},
		},
		components.TemplateMessageComponentButtonType{
			Type:    components.TemplateMessageComponentTypeButton,
			SubType: components.TemplateMessageButtonComponentTypeUrl,
			Index:   0,
			Parameters: &[]components.TemplateMessageParameter{
				components.TemplateMessageBodyAndHeaderParameter{
					Type: "text",
					Text: &code,
				},
			},
		},
	})

	if err != nil {
		return err
	}

	utils.LogMessage(utils.LevelInfo, "Message sent successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Message", res.Messages[0].ID)

	return nil
}
