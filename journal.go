package main

import (
  "context"
  pb "github.com/ckbball/smurfin-vendor/proto/vendor"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "log"
)

type journalRepository interface {
  CreateJournalEntry(event interface{}) error
}

type JournalRepository struct {
  collection *mongo.Collection
}

func (repository *JournalRepository) CreateJournalEntry(event interface{}) error {
  // Somehow determine which event it is, of the two
  v, ok := event.(*AccountTakenDownEvent)
  if ok {
    work := v
    _, err = repository.collection.InsertOne(context.Background(), work)
    return nil
  }
  w, ok := event.(*AccountSubmittedEvent)
  if ok {
    work := w
    _, err = repository.collection.InsertOne(context.Background(), work)
    return nil
  }
  return errors.New("Event does not match AccountTakenDownEvent or AccountSubmittedEvent in CreatingJournalEntry")
}
