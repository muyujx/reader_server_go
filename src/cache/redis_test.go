package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"muyu.com/reader_server_go/v1/src/config"
	"os"
	"runtime"
	"testing"
)

func initRedisEnv() {
	config.Config.Redis = config.RedisConfig{
		Host: "127.0.0.1",
	}
	InitRedis()
}

func TestConnect(t *testing.T) {
	initRedisEnv()

	ctx := context.Background()
	err := Redis().Set(ctx, "test", "test", 0).Err()
	if err != nil {
		t.Error(err)
	}

	val, err := Redis().Get(ctx, "test").Result()
	if err != nil {
		t.Error(err)
	}
	if val != "test" {
		t.Error("value not is test!")
	}
}

func TestPath(t *testing.T) {
	path, err := os.Executable()
	if err != nil {
		t.Error(err)
	}

	t.Log(path)

	_, path, _, ok := runtime.Caller(0)

	if !ok {
		t.Fail()
	}
	t.Log(path)
}

type TestRedis struct {
	Aa string `json:"aa"`
	Bb string `json:"bb"`
}

func TestObj(t *testing.T) {
	initRedisEnv()

	param := TestRedis{
		Aa: "testA",
		Bb: "testB",
	}

	ctx := context.Background()

	jsonVal, err := json.Marshal(param)
	if err != nil {
		t.Error(err)
	}
	val, err := Redis().Set(ctx, "test", jsonVal, 0).Result()
	if err != nil {
		t.Error(err)
	}
	t.Log(val)

	res, err := Redis().Get(ctx, "test").Result()

	if err != nil {
		t.Error(err)
	}

	result := TestRedis{}
	err = json.Unmarshal([]byte(res), &result)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func TestKeyNotExists(t *testing.T) {
	initRedisEnv()

	ctx := context.Background()
	res, err := Redis().Get(ctx, "aaattt").Result()

	if err == redis.Nil {
		t.Log("key not exists!")
	} else if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestZSet(t *testing.T) {
	initRedisEnv()
}
