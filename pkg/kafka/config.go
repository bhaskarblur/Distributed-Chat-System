package kafka

type Config struct {
	Brokers    string   // Comma-separated list of Kafka brokers
	GroupID    string   // Consumer group ID
	Topics     []string // Topics to consume
	AutoOffset string   // Auto offset reset (earliest/latest)
	ProducerID string   // Producer ID for logging
}

// DefaultConfig provides a default Kafka configuration
func DefaultConfig() *Config {
	return &Config{
		Brokers:    "localhost:9092",
		GroupID:    "default-group",
		Topics:     []string{"example-topic"},
		AutoOffset: "earliest",
		ProducerID: "default-producer",
	}
}
