package appconf

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type OTPConfig struct {
	EmailOTPExpired       time.Duration `validate:"required" name:"EMAIL_OTP_EXPIRED"`
	PhoneOTPExpired       time.Duration `validate:"required" name:"SMS_OTP_EXPIRED"`
	TemplateSMSOTPRegis   string        `validate:"required" name:"SEND_OTP_SMS_TEMPLATE_REGIS_ID"`
	OTPRegisExpired       time.Duration `validate:"required" name:"SMS_OTP_REGIS_EXPIRED"`
	TemplateEmailOTPRegis string        `validate:"required" name:"SEND_OTP_EMAIL_TEMPLATE_REGIS_ID"`
}

func OTPConfigInit() *OTPConfig {
	emailOTPExpired, err := time.ParseDuration(os.Getenv("EMAIL_OTP_EXPIRED"))
	if err != nil {
		logrus.Fatalf("invalid EMAIL_OTP_EXPIRED: %v", err)
	}
	phoneOTPExpired, err := time.ParseDuration(os.Getenv("SMS_OTP_EXPIRED"))
	if err != nil {
		logrus.Fatalf("invalid SMS_OTP_EXPIRED: %v", err)
	}
	OTPRegisExpired, err := time.ParseDuration(os.Getenv("SMS_OTP_REGIS_EXPIRED"))
	if err != nil {
		logrus.Fatalf("Invalid SMS_OTP_REGIS_EXPIRED: %v", err)
	}

	return &OTPConfig{
		EmailOTPExpired:       emailOTPExpired,
		PhoneOTPExpired:       phoneOTPExpired,
		TemplateEmailOTPRegis: os.Getenv("SEND_OTP_EMAIL_TEMPLATE_REGIS_ID"),
		OTPRegisExpired:       OTPRegisExpired,
		TemplateSMSOTPRegis:   os.Getenv("SEND_OTP_SMS_TEMPLATE_REGIS_ID"),
	}
}
