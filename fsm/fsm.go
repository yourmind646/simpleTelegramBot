package fsm

import (
	"context"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type FSM struct {
	rdb *redis.Client
	ctx context.Context
}

func NewFSM(addr string, password string, db int) *FSM {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr, // "localhost:6379"
		Password: password, // no password - ""
		DB:       db,  // use default DB - 0
		Protocol: 2,
	})
	return &FSM{
		rdb: rdb,
		ctx: context.Background(),
	}
}

func (fsm *FSM) SetState(user_id int64, state string) {
	user_redis_key := strconv.Itoa(int(user_id)) + "_" + "state"
	
	err := fsm.rdb.Set(fsm.ctx, user_redis_key, state, 0).Err()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Установлено состояние '%s' для пользователя id%d\n", state, user_id)
}

func (fsm *FSM) GetState(user_id int64) (string, bool) {
	user_redis_key := strconv.Itoa(int(user_id)) + "_" + "state"
	
	val, err := fsm.rdb.Get(fsm.ctx, user_redis_key).Result()
	if err != nil {
		return "", false
	}

	return val, true
}
