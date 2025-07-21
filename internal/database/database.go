package database

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	gocql "github.com/gocql/gocql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Pools  []*pgxpool.Pool
	Node   *snowflake.Node
	RDB    *redis.Client
	Scylla *gocql.Session
}

var (
	RDB    *redis.Client
	Scylla *gocql.Session
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("Error loading .env file", err)
	}
}

func NewApp() (*App, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Printf("snowflake node init failed: %v", err)
		return nil, err
	}
	pools := ConnectPostgres()
	rdb := ConnectRedis()
	session := ConnectScylla()

	return &App{Node: node, Pools: pools, RDB: rdb, Scylla: session}, nil
}

func ConnectPostgres() []*pgxpool.Pool {
	raw := os.Getenv("DATABASE_URLS")
	if raw == "" {
		log.Fatal("DATABASE_URLS is not set; should be comma-separated URLs")
		return nil
	}
	urls := strings.Split(raw, ",")
	pools := make([]*pgxpool.Pool, 0, 4)
	for i, u := range urls {
		pool, err := pgxpool.New(context.Background(), strings.TrimSpace(u))
		if err != nil {
			log.Fatalf("Unable to connect to database shard %d: %v", i, err)
		}
		log.Printf("Shard %d connected\n", i)
		pools = append(pools, pool)
	}
	return pools
}

func ConnectScylla() *gocql.Session {
	cluster := gocql.NewCluster("x32:9042")
	cluster.Keyspace = "chatapp"
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 10 * time.Second
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("unable to connect to Scylla: %v", err)
	}
	return session
}

func ConnectRedis() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASS")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Unable to connect to Redis: %v\n", err)
	}
	log.Println("Redis connected")
	return rdb
}

func (a *App) Close() {
	for i, pool := range a.Pools {
		pool.Close()
		log.Printf("Shard %d closed\n", i)
	}
	if a.RDB != nil {
		if err := RDB.Close(); err != nil {
			log.Printf("Error closing redis connection: %v\n", err)
		}
	}
	if a.Scylla != nil {
		a.Scylla.Close()
	}
}

// Websocket Stuff I think
func Connect() {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASS")

	Scylla = ConnectScylla()

	RDB = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword, // no password if empty
		DB:       0,             // default DB
	})

	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Unable to connect to Redis: %v\n", err)
	}
	log.Println("Redis connected")
}

func Close() {
	if RDB != nil {
		if err := RDB.Close(); err != nil {
			log.Printf("Error closing redis connection: %v\n", err)
		}
	}
}
