package kafka

import (
    "context"
    "encoding/json"
    "log"
    "github.com/IBM/sarama"
    "github.com/AkshatPandey-2004/4-in-a-row/pkg/models"
)

type Consumer struct {
    consumer sarama.ConsumerGroup
    topic    string
}

func NewConsumer(brokers []string, groupID, topic string) (*Consumer, error) {
    config := sarama.NewConfig()
    config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
    config.Consumer.Offsets.Initial = sarama.OffsetNewest

    consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
    if err != nil {
        return nil, err
    }

    log.Println("Kafka consumer connected")
    return &Consumer{
        consumer: consumer,
        topic:    topic,
    }, nil
}

func (c *Consumer) Start(ctx context.Context) {
    handler := &consumerHandler{}
    
    go func() {
        for {
            if err := c.consumer.Consume(ctx, []string{c.topic}, handler); err != nil {
                log.Printf("Error from consumer: %v", err)
            }
            
            if ctx.Err() != nil {
                return
            }
        }
    }()
}

func (c *Consumer) Close() error {
    return c.consumer.Close()
}

type consumerHandler struct{}

func (h *consumerHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for message := range claim.Messages() {
        var event models.GameEvent
        if err := json.Unmarshal(message.Value, &event); err != nil {
            log.Printf("Error unmarshaling event: %v", err)
            continue
        }

        // Log analytics event
        log.Printf("Analytics Event - Type: %s, GameID: %s, Time: %v", 
            event.Type, event.GameID, event.Timestamp)
        
        // Here you can store to database or process further
        
        session.MarkMessage(message, "")
    }
    return nil
}