package main

import (
  "context"
  pb "github.com/ckbball/smurfin-vendor/proto/vendor"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

type repository interface {
  UpdateAccount(request *pb.Request) (bool, error)
  ListAccounts(request *pb.Request) ([]*pb.Account, error) //
  CreateAccount(request *pb.Request) (*pb.Account, error)
  Create(vendor *pb.Vendor) error//
  Update(vendor *pb.Vendor) (*pb.Vendor, error)//
  Delete(vendor *pb.Vendor) error//
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

func (repo *VendorRepository) UpdateAccount(request *pb.Request) (bool, error) {
  // get account from db

  // delete account?
  deleteResult, err := repo.collection.Delete(context.TODO(), bson.D{{"Id", { "$eq", request.Account.Id}}})
  if err != nil {
    return err
  }

  _, err = repo.collection.InsertOne(context.Background(), request.Account)
  return true, err
}

func (repo *VendorRepository) CreateAccount(request *pb.Request) (*pb.Account, error) {
  _, err = repo.collection.InsertOne(context.Background(), request.Account)
  return request.Account, err
}

// Do password crypto stuff in handler
func (repo *VendorRepository) Create(vendor *pb.Vendor) (*pb.Vendor, error) {
  // check that email is not already in use
  rep, err := repo.collection.FindOne(context.TODO(), bson.D{{"VendorId", { "$eq", vendor.Id}}}, nil)
  if err != nil {
    return nil, error
  }
  if rep != nil {
    return nil, nil
  }

  // create new account
  _, err := repo.collection.InsertOne(context.Background(), vendor)

  // return vendor with updated model
  return vendor, err
}

func (repo *VendorRepository) Update(vendor *pb.Vendor) (*pb.Vendor, error) {
  deleteResult, err := repo.collection.Delete(context.TODO(), bson.D{{"VendorId", { "$eq", vendor.Id}}})
  if err != nil {
    return err
  }

  _, err = repo.collection.InsertOne(context.Background(), vendor)
  return vendor, err
}

func (repo *VendorRepository) Delete(vendor *pb.Vendor) error {
  deleteResult, err := repo.collection.Delete(context.TODO(), bson.D{{"VendorId", { "$eq", vendor.Id}}})
  if err != nil {
    return err
  }
  fmt.Printf("Deleted %v documents\n", deleteResult.DeletedCount)
  return nil
}

func (repo *VendorRepository) UpdateAccount(request *pb.Request) (bool, error) {
  
}