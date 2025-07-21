package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func (a *App) Init() error {
	sql := `
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY,
			username VARCHAR(255) NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);

		CREATE TABLE IF NOT EXISTS friends (
			user_id BIGINT NOT NULL,
			requester_id BIGINT NOT NULL,
			STATUS TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT now(),
			updated_at TIMESTAMP NOT NULL DEFAULT now(),
			PRIMARY KEY (user_id, requester_id)
		);
	`

	for i, pool := range a.Pools {
		if err := execOnShard(pool, sql); err != nil {
			return fmt.Errorf("shard %d: %w", i, err)
		}
	}

	csql := []string{`
		CREATE TABLE IF NOT EXISTS conversations (
  			conv_id     BIGINT PRIMARY KEY,
  			kind        text,         
  			title       text,
  			created_at  timestamp,
  			owner_id    bigint,
  			settings    map<text,text> 
		);`,
		`
		CREATE TABLE IF NOT EXISTS dm_lookup (
  			user_a   bigint,
  			user_b   bigint,
  			conv_id  bigint,
  			PRIMARY KEY ((user_a, user_b))
		);`,
		`
		CREATE TABLE IF NOT EXISTS participants_by_conv (
  			conv_id   bigint,
  			user_id   bigint,
  			joined_at timestamp,
  			role      text,
  			PRIMARY KEY ((conv_id), user_id)
		);`, `
		CREATE TABLE IF NOT EXISTS convs_by_user (
  			user_id   bigint,
  			conv_id   bigint,
  			joined_at timestamp,
  			role      text,
  			PRIMARY KEY ((user_id), conv_id)
		);`, `
		CREATE TABLE IF NOT EXISTS messages_by_conv (
  			conv_id   bigint,
  			ts        bigint,        
  			sender_id bigint,
  			body      text,
  			PRIMARY KEY ((conv_id), ts)
		) WITH CLUSTERING ORDER BY (ts DESC);
	`}
	for _, q := range csql {
		if err := a.Scylla.Query(q).Exec(); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func execOnShard(pool *pgxpool.Pool, sql string) error {
	_, err := pool.Exec(context.Background(), sql)
	return err
}
