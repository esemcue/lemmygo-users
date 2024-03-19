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

type credentials map[string]string

type User struct {
	Email     string                 `bson:"_id"`
	Password  string                 `bson:"password"`
	Instances map[string]credentials `bson:"instances"`
}

type mUserServer struct {
	pb.UnimplementedUsersServer
}

// TODO split
func (s mUserServer) Register(ctx context.Context, req *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	newUser := User{
		Email:    req.Email,
		Password: string(hashedPass),
	}
	res, err := db.Collection("Users").InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &pb.RegistrationResponse{
		Message: fmt.Sprintf("Info - %s", res.InsertedID),
	}, nil
}

func (s mUserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	fmt.Printf("Searching for %s...", req.Email)

	foundResult := db.Collection("Users").FindOne(ctx, bson.D{{Key: "_id", Value: req.Email}})
	if foundResult.Err() != nil {
		fmt.Print("User not found! ")
		return nil, foundResult.Err()
	}

	var foundUser User
	foundResult.Decode(&foundUser)
	fmt.Print("User found. Checking password...")

	loginErr := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password))
	if loginErr != nil {
		fmt.Printf("mismatch.")
		return nil, loginErr
	}

	fmt.Println("matched. Returning.")
	spew.Dump(foundUser)

	bytes, jsonError := json.Marshal(foundUser)
	if jsonError != nil {
		return nil, jsonError
	}
	return &pb.LoginResponse{
		Message: string(bytes),
	}, nil
}
