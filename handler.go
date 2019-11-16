package main

import (
  "context"
  "encoding/json"
  "errors"
  "fmt"
  "github.com/ThreeDotsLabs/watermill"
  "github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
  "github.com/ThreeDotsLabs/watermill/message"
  pb "github.com/ckbball/smurfin-vendor/proto/vendor"
  "golang.org/x/crypto/bcrypt"
  "log"
  "time"
)

type handler struct {
  repo         repository
  journal      journalRepository
  tokenService Authable
  subscriber   message.Subscriber
  publisher    message.Publisher
}

func (h *handler) SignUp(ctx context.Context, req *pb.Vendor, res *pb.Response) error {
  req.Status = "preliminary"

  // Generate hash of password
  hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
  if err != nil {
    return errors.New(fmt.Sprintf("error hashing password: %v", err))
  }

  req.Password = string(hashedPass)
  if err := h.repo.Create(req); err != nil {
    return errors.New(fmt.Sprintf("error creating user : %v", err))
  }

  token, err := h.tokenService.Encode(req)
  if err != nil {
    return err
  }

  res.Vendor = req
  res.Token = &pb.Token{Token: token}

  return nil
}

func (h *handler) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {
  // Decode token
  claims, err := h.tokenService.Decode(req.Token)

  if err != nil {
    return err
  }

  if claims.Vendor.Id == "" {
    return errors.New("Invalid Vendor")
  }

  res.Valid = true

  return nil
}

func (h *handler) Login(ctx context.Context, req *pb.Vendor, res *pb.Token) error {
  vendor, err := h.repo.GetByEmail(req.Email)
  log.Println(vendor, err)
  if err != nil {
    return err
  }

  // Compare given password to stored hash
  if err := bcrypt.CompareHashAndPassword([]byte(vendor.Password), []byte(req.Password)); err != nil {
    return err
  }

  token, err := h.tokenService.Encode(vendor)
  if err != nil {
    return err
  }

  res.Token = token
  return nil
}

func (h *handler) UpdateVendor(ctx context.Context, req *pb.Vendor, res *pb.Response) error {
  if vendor, err := h.repo.Update(req); err != nil {
    return err
  }
  res.Vendor = vendor
  return nil
}

func (h *handler) GetVendor(ctx context.Context, req *pb.Vendor, res *pb.Response) error {
  if vendor, err := h.repo.Get(req.Id); err != nil {
    return err
  }
  res.Vendor = vendor
  return nil
}

// Submit an account to be added to catalog, must be validated by smurfin agents
func (h *handler) PutAccountUp(ctx context.Context, req *pb.Request, res *pb.Response) error {
  // validate vendor id
  if vendor, err := h.repo.Get(req.VendorId); err != nil {
    return err
  }
  if vendor == "" {
    return errors.New("Invalid Vendor Id")
  }

  // then add account to db with pending status
  req.Account.status = "pending"
  if account, err := h.repo.CreateAccount(req); err != nil {
    return err
  }
  res.Account = account

  // then publish event AccountSubmitted which the agent-service is listening to so it can be handed off to agents to complete validation process

  return nil
}

// Take account off of catalog?
// idk what this supposed to be
func (h *handler) TakeAccountDown(ctx context.Context, req *pb.Request, res *pb.Response) error {
  // validate vendor id
  if vendor, err := h.repo.Get(req.VendorId); err != nil {
    return err
  }
  if vendor == "" {
    return errors.New("Invalid Vendor Id")
  }

  // find account
  req.Status = "down"
  if account, err := h.repo.UpdateAccount(req); err != nil {
    return err
  }

  // set response
  res.Account = account

  // publish event AccountTakenDown, catalog listening so that it will remove this account from its db
  // type AccountTakenDown struct {
  //   Id: account.Id
  // }

  // return
  return nil
}

// List all accounts for a specific vendor,
func (h *handler) ListAccounts(ctx context.Context, req *pb.Request, res *pb.Response) error {
  // grab accounts
  if accounts, err := h.repo.ListAccounts(req); err != nil {
    return err
  }

  // set response
  res.Accounts = accounts
  //return
  return nil
}
