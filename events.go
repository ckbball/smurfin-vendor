// 2 event structs

// take down

package main

import (
  "context"
  "github.com/ThreeDotsLabs/watermill"
  "github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
  "github.com/ThreeDotsLabs/watermill/message"
  pb "github.com/ckbball/smurfin-vendor/proto/vendor"
  "log"
  "time"
)

// to catalog
type AccountTakenDownEvent struct {
  Id string
}

// to agent
type AccountSubmittedEvent struct {
  Account *pb.Account
}

func InitSubscriber(config kafka.SubscriberConfig) *kafka.Subscriber {
  subscriber, err := kafka.NewSubscriber(
    kafka.SubscriberConfig{
      Brokers:               []string{"kafka:9092"},
      Unmarshaler:           kafka.DefaultMarshaler{},
      OverwriteSaramaConfig: config,
      ConsumerGroup:         "test_consumer_group",
    },
    watermill.NewStdLogger(false, false),
  )
  if err != nil {
    panic(err)
  }
  return subscriber
}

func InitPublisher() *kafka.Publisher {
  publisher, err := kafka.NewPublisher(
    kafka.PublisherConfig{
      Brokers:   []string{"kafka:9092"},
      Marshaler: kafka.DefaultMarshaler{},
    },
    watermill.NewStdLogger(false, false),
  )
  if err != nil {
    panic(err)
  }
  return publisher
}
