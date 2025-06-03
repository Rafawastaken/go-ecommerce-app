package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort        string
	Dsn               string
	AppSecret         string
	TwilioAccountSid  string
	TwilioAuthToken   string
	TwilioPhoneNumber string
}

func SetupEnv() (cfg AppConfig, err error) {
	if os.Getenv("APP_ENV") == "dev" {
		godotenv.Load()
	}

	httpPort := os.Getenv("HTTP_PORT")

	if len(httpPort) < 1 {
		return AppConfig{}, errors.New("env variables not found")
	}

	Dsn := os.Getenv("DSN")
	appSecret := os.Getenv("APP_SECRET")
	TwilioAccountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	TwilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	TwilioPhoneNumber := os.Getenv("TWILIO_PHONE_NUMBER")

	if len(Dsn) < 1 {
		return AppConfig{}, errors.New("dsn variables not found")
	}
	if len(appSecret) < 1 {
		return AppConfig{}, errors.New("appSecret variable not found")
	}

	return AppConfig{ServerPort: httpPort, Dsn: Dsn, AppSecret: appSecret, TwilioAccountSid: TwilioAccountSid, TwilioAuthToken: TwilioAuthToken, TwilioPhoneNumber: TwilioPhoneNumber}, nil
}
