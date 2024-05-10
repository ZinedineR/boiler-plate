package appconf

import "os"

type CustomerServiceConfig struct {
	CustomerSvcGetDetailInvestorURL string `validate:"required" name:"CUSTOMER_SERVICE_GET_DETAIL_INVESTOR_URL"`
	CustomerSvcAPIKey               string `validate:"required" name:"CUSTOMER_SERVICE_API_KEY"`
}

func CustomerServiceConfigInit() *CustomerServiceConfig {
	return &CustomerServiceConfig{
		CustomerSvcGetDetailInvestorURL: os.Getenv("CUSTOMER_SERVICE_GET_DETAIL_INVESTOR_URL"),
		CustomerSvcAPIKey:               os.Getenv("CUSTOMER_SERVICE_API_KEY"),
	}
}
