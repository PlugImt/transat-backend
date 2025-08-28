package services

import (
	"fmt"
	"time"

	"github.com/plugimt/transat-backend/utils"
)

type PhoneVerificationService struct {
	whatsappService *WhatsAppService
}

func NewPhoneVerificationService(whatsappService *WhatsAppService) *PhoneVerificationService {
	return &PhoneVerificationService{
		whatsappService: whatsappService,
	}
}

func (s *PhoneVerificationService) SendPhoneVerificationCode(phoneNumber, language string) (string, error) {
	verificationCode := utils.Generate2FACode(6)

	err := s.whatsappService.SendVerificationCode(phoneNumber, verificationCode, language)
	if err != nil {
		return "", fmt.Errorf("failed to send WhatsApp verification code: %w", err)
	}

	utils.LogMessage(utils.LevelInfo, "Phone verification code sent successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Phone", phoneNumber)
	utils.LogLineKeyValue(utils.LevelInfo, "Code", verificationCode)

	return verificationCode, nil
}

func (s *PhoneVerificationService) IsVerificationCodeExpired(expirationTime string) bool {
	expiration, err := time.Parse(time.RFC3339, expirationTime)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse expiration time")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return true // je vois pas comment ça pourrait arriver (null?), mais dans le doute on refuse
	}

	return time.Now().After(expiration)
}

func (s *PhoneVerificationService) ValidateVerificationCode(storedCode, providedCode, expirationTime string) bool {
	if s.IsVerificationCodeExpired(expirationTime) {
		utils.LogMessage(utils.LevelWarn, "Verification code expired")
		return false
	}

	if storedCode != providedCode {
		utils.LogMessage(utils.LevelWarn, "Invalid verification code")
		utils.LogLineKeyValue(utils.LevelWarn, "Stored", storedCode)
		utils.LogLineKeyValue(utils.LevelWarn, "Provided", providedCode)
		return false
	}

	utils.LogMessage(utils.LevelInfo, "Verification code validated successfully")
	return true
}
