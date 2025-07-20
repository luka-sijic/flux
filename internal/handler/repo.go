package handler

import "github.com/luka-sijic/flux/internal/service"

type UserHandler struct {
	svc service.UserService
}

type FriendHandler struct {
	svc service.FriendService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{svc: s}
}

func NewFriendHandler(s service.FriendService) *FriendHandler {
	return &FriendHandler{svc: s}
}
