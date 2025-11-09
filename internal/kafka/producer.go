package kafka

import (
    "encoding/json"
    "log"
    "github.com/IBM/sarama"
    "github.com/AkshatPandey-2004/4-in-a-row/pkg/models"
)

type Producer struct {
    producer sarama.SyncProducer
    topic    string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Retry.Max = 5

    producer, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        return nil, err
    }

    log.Println("Kafka producer connected")
    return &Producer{
        producer: producer,
        topic:    topic,
    }, nil
}

func (p *Producer) SendGameEvent(event *models.GameEvent) error {
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }

    msg := &sarama.ProducerMessage{
        Topic: p.topic,
        Value: sarama.ByteEncoder(data),
    }

    _, _, err = p.producer.SendMessage(msg)
    if err != nil {
        log.Printf("Failed to send message: %v", err)
        return err
    }

    return nil
}

func (p *Producer) Close() error {
    return p.producer.Close()
}