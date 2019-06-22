package service

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/abhishekkr/gol/golconv"
	"github.com/abhishekkr/gol/golenv"
	"github.com/go-redis/redis"
)

var (
	RedisHost       = golenv.OverrideIfEnv("TYSON_REDIS_HOST", "127.0.0.1:6379")
	RedisPassword   = golenv.OverrideIfEnv("TYSON_REDIS_PASSWORD", "")
	RedisDB         = golenv.OverrideIfEnv("TYSON_REDIS_DB", "0")
	RedisKey        = golenv.OverrideIfEnv("TYSON_REDIS_KEY", "tyson")
	RedisKeyExpiry  = time.Duration(golconv.StringToInt(golenv.OverrideIfEnv("TYSON_REDIS_KEY_EXPIRY", ""), 0))
	RedisValPrefix  = golenv.OverrideIfEnv("TYSON_REDIS_VALUE_PREFIX", "peek-a-boo")
	RedisValueCount = golenv.OverrideIfEnv("TYSON_REDIS_VALUE_COUNT", "5000000")
	RedisCall       = golenv.OverrideIfEnv("TYSON_REDIS_Calls", "set")

	RedisCalls = map[string]func(int, *sync.WaitGroup){}
)

/*
init registers Redis to ServiceEngines.
*/
func init() {
	db := golconv.StringToInt(RedisDB, 0)
	if db < 0 || db > 15 {
		log.Printf("[warning] wrong redis db: %d available, using 0", db)
		db = 0
	}

	client := redis.NewClient(&redis.Options{
		Addr:     RedisHost,
		Password: RedisPassword,
		DB:       db,
	})
	redisService := RedisService{Client: client}
	RedisCalls["set"] = redisService.set
	RedisCalls["get"] = redisService.get
	RedisCalls["del"] = redisService.del
	RedisCalls["sadd"] = redisService.sadd
	RedisCalls["smembers"] = redisService.smembers

	RegisterService("redis", &redisService)
}

type RedisService struct {
	Client *redis.Client
}

func (svc *RedisService) Ping() error {
	pong, err := svc.Client.Ping().Result()
	if err == nil {
		fmt.Println(pong)
	} else {
		fmt.Println(err.Error())
	}
	return err
}

func (svc *RedisService) Execute() {
	valCount := golconv.StringToInt(RedisValueCount, 5000000)
	fmt.Printf("starting to run %d count of %s calls", valCount, RedisCall)
	for idx := 0; idx < valCount; {
		var wg sync.WaitGroup
		for limit := 1000; limit > 0 && idx < valCount; limit-- {
			wg.Add(1)
			go RedisCalls[RedisCall](idx, &wg)
			idx++
		}
		wg.Wait()
	}
	fmt.Printf("done %s for %d entries\n", RedisCall, valCount)
}

func (svc *RedisService) Help() {
	fmt.Println(`Configurable Redis env-vars for Tyson are:
* TYSON_REDIS_HOST:         default(127.0.01:6379)
* TYSON_REDIS_PASSWORD:     default("")
* TYSON_REDIS_DB:           default(0), allowed values: 0-15
* TYSON_REDIS_CALL:         default("set"), set/get/del/sadd/smembers
* TYSON_REDIS_KEY:          default("tyson")
* TYSON_REDIS_KEY_EXPIRY:   default(0), whatever Duration wanna give for set key expiration
* TYSON_REDIS_VALUE_PREFIX: default("peek-a-boo")
* TYSON_REDIS_VALUE_COUNT:  default(5000000)
`)
}

func (svc *RedisService) set(index int, wg *sync.WaitGroup) {
	key := fmt.Sprintf("%s-%d", RedisKey, index)
	val := fmt.Sprintf("%s-%d", RedisValPrefix, index)
	result := svc.Client.Set(key, val, RedisKeyExpiry)
	if result.Err() != nil {
		log.Printf("[error] %s\n", result.Err())
	}
	wg.Done()
}

func (svc *RedisService) get(index int, wg *sync.WaitGroup) {
	key := fmt.Sprintf("%s-%d", RedisKey, index)
	result := svc.Client.Get(key)
	if result.Err() != nil {
		log.Printf("[error] %s\n", result.Err())
	}
	wg.Done()
}

func (svc *RedisService) del(index int, wg *sync.WaitGroup) {
	key := fmt.Sprintf("%s-%d", RedisKey, index)
	result := svc.Client.Del(key)
	if result.Err() != nil {
		log.Printf("[error] %s\n", result.Err())
	}
	wg.Done()
}

func (svc *RedisService) sadd(index int, wg *sync.WaitGroup) {
	val := fmt.Sprintf("%s-%d", RedisValPrefix, index)
	result := svc.Client.SAdd(RedisKey, val)
	if result.Err() != nil {
		log.Printf("[error] %s\n", result.Err())
	}
	wg.Done()
}

func (svc *RedisService) smembers(index int, wg *sync.WaitGroup) {
	result := svc.Client.SMembers(RedisKey)
	if result.Err() != nil {
		log.Printf("[error] failed for index:%s\n%s\n", index, result.Err())
	}
	wg.Done()
}
