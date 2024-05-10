package appconf

import (
	"os"
	"strings"
)

type KafkaConfig struct {
	KafkaSecurityProtocol    string   `validate:"required,eq=SASL_SSL|eq=SASL_PLAIN|eq=PLAIN" name:"KAFKA_SECURITY_PROTOCOL"`
	KafkaUsername            string   `validate:"required" name:"KAFKA_USERNAME"`
	KafkaPassword            string   `validate:"required" name:"KAFKA_PASSWORD"`
	KafkaBroker              []string `validate:"required" name:"KAFKA_BROKER"`
	KafkaGroupId             string   `validate:"required" name:"KAFKA_GROUP_ID"`
	KafkaUpdateCustomerTopic string   `validate:"required" name:"KAFKA_UPDATE_CUSTOMER_TOPIC"`
}

func KafkaConfigInit() *KafkaConfig {
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if len(brokers) == 0 {
		brokers = nil
	}
	return &KafkaConfig{
		KafkaSecurityProtocol:    os.Getenv("KAFKA_SECURITY_PROTOCOL"),
		KafkaUsername:            os.Getenv("KAFKA_USERNAME"),
		KafkaPassword:            os.Getenv("KAFKA_PASSWORD"),
		KafkaBroker:              brokers,
		KafkaGroupId:             os.Getenv("KAFKA_GROUP_ID"),
		KafkaUpdateCustomerTopic: os.Getenv("KAFKA_UPDATE_CUSTOMER_TOPIC"),
	}
}
