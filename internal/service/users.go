package service

import (
	"context"

	"github.com/luka-sijic/flux/internal/database"
	"github.com/luka-sijic/flux/internal/models"
	"github.com/luka-sijic/flux/pkg/hash"

	"log"

	"github.com/bwmarrin/snowflake"
)

func (infra *Infra) CreateUser(user *models.UserDTO) bool {
	res, err := infra.RDB.Do(context.Background(), "BF.EXISTS", "users", user.Username).Bool()
	if err != nil || res {
		log.Println(err)
		return false
	}

	hashedPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		log.Println(err)
		return false
	}

	id := infra.Node.Generate()
	db := database.GetShardPool(infra.Pools, id)

	_, err = db.Exec(context.Background(), "INSERT INTO users (id, username, password) VALUES ($1,$2,$3)", id, user.Username, hashedPassword)
	if err != nil {
		log.Println("Failed to create user: ", err)
		return false
	}

	infra.RDB.HSet(context.Background(), "user:id_map", user.Username, id.String())
	infra.RDB.HSet(context.Background(), "user:username_map", id.String(), user.Username)
	infra.RDB.Do(context.Background(), "BF.ADD", "users", user.Username)

	return true
}

func (infra *Infra) LoginUser(user *models.UserDTO) *models.User {
	idStr, err := infra.RDB.HGet(context.Background(), "user:id_map", user.Username).Result()
	if err != nil {
		log.Println("Error 1:", err)
		return nil
	}
	id, err := snowflake.ParseString(idStr)
	if err != nil {
		log.Println("Error 2:", err)
		return nil
	}
	db := database.GetShardPool(infra.Pools, id)

	var storedHash string
	err = db.QueryRow(context.Background(), "SELECT password FROM users WHERE username=$1", user.Username).Scan(&storedHash)
	if err != nil || !hash.CheckPasswordHash(user.Password, storedHash) {
		log.Println("Error 3:", err)
		return nil
	}

	return &models.User{ID: idStr, Username: user.Username}
}
