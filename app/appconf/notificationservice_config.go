package appconf

import "os"

type NotificationServiceConfig struct {
	SendEmailOTPURL           string `validate:"required" name:"NOTIFICATION_SERVICE_SEND_EMAIL_URL"`
	SendSMSOTPURL             string `validate:"required" name:"NOTIFICATION_SERVICE_SEND_SMS_URL"`
	TemplateEmailOTP          string `validate:"required" name:"OTP_EMAIL_TEMPLATE_ID"`
	TemplateSMSOTP            string `validate:"required" name:"OTP_SMS_TEMPLATE_ID"`
	NotificationServiceAPIKey string `validate:"required" name:"NOTIFICATION_SERVICE_API_KEY"`
}

func NotificationServiceConfigInit() *NotificationServiceConfig {
	return &NotificationServiceConfig{
		SendEmailOTPURL:           os.Getenv("NOTIFICATION_SERVICE_SEND_EMAIL_URL"),
		SendSMSOTPURL:             os.Getenv("NOTIFICATION_SERVICE_SEND_SMS_URL"),
		TemplateEmailOTP:          os.Getenv("OTP_EMAIL_TEMPLATE_ID"),
		TemplateSMSOTP:            os.Getenv("OTP_SMS_TEMPLATE_ID"),
		NotificationServiceAPIKey: os.Getenv("NOTIFICATION_SERVICE_API_KEY"),
	}
}
