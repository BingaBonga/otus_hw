package kafka

//nolint:depguard
import (
	"encoding/json"

	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/configs"
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Producer struct {
	producer sarama.SyncProducer
	logger   *zap.Logger
	config   configs.KafkaConfig
}

func NewProducer(config configs.KafkaConfig, logger *zap.Logger) (*Producer, error) {
	producer, err := sarama.NewSyncProducer([]string{config.URL}, getProducerConfig(config))
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer, config: config, logger: logger}, nil
}

func (producer Producer) SendEventMessage(event *storage.Event) error {
	producer.logger.Info("send event to kafka", zap.Any("event", &event))

	jsonData, _ := json.Marshal(event)
	msg := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(event.ID),
		Topic: producer.config.ProduceTopic,
		Value: sarama.ByteEncoder(jsonData),
	}

	_, _, err := producer.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

func getProducerConfig(kafkaConfig configs.KafkaConfig) *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Idempotent = true
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 5                    // Retry up to 5 times to produce the message
	config.Producer.Return.Successes = true
	config.ClientID = kafkaConfig.ServiceName
	config.Net.MaxOpenRequests = 1
	return config
}
