package main

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	pb "gitlab.com/lemmyGo/lemmyGoUsers/proto"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string `bson:"username"`
	Password string `bson:"password"`
	Email    string `bson:"_id"`
}

type mUserServer struct {
	pb.UnimplementedUsersServer
}

func (s mUserServer) Register(ctx context.Context, req *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 4)
	if err != nil {
		return nil, err
	}

	newUser := User{
		Name:     req.Name,
		Password: string(hashed),
		Email:    req.Email,
	}
	res, err := db.Collection("Users").InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &pb.RegistrationResponse{
		Message: fmt.Sprintf("Info - %s:%s", res.InsertedID, newUser.Name),
	}, nil
}

func (s mUserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	foundResult := db.Collection("Users").FindOne(ctx, bson.D{{Key: "username", Value: req.Name}})
	if foundResult.Err() != nil {
		return nil, foundResult.Err()
	}

	spew.Dump(foundResult)
	return &pb.LoginResponse{
		Message: fmt.Sprintf("Info - %s:%s", req.Name, req.Password),
	}, nil
}
