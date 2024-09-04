package models

import pb "user-service/internal/proto/usersPB"

type User struct {
	Id        int64  `json:"id" gorm:"id" fake:"{number:1,100}"`
	Username  string `json:"username" gorm:"username" fake:"{username}"`
	FirstName string `json:"firstName" gorm:"first_name" fake:"{firstname}"`
	LastName  string `json:"lastName" gorm:"last_name" fake:"{lastname}"`
	Email     string `json:"email" gorm:"email" fake:"{email}"`
	Phone     string `json:"phone" gorm:"phone" fake:"{phone}"`
	City      string `json:"city" gorm:"city" fake:"{city}"`
	Password  string `json:"password" gorm:"password" fake:"{password:true,true,true,true,true,8}"`
	Role      string `json:"role" gorm:"role"`
	AvatarURL string `json:"avatarURL" gorm:"avatar_url"`
	CreatedAt int64  `json:"createdAt" gorm:"created_at"`
	DeletedAt int64  `json:"deletedAt" gorm:"deleted_at"`
}

type UserProfile struct {
	Id        int64  `json:"id" gorm:"id"`
	Username  string `json:"username" gorm:"username"`
	FirstName string `json:"firstName" gorm:"first_name"`
	LastName  string `json:"lastName" gorm:"last_name"`
	Email     string `json:"email" gorm:"email"`
	Phone     string `json:"phone" gorm:"phone"`
	City      string `json:"city" gorm:"city"`
	Role      string `json:"role" gorm:"role"`
	AvatarURL string `json:"avatarURL" gorm:"avatar_url"`
}

func MarshalUserPb(user *User) *pb.User {
	if user == nil {
		return nil
	}

	return &pb.User{
		Id:        user.Id,
		Username:  user.Username,
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		City:      user.City,
		Password:  user.Password,
		Role:      user.Role,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
		DeletedAt: user.DeletedAt,
	}
}

func MarshalUsersListPb(users []User) []*pb.User {
	resp := make([]*pb.User, 0, len(users))

	for _, user := range users {
		resp = append(resp, MarshalUserPb(&user))
	}

	return resp
}
