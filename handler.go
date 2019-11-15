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
  "time"
)

type handler struct {
  repo         repository
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
