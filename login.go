package main

import (
	"context"
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
	Name           string          `bson:"username"`
	Password       string          `bson:"password"`
	Email          string          `bson:"_id"`
	LemmyInstances []lemmyInstance `bson:"lemmyInstances"`
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
		Name:           req.Name,
		Password:       string(hashed),
		Email:          req.Email,
		LemmyInstances: []lemmyInstance{},
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

	var foundUser User
	foundResult.Decode(&foundUser)

	loginErr := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password))
	if loginErr != nil {
		return nil, loginErr
	}

	spew.Dump(foundUser)
	return &pb.LoginResponse{
		Message: foundUser.Email,
	}, nil
}
