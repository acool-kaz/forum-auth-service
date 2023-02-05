package grpc

import (
	"context"
	"log"

	"github.com/acool-kaz/forum-auth-service/internal/models"
	"github.com/acool-kaz/forum-auth-service/internal/service"
	auth_svc_pb "github.com/acool-kaz/forum-auth-service/pkg/auth_svc/pb"
)

type AuthSvcHandler struct {
	auth_svc_pb.UnimplementedAuthServiceServer
	service *service.Service
}

func InitAuthSvcHandler(service *service.Service) *AuthSvcHandler {
	log.Println("init grpc auth service handler")
	return &AuthSvcHandler{
		service: service,
	}
}

func (a *AuthSvcHandler) Register(ctx context.Context, req *auth_svc_pb.RegisterRequest) (*auth_svc_pb.RegisterResponse, error) {
	user := models.User{
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Email:     req.GetEmail(),
		Username:  req.GetUsername(),
		Password:  req.GetPassword(),
	}

	id, err := a.service.Auth.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	return &auth_svc_pb.RegisterResponse{Id: int64(id)}, nil
}

func (a *AuthSvcHandler) Login(ctx context.Context, req *auth_svc_pb.LoginRequest) (*auth_svc_pb.LoginResponse, error) {
	user := models.User{
		Email:    *req.Email,
		Username: *req.Username,
		Password: req.Password,
	}

	access, refresh, err := a.service.Auth.Login(ctx, user)
	if err != nil {
		return nil, err
	}

	return &auth_svc_pb.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (a *AuthSvcHandler) Validate(ctx context.Context, req *auth_svc_pb.ValidateRequest) (*auth_svc_pb.ValidateResponse, error) {
	userId, err := a.service.Auth.ParseToken(ctx, req.GetAccessToken())
	if err != nil {
		return nil, err
	}

	return &auth_svc_pb.ValidateResponse{UserId: int32(userId)}, nil
}

func (a *AuthSvcHandler) Refresh(ctx context.Context, req *auth_svc_pb.RefreshRequest) (*auth_svc_pb.RefreshResponse, error) {
	return nil, nil
}
