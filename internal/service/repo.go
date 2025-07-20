package service

import (
	"github.com/luka-sijic/flux/internal/database"
	"github.com/luka-sijic/flux/internal/models"

	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type UserService interface {
	CreateUser(user *models.UserDTO) bool
	LoginUser(user *models.UserDTO) bool
}

type FriendService interface {
	AddFriend(username string, user *models.FriendDTO) bool
	GetFriends(username string) []models.FriendDTO
	GetRequests(username string) []models.FriendDTO
	FriendResponse(username string, action *models.FriendActionDTO) bool
	GetLog(username1, username2 string) []models.Messages
}

type Infra struct {
	Pools []*pgxpool.Pool
	Node  *snowflake.Node
	RDB   *redis.Client
}

func NewService(app *database.App) *Infra {
	return &Infra{Pools: app.Pools, Node: app.Node, RDB: app.RDB}
}
