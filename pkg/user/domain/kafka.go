package domain

type KafkaRepository interface {
	ProduceMsg(topic string, message map[string]interface{}) error
}
