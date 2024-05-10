package appconf

import (
	"boiler-plate/pkg/xvalidator"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	AppEnvConfig              *AppConfig
	DatabaseConfig            *DatabaseConfig
	KafkaConfig               *KafkaConfig
	AuthConfig                *AuthConfig
	NotificationServiceConfig *NotificationServiceConfig
	OTPConfig                 *OTPConfig
	CustomerServiceConfig     *CustomerServiceConfig
}

func (c Config) IsStaging() bool {
	return c.AppEnvConfig.AppEnv != "production"
}

func (c Config) IsProd() bool {
	return c.AppEnvConfig.AppEnv == "production"
}

func (c Config) IsDebug() bool {
	return c.AppEnvConfig.AppDebug == "True"
}

func InitAppConfig(validate *xvalidator.Validator) *Config {

	c := Config{
		AppEnvConfig:              AppConfigInit(),
		DatabaseConfig:            DatabaseConfigInit(),
		KafkaConfig:               KafkaConfigInit(),
		AuthConfig:                AuthConfigInit(),
		NotificationServiceConfig: NotificationServiceConfigInit(),
		OTPConfig:                 OTPConfigInit(),
		CustomerServiceConfig:     CustomerServiceConfigInit(),
	}

	// NOTIFICATION SERVICE CONFIG
	//c.SendEmailOTPURL = os.Getenv("NOTIFICATION_SERVICE_SEND_EMAIL_URL")
	//c.SendSMSOTPURL = os.Getenv("NOTIFICATION_SERVICE_SEND_SMS_URL")
	//c.TemplateEmailOTP = os.Getenv("OTP_EMAIL_TEMPLATE_ID")
	//c.TemplateSMSOTP = os.Getenv("OTP_SMS_TEMPLATE_ID")
	//c.NotificationServiceAPIKey = os.Getenv("NOTIFICATION_SERVICE_API_KEY")

	// OTP CONFIG
	//emailOTPExpired, err := time.ParseDuration(os.Getenv("EMAIL_OTP_EXPIRED"))
	//if err != nil {
	//	logrus.Fatalf("invalid EMAIL_OTP_EXPIRED: %v", err)
	//}
	//c.EmailOTPExpired = emailOTPExpired
	//phoneOTPExpired, err := time.ParseDuration(os.Getenv("SMS_OTP_EXPIRED"))
	//if err != nil {
	//	logrus.Fatalf("invalid SMS_OTP_EXPIRED: %v", err)
	//}
	//c.PhoneOTPExpired = phoneOTPExpired

	// KAFKA CONFIG
	//c.KafkaSecurityProtocol = os.Getenv("KAFKA_SECURITY_PROTOCOL")
	//c.Username = os.Getenv("KAFKA_USERNAME")
	//c.Password = os.Getenv("KAFKA_PASSWORD")
	//brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	//if len(brokers) == 0 {
	//	brokers = nil
	//}
	//c.Brokers = brokers
	//c.GroupID = os.Getenv("KAFKA_GROUP_ID")
	//c.UpdateCustomerTopic = os.Getenv("KAFKA_UPDATE_CUSTOMER_TOPIC")

	// CUSTOMER SERVICE CONFIG
	//c.CustomerSvcGetDetailInvestorURL = os.Getenv("CUSTOMER_SERVICE_GET_DETAIL_INVESTOR_URL")
	//c.CustomerSvcAPIKey = os.Getenv("CUSTOMER_SERVICE_API_KEY")

	//c.TemplateSMSOTPRegis = os.Getenv("SEND_OTP_SMS_TEMPLATE_REGIS_ID")
	//OTPRegisExpired, err := time.ParseDuration(os.Getenv("SMS_OTP_REGIS_EXPIRED"))
	//if err != nil {
	//	logrus.Fatalf("Invalid SMS_OTP_REGIS_EXPIRED: %v", err)
	//}
	//logrus.Info(OTPRegisExpired)
	//c.OTPRegisExpired = OTPRegisExpired
	//c.TemplateEmailOTPRegis = os.Getenv("SEND_OTP_EMAIL_TEMPLATE_REGIS_ID")

	// JWT ACCESS TOKEN
	//c.JWTSecretAccessToken = []byte(os.Getenv("JWT_SECRET_ACCESS_TOKEN"))
	errs := validate.Struct(c)
	if errs != nil {
		for _, e := range errs {
			logrus.Error(fmt.Sprintf("Failed to load env: %s", e))
		}
		os.Exit(1)
	}
	return &c

}
