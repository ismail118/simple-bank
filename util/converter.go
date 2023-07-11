package util

import (
	"github.com/ismail118/simple-bank/models"
	pb "github.com/ismail118/simple-bank/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertUser(user models.Users) *pb.User {
	return &pb.User{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
