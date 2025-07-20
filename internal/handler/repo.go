package handler

import "github.com/luka-sijic/flux/internal/service"

type UserHandler struct {
	svc service.Service
}

type FriendHandler struct {
}

func NewHandler(s service.Service) *UserHandler {
	return &UserHandler{svc: s}
}
