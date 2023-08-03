package main

import (
	"context"
	"fmt"

	pb "gitlab.com/lemmyGo/lemmyGoUsers/proto"
	"go.mongodb.org/mongo-driver/bson"
)

type mUserServer struct {
	pb.UnimplementedUsersServer
}

func (s mUserServer) Register(ctx context.Context, req *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {
	db := GetMongoDb()
	db.Collection("Users").InsertOne(ctx, bson.D{{Key: "username", Value: req.Name}, {Key: "password", Value: req.Password}, {Key: "email", Value: req.Email}})
	return &pb.RegistrationResponse{
		Message: fmt.Sprintf("Info - %s:%s", req.Name, req.Password),
	}, nil
}

func (s mUserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	db := GetMongoDb()
	db.Collection("Users").FindOne(ctx, bson.D{{Key: "username", Value: req.Name}})
	return &pb.LoginResponse{
		Message: fmt.Sprintf("Info - %s:%s", req.Name, req.Password),
	}, nil
}
