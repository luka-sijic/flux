package handler

import "github.com/luka-sijic/flux/internal/service"

type UserHandler struct {
	svc service.Service
}

type FriendHandler struct {
	svc service.Service
}

func NewUserHandler(s service.Service) *UserHandler {
	return &UserHandler{svc: s}
}

func NewFriendHandler(s service.Service) *FriendHandler {
	return &FriendHandler{svc: s}
}
