package service

import (
	"fmt"
	"hash/fnv"

	"github.com/gocql/gocql"
	"github.com/luka-sijic/flux/internal/database"
	"github.com/luka-sijic/flux/internal/models"

	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Service interface {
	CreateUser(user *models.UserDTO) bool
	LoginUser(user *models.UserDTO) *models.User
	Profile(username string) bool

	AddFriend(username string, user *models.FriendDTO) bool
	GetFriends(username string) []models.FriendDTO
	GetRequests(username string) []models.FriendDTO
	FriendResponse(username string, action *models.FriendActionDTO) bool
	GetLog(username1, username2 string) []models.Messages
}

type Infra struct {
	Pools  []*pgxpool.Pool
	Node   *snowflake.Node
	RDB    *redis.Client
	Scylla *gocql.Session
}

func (infra *Infra) GetShardPool(key snowflake.ID) *pgxpool.Pool {
	//id := (key >> 12) & ((1 << 10) - 1)
	//fmt.Println("\033[32m POOL ID: ", id, " \033[0m")
	h := fnv.New32a()
	h.Write([]byte(key.String()))
	idx := int(h.Sum32()) % len(infra.Pools)
	return infra.Pools[idx]
}

func sortUsernames(user1, user2 string) string {
	var key string
	if user1 > user2 {
		key = fmt.Sprintf("%s:%s", user1, user2)
	} else {
		key = fmt.Sprintf("%s:%s", user2, user1)
	}
	return key
}

func NewService(app *database.App) *Infra {
	return &Infra{Pools: app.Pools, Node: app.Node, RDB: app.RDB, Scylla: app.Scylla}
}
