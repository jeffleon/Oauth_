package domain

type KafkaRepository interface {
	ProduceMsg(topic, key string, message map[string]interface{}) error
}
