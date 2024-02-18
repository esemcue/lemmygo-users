package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	pb "gitlab.com/lemmyGo/lemmyGoUsers/proto"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type lemmyInstance struct {
	Url      string `bson:"url"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}

type User struct {
	Email          string          `bson:"_id"`
	Password       string          `bson:"password"`
	LemmyInstances []lemmyInstance `bson:"lemmyInstances"`
}

type mUserServer struct {
	pb.UnimplementedUsersServer
}

// TODO split
func (s mUserServer) Register(ctx context.Context, req *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 4)
	if err != nil {
		return nil, err
	}

	newUser := User{
		Password:       string(hashed),
		Email:          req.Email,
		LemmyInstances: []lemmyInstance{},
	}
	res, err := db.Collection("Users").InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &pb.RegistrationResponse{
		Message: fmt.Sprintf("Info - %s:%s", res.InsertedID),
	}, nil
}

func (s mUserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	foundResult := db.Collection("Users").FindOne(ctx, bson.D{{Key: "_id", Value: req.Email}})
	if foundResult.Err() != nil {
		fmt.Printf("User not found: %s\n", req.Email)
		return nil, foundResult.Err()
	}

	var foundUser User
	foundResult.Decode(&foundUser)

	loginErr := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password))
	if loginErr != nil {
		fmt.Printf("Password mismatch for user %s\n", req.Email)
		return nil, loginErr
	}

	spew.Dump(foundUser)
	bytes, jsonError := json.Marshal(foundUser)
	if jsonError != nil {
		return nil, jsonError
	}

	return &pb.LoginResponse{
		Message: string(bytes),
	}, nil
}
