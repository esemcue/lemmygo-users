package main

import (
	"context"
	"encoding/json"
	"fmt"

	pb "gitlab.com/lemmyGo/lemmyGoUsers/proto"
	"go.mongodb.org/mongo-driver/bson"
)

func (s mUserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	fmt.Printf("Updating user: %s...", req.Email)

	// Parse the user data from JSON
	var updatedUser User
	if err := json.Unmarshal([]byte(req.UserData), &updatedUser); err != nil {
		fmt.Printf("Error parsing user data: %v", err)
		return nil, err
	}

	// Update the user document in MongoDB
	filter := bson.D{{Key: "_id", Value: req.Email}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "instances", Value: updatedUser.Instances},
	}}}

	result, err := db.Collection("Users").UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Printf("Error updating user: %v", err)
		return nil, err
	}

	fmt.Printf("User updated successfully. Matched: %d, Modified: %d", result.MatchedCount, result.ModifiedCount)

	return &pb.UpdateUserResponse{
		Message: fmt.Sprintf("User updated successfully. Modified: %d", result.ModifiedCount),
	}, nil
}
