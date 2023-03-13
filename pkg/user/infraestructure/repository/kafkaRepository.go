package repository

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jeffleon/oauth-microservice/pkg/user/domain"
	"github.com/sirupsen/logrus"
)

type kafkaRepository struct {
	producer *kafka.Producer
}

func NewKafkaRepository(producer *kafka.Producer) domain.KafkaRepository {
	return &kafkaRepository{
		producer,
	}
}

func (k *kafkaRepository) ProduceMsg(topic, key string, message map[string]interface{}) error {

	go func(topic string, message map[string]interface{}) {
		byteMsg, err := json.Marshal(message)
		if err != nil {
			logrus.Errorf("Error marshall error sending to topic %s msg %v", topic, message)
		}

		err = k.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte(key),
			Value:          byteMsg,
		}, nil)

		if err != nil {
			logrus.Errorf("Error producing msg %v to topic %s", message, topic)
		}

		k.producer.Flush(15000)
	}(topic, message)

	return nil

}
