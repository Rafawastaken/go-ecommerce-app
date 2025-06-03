package notifications

import (
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"go-ecommerce-app/config"
)

type NotificationClient interface {
	SendSMS(phone string, message string) error
}

type notificationClient struct {
	config config.AppConfig
}

// Twilio

func (c notificationClient) SendSMS(phone string, message string) error {
	twilioAccountSid := c.config.TwilioAccountSid
	twilioAuthToken := c.config.TwilioAuthToken
	twilioPhoneNumber := c.config.TwilioPhoneNumber

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: twilioAccountSid,
		Password: twilioAuthToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(phone)
	params.SetFrom(twilioPhoneNumber)
	params.SetBody(message)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return err
	}

	return nil
}

func NewNotificationClient(config config.AppConfig) NotificationClient {
	return &notificationClient{
		config: config,
	}
}
