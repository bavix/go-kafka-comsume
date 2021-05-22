package configure

import (
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"os"
	"strconv"
	"strings"
)

type ConsumerConfig struct {
	GroupID string
}

type SchemaRegistry struct {
	URL string
}

type RetailRocketConfig struct {
	Topic     string
	URL       string
	ApiKey    string
	PartnerID string
}

type Config struct {
	Sarama         *sarama.Config
	Consumer       *ConsumerConfig
	RetailRocket   *RetailRocketConfig
	SchemaRegistry *SchemaRegistry
	BrokerList     []string
}

func New() *Config {
	saramaConfig := kafka.DefaultSaramaSubscriberConfig()
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	if getEnvAsBool("BROKER_SASL_ENABLE", false) {
		saramaConfig.Net.SASL.Enable = true
		saramaConfig.Net.SASL.Mechanism = sarama.SASLMechanism(getEnv("BROKER_SASL_MECHANISM", "PLAIN"))
		saramaConfig.Net.SASL.User = getEnv("BROKER_SASL_USERNAME", "username")
		saramaConfig.Net.SASL.Password = getEnv("BROKER_SASL_PASSWORD", "password")
	}

	return &Config{
		Sarama: saramaConfig,
		Consumer: &ConsumerConfig{
			GroupID: getEnv("CONSUMER_GROUP_ID", "vi-tech-group"),
		},
		RetailRocket: &RetailRocketConfig{
			Topic:     getEnv("RETAIL_ROCKET_TRACKING_TOPIC", "retail_rocket.topic"),
			URL:       getEnv("RETAIL_ROCKET_TRACKING_URL", "https://apptracking.retailrocket.net/1.0/"),
			ApiKey:    getEnv("RETAIL_ROCKET_TRACKING_API_KEY", ""),
			PartnerID: getEnv("RETAIL_ROCKET_TRACKING_PARTNER_ID", ""),
		},
		SchemaRegistry: &SchemaRegistry{
			URL: getEnv("SCHEMA_REGISTRY_URL", "http://schema-registry:8081"),
		},
		BrokerList: getEnvAsSlice("BROKER_LIST", []string{}, ","),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}
