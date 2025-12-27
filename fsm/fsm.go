package fsm

import (
	"context"
	"encoding/json"
	"log"
	"sort"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type FSM struct {
	rdb *redis.Client
	ctx context.Context
}

func NewFSM(addr string, password string, db int) *FSM {
	// sort states
	for k, v := range availableStates {
		sort.Strings(v)
		availableStates[k] = v
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,     // "localhost:6379"
		Password: password, // no password - ""
		DB:       db,       // use default DB - 0
		Protocol: 2,
	})

	return &FSM{
		rdb: rdb,
		ctx: context.Background(),
	}
}

func getUserKey(user_id int64, scope string) string {
	user_redis_key := strconv.Itoa(int(user_id)) + "_" + scope

	return user_redis_key
}

func isStateValueExists(stateClass, stateValue string) bool {
	stateValues, ok := availableStates[stateClass]
	if !ok {
		return false
	}

	i := sort.SearchStrings(stateValues, stateValue)
	return i < len(stateValues) && stateValues[i] == stateValue
}

func (fsm *FSM) SetState(user_id int64, stateClass string, stateValue string) {
	user_redis_key := getUserKey(user_id, "state")
	state := stateClass + "." + stateValue

	if _, ok := availableStates[stateClass]; !ok {
		log.Printf("Ошибка установки состояния `%s` для пользователя id%d: класс не найден", state, user_id)
		return
	}
	if !isStateValueExists(stateClass, stateValue) {
		log.Printf("Ошибка установки состояния `%s` для пользователя id%d: значение состояния не найдено", state, user_id)
		return
	}

	err := fsm.rdb.Set(fsm.ctx, user_redis_key, state, 0).Err()
	if err != nil {
		log.Panic(err)
	}

	// debug
	log.Printf("Установлено состояние `%s` для пользователя id%d\n", state, user_id)
}

func (fsm *FSM) GetState(user_id int64) (string, bool) {
	user_redis_key := getUserKey(user_id, "state")

	val, err := fsm.rdb.Get(fsm.ctx, user_redis_key).Result()
	if err != nil {
		return "", false
	}

	return val, true
}

func (fsm *FSM) GetData(user_id int64) (map[string]string, bool) {
	user_redis_key := getUserKey(user_id, "data")

	val, err := fsm.rdb.Get(fsm.ctx, user_redis_key).Result()
	if err != nil {
		return map[string]string{}, false
	}

	var result map[string]string
	data_bytes := []byte(val)

	if json.Valid(data_bytes) {
		json.Unmarshal(data_bytes, &result)
		return result, true
	} else {
		return map[string]string{}, false
	}
}

func (fsm *FSM) UpdateData(user_id int64, data map[string]string) {
	user_redis_key := getUserKey(user_id, "data")
	old_data, _ := fsm.GetData(user_id)

	for k, v := range data {
		old_data[k] = v
	}

	result, err := json.Marshal(old_data)
	if err != nil {
		log.Panic(err)
	}

	seterr := fsm.rdb.Set(fsm.ctx, user_redis_key, result, 0).Err()
	if seterr != nil {
		log.Panic(seterr)
	}
}

func (fsm *FSM) ClearData(user_id int64) {
	user_redis_key := getUserKey(user_id, "data")

	seterr := fsm.rdb.Set(fsm.ctx, user_redis_key, "{}", 0).Err()
	if seterr != nil {
		log.Panic(seterr)
	}
}
