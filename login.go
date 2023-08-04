package main

import (
	"context"
	"fmt"

	pb "gitlab.com/lemmyGo/lemmyGoUsers/proto"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Name     string `bson:"username"`
	Password string `bson:"password"`
	Email    string `bson:"email,omitempty"`
}

type mUserServer struct {
	pb.UnimplementedUsersServer
}

func (s mUserServer) Register(ctx context.Context, req *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {
	newUser := User{
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
	}
	res, err := db.Collection("Users").InsertOne(ctx, newUser)
	if err != nil {
		panic(err)
	}

	return &pb.RegistrationResponse{
		Message: fmt.Sprintf("Info - %s:%s", res.InsertedID, newUser.Name),
	}, nil
}

func (s mUserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	db.Collection("Users").FindOne(ctx, bson.D{{Key: "username", Value: req.Name}})
	return &pb.LoginResponse{
		Message: fmt.Sprintf("Info - %s:%s", req.Name, req.Password),
	}, nil
}
