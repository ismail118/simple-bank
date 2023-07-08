package api

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ismail118/simple-bank/models"
	pb "github.com/ismail118/simple-bank/proto"
	"github.com/ismail118/simple-bank/repository"
	"github.com/ismail118/simple-bank/token"
	"github.com/ismail118/simple-bank/util"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GrpcServer serves gRPC request for our banking service
type GrpcServer struct {
	pb.UnimplementedSimpleBankServer
	store      repository.Store
	repo       repository.Repository
	tokenMaker token.Maker
	config     *util.Config
}

// NewGrpcServer create a new gRPC server and setup routing
func NewGrpcServer(
	store repository.Store,
	repo repository.Repository,
	tokenMaker token.Maker,
	config *util.Config,
) *GrpcServer {
	server := &GrpcServer{
		store:      store,
		repo:       repo,
		tokenMaker: tokenMaker,
		config:     config,
	}

	return server
}

func (s *GrpcServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s:%s", violations[0].Field, violations[0].Description)
	}

	hashedPassword, err := util.HashedPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password err:%s", err)
	}
	user := models.Users{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	u, err := s.repo.GetUsersByUsername(ctx, user.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}
	if u.Username != "" {
		return nil, status.Errorf(codes.AlreadyExists, "user with username %s is exists", u.Username)
	}

	u, err = s.repo.GetUsersByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}
	if u.Username != "" {
		return nil, status.Errorf(codes.AlreadyExists, "email %s already being user", req.Email)
	}

	err = s.repo.InsertUsers(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	u, err = s.repo.GetUsersByUsername(ctx, user.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	return &pb.CreateUserResponse{
		User: util.ConvertUser(user),
	}, nil
}

func (s *GrpcServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s:%s", violations[0].Field, violations[0].Description)
	}

	user, err := s.repo.GetUsersByUsername(ctx, req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "err:%s", err)
	}
	if user.Username == "" {
		return nil, status.Error(codes.NotFound, "username not found")
	}

	if user.Email != req.GetEmail() {
		u, err := s.repo.GetUsersByEmail(ctx, req.GetEmail())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "err:%s", err)
		}
		if u.Username != "" {
			return nil, status.Errorf(codes.InvalidArgument, "email %s already being used", req.GetEmail())
		}
	}

	// updated
	user.FullName = req.GetFullName()
	user.Email = req.GetEmail()

	err = s.repo.UpdateUsers(ctx, repository.UpdateUserParam{
		Username: user.Username,
		FullName: sql.NullString{String: user.FullName, Valid: true},
		Email:    sql.NullString{String: user.Email, Valid: true},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "err:%s", err)
	}

	return &pb.UpdateUserResponse{
		User: util.ConvertUser(user),
	}, nil
}

func (s *GrpcServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	violations := validateLoginRequest(req)
	if violations != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s:%s", violations[0].Field, violations[0].Description)
	}

	user, err := s.repo.GetUsersByUsername(ctx, req.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}
	if user.Username == "" {
		err = errors.New("username not found")
		return nil, status.Errorf(codes.NotFound, "%s", err)
	}

	err = util.ComparePassword(user.HashedPassword, req.Password)
	if err != nil {
		err = errors.New("wrong password")
		return nil, status.Errorf(codes.NotFound, "%s", err)
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	mtdt := s.extractMetadata(ctx)

	session := models.Sessions{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiredAt:    refreshPayload.ExpiredAt,
	}

	err = s.repo.InsertSessions(ctx, session)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	return &pb.LoginResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiredAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiredAt: timestamppb.New(accessPayload.ExpiredAt),
		User:                  util.ConvertUser(user),
	}, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) []*errdetails.BadRequest_FieldViolation {
	var violation []*errdetails.BadRequest_FieldViolation
	err := util.ValidateUsername(req.GetUsername())
	if err != nil {
		violation = append(violation, util.FieldViolation("username", err))
	}

	err = util.ValidatePassword(req.GetPassword())
	if err != nil {
		violation = append(violation, util.FieldViolation("password", err))
	}

	err = util.ValidateEmail(req.GetEmail())
	if err != nil {
		violation = append(violation, util.FieldViolation("email", err))
	}

	err = util.ValidateFullName(req.GetFullName())
	if err != nil {
		violation = append(violation, util.FieldViolation("full_name", err))
	}

	return violation
}

func validateLoginRequest(req *pb.LoginRequest) []*errdetails.BadRequest_FieldViolation {
	var violation []*errdetails.BadRequest_FieldViolation
	err := util.ValidateUsername(req.GetUsername())
	if err != nil {
		violation = append(violation, util.FieldViolation("username", err))
	}

	err = util.ValidatePassword(req.GetPassword())
	if err != nil {
		violation = append(violation, util.FieldViolation("password", err))
	}

	return violation
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) []*errdetails.BadRequest_FieldViolation {
	var violation []*errdetails.BadRequest_FieldViolation
	err := util.ValidateUsername(req.GetUsername())
	if err != nil {
		violation = append(violation, util.FieldViolation("username", err))
	}

	//err = util.ValidatePassword(req.GetPassword())
	//if err != nil {
	//	violation = append(violation, util.FieldViolation("password", err))
	//}

	err = util.ValidateEmail(req.GetEmail())
	if err != nil {
		violation = append(violation, util.FieldViolation("email", err))
	}

	err = util.ValidateFullName(req.GetFullName())
	if err != nil {
		violation = append(violation, util.FieldViolation("full_name", err))
	}

	return violation
}
