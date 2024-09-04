package controller

import (
	"context"
	"fmt"
	"user-service/internal/models"
	pb "user-service/internal/proto/usersPB"
	"user-service/internal/service"

	"go.uber.org/zap"
)

type GRPCInterface interface {
	SearchUser(ctx context.Context, req *pb.SearchRequest) (*pb.ListResponse, error)
	RegisterUser(ctx context.Context, user *pb.User) (*pb.RegisterResponse, error)
	LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.User, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.MsgResponse, error)
	DeleteUser(ctx context.Context, req *pb.GetByIdRequest) (*pb.MsgResponse, error)
	ListUsers(context.Context, *pb.ListRequest) (*pb.ListResponse, error)
	GetProfileById(ctx context.Context, req *pb.GetByIdRequest) (*pb.Profile, error)
	UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.MsgResponse, error)
	UploadAvatar(ctx context.Context, req *pb.UploadAvatarRequest) (*pb.MsgResponse, error)
	pb.UsersServiceServer
}

type GRPCObj struct {
	service service.UserServiceInterface
	log     *zap.Logger
	pb.UsersServiceServer
}

func NewGRPCObj(userService service.UserServiceInterface, log *zap.Logger) *GRPCObj {
	return &GRPCObj{
		service: userService,
		log:     log,
	}
}

func (g *GRPCObj) RegisterUser(ctx context.Context, user *pb.User) (*pb.RegisterResponse, error) {
	u := models.User{
		Username:  user.Username,
		FirstName: user.Firstname,
		LastName:  user.Lastname,
		Email:     user.Email,
		Phone:     user.Phone,
		City:      user.City,
		Password:  user.Password,
		Role:      user.Role,
	}
	id, err := g.service.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	resp := &pb.RegisterResponse{Id: id}

	return resp, nil
}

func (g *GRPCObj) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, exp, err := g.service.AuthenticateUser(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	resp := &pb.LoginResponse{
		Token:     token,
		ExpiredAt: exp,
	}
	return resp, nil
}

func (g *GRPCObj) GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.User, error) {
	user, err := g.service.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return models.MarshalUserPb(&user), nil
}

func (g *GRPCObj) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.MsgResponse, error) {
	if req.UpdatedUser == nil {
		return nil, fmt.Errorf("empty request")
	}
	u := models.User{
		Id:        req.UpdatedUser.Id,
		Username:  req.UpdatedUser.Username,
		FirstName: req.UpdatedUser.Firstname,
		LastName:  req.UpdatedUser.Lastname,
		Email:     req.UpdatedUser.Email,
		Phone:     req.UpdatedUser.Phone,
		City:      req.UpdatedUser.City,
		Password:  req.UpdatedUser.Password,
		Role:      req.UpdatedUser.Role,
		CreatedAt: req.UpdatedUser.CreatedAt,
		DeletedAt: req.UpdatedUser.DeletedAt,
	}
	err := g.service.UpdateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	resp := &pb.MsgResponse{Msg: fmt.Sprintf("user Id = %v updated", req.UpdatedUser.Id)}

	return resp, nil
}

// Устанавливает DeletedAt, не удалая саму запись
func (g *GRPCObj) DeleteUser(ctx context.Context, req *pb.GetByIdRequest) (*pb.MsgResponse, error) {
	err := g.service.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	resp := &pb.MsgResponse{Msg: fmt.Sprintf("user Id = %v delete", req.Id)}

	return resp, nil
}

func (g *GRPCObj) ListUsers(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	list, err := g.service.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ListResponse{Users: models.MarshalUsersListPb(list)}, nil
}

func (g *GRPCObj) GetProfileById(ctx context.Context, req *pb.GetByIdRequest) (*pb.Profile, error) {
	user, err := g.service.GetUserProfile(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	resp := &pb.Profile{
		Id:        user.Id,
		Username:  user.Username,
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		City:      user.City,
		Role:      user.Role,
		AvatarURL: user.AvatarURL,
	}
	return resp, nil
}

func (g *GRPCObj) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.MsgResponse, error) {
	u := models.User{
		Id:        req.UpdatedUser.Id,
		Username:  req.UpdatedUser.Username,
		FirstName: req.UpdatedUser.Firstname,
		LastName:  req.UpdatedUser.Lastname,
		Email:     req.UpdatedUser.Email,
		Phone:     req.UpdatedUser.Phone,
		City:      req.UpdatedUser.City,
		Role:      req.UpdatedUser.Role,
		AvatarURL: req.UpdatedUser.AvatarURL,
	}
	err := g.service.UpdateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	resp := &pb.MsgResponse{Msg: fmt.Sprintf("user profile Id = %v updated", req.UpdatedUser.Id)}

	return resp, nil
}

func (g *GRPCObj) UploadAvatar(ctx context.Context, req *pb.UploadAvatarRequest) (*pb.MsgResponse, error) {

	err := g.service.UploadAvatar(ctx, req.Img, req.Filename, req.Id)
	if err != nil {
		return nil, err
	}
	resp := &pb.MsgResponse{Msg: fmt.Sprintf("upload avatar for user Id = %v", req.Id)}

	return resp, nil
}

func (g *GRPCObj) SearchUser(ctx context.Context, req *pb.SearchRequest) (*pb.ListResponse, error) {
	list, err := g.service.SearchUser(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	return &pb.ListResponse{Users: models.MarshalUsersListPb(list)}, nil
}
