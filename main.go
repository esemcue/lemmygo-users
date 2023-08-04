package main

import (
	"context"
	"fmt"
	"net"
	"os"

	pb "gitlab.com/lemmyGo/lemmyGoUsers/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DbURI string `yaml:"dbUri"`
}

var db *mongo.Database

func main() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	var config Config
	yaml.Unmarshal(yamlFile, &config)
	opts := options.Client().ApplyURI(config.DbURI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	db = client.Database("Likky")

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}
	serverRegistrar := grpc.NewServer()
	service := &mUserServer{}
	pb.RegisterUsersServer(serverRegistrar, service)
	sErr := serverRegistrar.Serve(lis)
	if err != nil {
		panic(sErr)
	}

}
