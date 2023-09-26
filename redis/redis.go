package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"new_project/config"
	"os"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

type redisLog struct {
	Type 	string	`json:"type"`
	Key 	string	`json:"key"`
}

func RedisInit() {
	opt, err := redis.ParseURL(os.Getenv("REDIS"))
	if err != nil {
		log.Fatal(err)
	}
	rdb = redis.NewClient(opt)
}

func Get(key string) ([]byte, error) {
	ctx := context.Background()
	val, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("key does not exists")
			return nil, nil
		}
		return nil, err
	}
	log := &redisLog{
		Type: "Get",
		Key: key,
	}
	jsonstr, _ := json.Marshal(log)
	fmt.Println(string(jsonstr))
	return val, nil
}

func Set(key string, val string) error {
	ctx := context.Background()
	err := rdb.Set(ctx, key, val, config.RedisTimeDuration).Err()
	if err != nil {
		return err
	}
	log := &redisLog{
		Type: "Set",
		Key: key,
	}
	jsonstr, _ := json.Marshal(log)
	fmt.Println(string(jsonstr))
	return nil
}

func Delete(key string) {
	ctx := context.Background()
	rdb.Del(ctx, key)
	log := &redisLog{
		Type: "Delete",
		Key: key,
	}
	jsonstr, _ := json.Marshal(log)
	fmt.Println(string(jsonstr))
}
