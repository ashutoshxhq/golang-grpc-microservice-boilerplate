package user

import (
	"egnite.app/microservices/user/config"
	"egnite.app/microservices/user/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

// Server represents the gRPC server
type Server struct {
}

// GetUsers fuctions returns list of all users
func (s *Server) GetUsers(ctx context.Context, req *GetUsersRequest) (*GetUsersResponse, error) {
	client := config.Database
	var users []*User
	collection := client.Database("database").Collection("users")
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return &GetUsersResponse{Success: false, Err: "users_not_found"}, err
	}
	for cur.Next(context.TODO()) {
		var user models.User
		_ = cur.Decode(&user)
		users = append(users, &User{Id: user.ID.Hex(), Name: user.Name, Email: user.Email, Username: user.Username, Phone: user.Phone, Role: user.Role})
	}
	return &GetUsersResponse{Success: true, Users: users}, nil
}
