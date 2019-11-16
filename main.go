package main

import (
  "context"
  "fmt"
  pb "github.com/ckbball/smurfin-vendor/proto/vendor"
  "github.com/micro/go-micro"
  "github.com/Shopify/sarama"
  "log"
  "os"
)

const (
  defaultHost = "datastore:27018"
)

func main() {
  srv := micro.NewService(
    micro.Name("smurfin.vendor")
  )

  srv.Init()

  uri := os.Getenv("DB_HOST")
  if uri == "" {
    uri = defaultHost
  }

  client, err := CreateClient(uri)
  if err != nil {
    log.Panic(err)
  }
  defer client.Disconnect(context.TODO())

  journalCollection := client.Database("smurfin-vendor").Collection("journal")
  jRepository := &JournalRepository{
    journalCollection,
  }

  vendorCollection := client.Database("smurfin-vendor").Collection("vendor")
  vRepo := &VendorRepository{
    vendorCollection,
  }

  // Make subscriber config here
  saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
  saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

  // Make subscriber pointer here
  subscriber := InitSubscriber(saramaSubscriberConfig)

  // Make publisher pointer here
  publisher := InitPublisher()

  // Make token service here
  tokenService := &TokenService{vRepo}

  // Make handler here with stuff
  h := &handler{vRepo, jRepository, tokenService, subscriber, publisher}

  // Register handler and server
  pb.RegisterVendorServiceHandler(srv.Server(), h)

  // Run Server
  if err := srv.Run(); err != nil {
    fmt.Println(err)
  }
}