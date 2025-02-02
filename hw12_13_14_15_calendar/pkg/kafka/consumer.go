package kafka

//nolint:depguard
import (
	"context"

	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/configs"
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/pkg/app"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Consumer struct {
	Logger         zap.Logger
	Config         configs.KafkaConfig
	ConsumeMessage func(context.Context, app.Storage, *sarama.ConsumerMessage) error
	Storage        app.Storage
	Context        context.Context
}

func (consumer Consumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumer Consumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (consumer Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		consumer.Logger.Info("consume message:", zap.ByteString("body:", msg.Value))
		err := consumer.ConsumeMessage(consumer.Context, consumer.Storage, msg)
		if err != nil {
			consumer.Logger.Error("consume message failed", zap.Error(err))
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}

func (consumer Consumer) getConsumerGroup() (sarama.ConsumerGroup, error) {
	return sarama.NewConsumerGroup([]string{consumer.Config.URL}, consumer.Config.Group, consumer.getConsumerConfig())
}

func (consumer Consumer) RegisterMessageConsumer() error {
	consumerGroup, err := consumer.getConsumerGroup()
	if err != nil {
		return err
	}

	consumer.Logger.Info("message consumer start messaging in address: " + consumer.Config.URL)
	for {
		err := consumerGroup.Consume(consumer.Context, []string{consumer.Config.ConsumeTopic}, consumer)
		if err != nil {
			return err
		}
	}
}

func (consumer Consumer) getConsumerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.DefaultVersion // specify appropriate Kafka version
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	return config
}
