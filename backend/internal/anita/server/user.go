package server

import (
	"context"

	"connectrpc.com/connect"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain/value"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/service"
	v1 "github.com/ictsc/ictsc-outlands/backend/internal/proto/anita/v1"
	"github.com/ictsc/ictsc-outlands/backend/internal/proto/anita/v1/v1connect"
)

// UserServiceHandler ユーザーサービスのハンドラ
type UserServiceHandler struct {
	srv service.UserService
}

var _ v1connect.UserServiceHandler = (*UserServiceHandler)(nil)

// NewUserServiceHandler ユーザーサービスのハンドラを作成する
func NewUserServiceHandler(srv service.UserService) *UserServiceHandler {
	return &UserServiceHandler{srv: srv}
}

// GetUser ユーザーを取得する
func (s *UserServiceHandler) GetUser(ctx context.Context, req *connect.Request[v1.GetUserRequest]) (*connect.Response[v1.GetUserResponse], error) {
	userID, err := value.NewUserID(req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	user, err := s.srv.ReadUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.GetUserResponse{
		User: fromDomainUser(user),
	}), nil
}

// GetUsers ユーザー一覧を取得する
func (s *UserServiceHandler) GetUsers(ctx context.Context, _ *connect.Request[v1.GetUsersRequest]) (*connect.Response[v1.GetUsersResponse], error) {
	users, err := s.srv.ReadUsers(ctx)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.GetUsersResponse{
		Users: fromDomainUsers(users),
	}), nil
}

// PostUser ユーザーを作成する
func (s *UserServiceHandler) PostUser(ctx context.Context, req *connect.Request[v1.PostUserRequest]) (*connect.Response[v1.PostUserResponse], error) {
	id, err := value.NewUserID(req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	name, err := value.NewUserName(req.Msg.GetName())
	if err != nil {
		return nil, err
	}

	invCode, err := value.NewTeamInvitationCode(req.Msg.GetInvitationCode())
	if err != nil {
		return nil, err
	}

	user, err := s.srv.CreateUser(ctx, id, name, invCode)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.PostUserResponse{
		User: fromDomainUser(user),
	}), nil
}

// PatchUser ユーザーを更新する
func (s *UserServiceHandler) PatchUser(
	ctx context.Context, req *connect.Request[v1.PatchUserRequest],
) (*connect.Response[v1.PatchUserResponse], error) {
	id, err := value.NewUserID(req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	name, err := value.NewUserName(req.Msg.GetName())
	if err != nil {
		return nil, err
	}

	user, err := s.srv.UpdateUser(ctx, id, name)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.PatchUserResponse{
		User: fromDomainUser(user),
	}), nil
}

// DeleteUser ユーザーを削除する
func (s *UserServiceHandler) DeleteUser(
	ctx context.Context, req *connect.Request[v1.DeleteUserRequest],
) (*connect.Response[v1.DeleteUserResponse], error) {
	id, err := value.NewUserID(req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	err = s.srv.DeleteUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.DeleteUserResponse{}), nil
}
