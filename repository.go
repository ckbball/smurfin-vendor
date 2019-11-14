package main

import (
  "context"
  pb "github.com/ckbball/smurfin-vendor/proto/vendor"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

type repository interface {
  SignUp(vendor *pb.Vendor) (*pb.Vendor, error)
  Login(vendor *pb.Vendor) (*pb.Token, error)
  TakeAccountDown(request *pb.Request) (bool, error)
  ListAccounts(request *pb.Request) ([]*pb.Account, error)
  PutAccountUp(request *pb.Request) (*pb.Account, error)
  Withdraw(request *pb.Request) (string, error)
  ValidateToken(token *pb.Token) (*pb.Token, error)
}

type VendorRepository struct {
  collection *mongo.Collection
}

func (repository *VendorRepository) ListAccounts(req *pb.Request) ([]*pb.Account, error) {
  // make bson filter object
  filter := bson.D{
    {
      "VendorId",
      {
        "$eq",
        req.VendorId,
      }
    }
  }

  findOptions := options.Find()
  findOptions.SetLimit(20)
  findOptions.SetSort(SortBy.Descending("id"))
  findOptions.SetSkip(req.PageNum)

  var items []*pb.Account
  cur, err := repository.collection.Find(context.TODO(), filter, findOptions)
  if err != nil {
    return nil, err
  }
  defer cur.Close(context.TODO())

  for cur.Next(context.TODO()) {
    var elem *pb.Account
    err := cur.Decode(&elem)
    if err != nil {
      return nil, err
    }

    items = append(items, &elem)
  }

  if err := cur.Err(); err != nil {
    return items, err
  }

  return items, nil
}
