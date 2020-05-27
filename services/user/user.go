package user

import (
	"time"

	"egnite.app/microservices/user/config"
	"egnite.app/microservices/user/helpers"
	"egnite.app/microservices/user/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

// Server represents the gRPC server
type Server struct {
}

// GetUsers fuctions returns list of all users
func (s *Server) GetUsers(ctx context.Context, req *GetUsersRequest) (*GetUsersResponse, error) {
	client := config.Database
	var users []*User
	collection := client.Database("egnite").Collection("users")
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

// GetUserByID returns single user by id
func (s *Server) GetUserByID(ctx context.Context, req *GetUserByIDRequest) (*GetUserByIDResponse, error) {
	var user models.User
	client := config.Database
	collection := client.Database("egnite").Collection("users")
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return &GetUserByIDResponse{Success: false, Err: "id_not_valid"}, err
	}
	documentReturned := collection.FindOne(context.TODO(), bson.M{"_id": id})
	err = documentReturned.Decode(&user)
	if err != nil {
		return &GetUserByIDResponse{Success: false, Err: "user_not_found"}, nil
	}

	return &GetUserByIDResponse{Success: true, User: &User{Id: user.ID.Hex(), Name: user.Name, Email: user.Email, Username: user.Username, Phone: user.Phone, Role: user.Role}}, nil
}

// GetUserByUsername returns single user by username
func (s *Server) GetUserByUsername(ctx context.Context, req *GetUserByUsernameRequest) (*GetUserByUsernameResponse, error) {
	var user *models.User
	client := config.Database
	collection := client.Database("egnite").Collection("users")
	documentReturned := collection.FindOne(context.TODO(), bson.M{"username": req.Username})
	err := documentReturned.Decode(&user)
	if err != nil {
		return &GetUserByUsernameResponse{Success: false, Err: "user_not_found"}, nil
	}
	return &GetUserByUsernameResponse{Success: true, User: &User{Id: user.ID.Hex(), Name: user.Name, Email: user.Email, Username: user.Username, Phone: user.Phone, Role: user.Role}}, nil
}

// GetUserByUsernameWithPassword returns single user by username
func (s *Server) GetUserByUsernameWithPassword(ctx context.Context, req *GetUserByUsernameWithPasswordRequest) (*GetUserByUsernameWithPasswordResponse, error) {
	var user *models.User
	client := config.Database
	collection := client.Database("egnite").Collection("users")
	documentReturned := collection.FindOne(context.TODO(), bson.M{"username": req.Username})
	err := documentReturned.Decode(&user)
	if err != nil {
		return &GetUserByUsernameWithPasswordResponse{Success: false, Err: "user_not_found"}, nil
	}
	return &GetUserByUsernameWithPasswordResponse{Success: true, User: &User{Id: user.ID.Hex(), Name: user.Name, Email: user.Email, Username: user.Username, Phone: user.Phone, Role: user.Role}, Password: user.Password}, nil
}

// CreateUser stores new user in database and returns id
func (s *Server) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	var user models.User
	user.Name = req.User.Name
	user.Username = req.User.Username
	user.Email = req.User.Email
	user.Phone = req.User.Phone
	user.Role = req.User.Role
	user.Password, _ = helpers.HashPassword(req.Password)
	user.CreatedAt = time.Now()

	client := config.Database
	collection := client.Database("egnite").Collection("users")
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return &CreateUserResponse{Success: false, Err: "unable_to_inser_user"}, nil
	}
	return &CreateUserResponse{Success: true, Id: result.InsertedID.(primitive.ObjectID).Hex()}, nil
}

// UpdateUserByUsername updates user by username anad returns success state
func (s *Server) UpdateUserByUsername(ctx context.Context, req *UpdateUserByUsernameRequest) (*UpdateUserByUsernameResponse, error) {
	update := make(map[string]string)
	if req.User.Name != "" {
		update["name"] = req.User.Name
	}
	if req.User.Email != "" {
		update["email"] = req.User.Email
	}
	if req.User.Phone != "" {
		update["phone"] = req.User.Phone
	}
	if req.User.Username != "" {
		update["username"] = req.User.Username
	}

	client := config.Database
	collection := client.Database("egnite").Collection("users")
	_, err := collection.UpdateOne(context.TODO(), bson.M{"username": req.User.Username}, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return &UpdateUserByUsernameResponse{Success: false, Err: "unable_to_update_user"}, nil
	}
	return &UpdateUserByUsernameResponse{Success: true}, nil
}

// DeleteUserByID deletes a user by id
func (s *Server) DeleteUserByID(ctx context.Context, req *DeleteUserByIDRequest) (*DeleteUserByIDResponse, error) {
	client := config.Database
	collection := client.Database("egnite").Collection("users")
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return &DeleteUserByIDResponse{Success: false, Err: "id_not_valid"}, nil
	}
	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return &DeleteUserByIDResponse{Success: false, Err: "unable_to_delete_user"}, nil
	}
	return &DeleteUserByIDResponse{Success: true}, nil
}
