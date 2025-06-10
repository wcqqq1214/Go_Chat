package kafka

import (
	myconfig "go_chat/internal/config"
	"go_chat/pkg/zlog"
	"time"

	"github.com/segmentio/kafka-go"
)

type kafkaService struct {
	ChatWriter *kafka.Writer
	ChatReader *kafka.Reader
	kafkaConn  *kafka.Conn
}

var KafkaService = new(kafkaService)

// KafkaInit 初始化 kafka
func (k *kafkaService) KafkaInit() {
	// k.CreateTopic()
	kafkaConfig := myconfig.GetConfig().KafkaConfig

	k.ChatWriter = &kafka.Writer{
		Addr:                   kafka.TCP(kafkaConfig.HostPort),
		Topic:                  kafkaConfig.ChatTopic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           kafkaConfig.Timeout * time.Second,
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: false,
	}

	k.ChatReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaConfig.HostPort},
		Topic:          kafkaConfig.ChatTopic,
		CommitInterval: kafkaConfig.Timeout * time.Second,
		GroupID:        "chat",
		StartOffset:    kafka.LastOffset,
	})
}

func (k *kafkaService) KafkaClose() {
	if err := k.ChatWriter.Close(); err != nil {
		zlog.Error(err.Error())
	}

	if err := k.ChatReader.Close(); err != nil {
		zlog.Error(err.Error())
	}
}

// CreateTopic 创建 topic
func (k *kafkaService) CreateTopic() {
	// 如果已经有 topic 了，就不创建了
	kafkaConfig := myconfig.GetConfig().KafkaConfig

	chatTopic := kafkaConfig.ChatTopic

	// 连接至任意 kafka 节点
	var err error
	k.kafkaConn, err = kafka.Dial("tcp", kafkaConfig.HostPort)
	if err != nil {
		zlog.Error(err.Error())
	}

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             chatTopic,
			NumPartitions:     kafkaConfig.Partition,
			ReplicationFactor: 1,
		},
	}

	// 创建 topic
	if err = k.kafkaConn.CreateTopics(topicConfigs...); err != nil {
		zlog.Error(err.Error())
	}
}
