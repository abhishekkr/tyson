package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/abhishekkr/gol/golconv"
	"github.com/abhishekkr/gol/golenv"
	"github.com/abhishekkr/gol/gollog"
	"github.com/go-redis/redis"
)

var (
	RedisHost      = golenv.OverrideIfEnv("TYSON_REDIS_HOST", "127.0.0.1:6379")
	RedisPassword  = golenv.OverrideIfEnv("TYSON_REDIS_PASSWORD", "")
	RedisDB        = golenv.OverrideIfEnv("TYSON_REDIS_DB", "0")
	RedisKey       = golenv.OverrideIfEnv("TYSON_REDIS_KEY", "tyson")
	RedisKeyExpiry = time.Duration(golconv.StringToInt(golenv.OverrideIfEnv("TYSON_REDIS_KEY_EXPIRY", ""), 0))
	RedisValPrefix = golenv.OverrideIfEnv("TYSON_REDIS_VALUE_PREFIX", "peek-a-boo")
	RedisCall      = golenv.OverrideIfEnv("TYSON_REDIS_CALL", "set")

	RedisCalls = map[string]func(int, *sync.WaitGroup){}
)

/*
init registers Redis to ServiceEngines.
*/
func init() {
	db := golconv.StringToInt(RedisDB, 0)
	if db < 0 || db > 15 {
		gollog.Warnf("wrong redis db: %d available, using 0", db)
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
	Client     *redis.Client
	StartTime  time.Time
	EndTime    time.Time
	ErrorCount uint64
	sync.Mutex
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
	gollog.Infof("starting to run %d count of %s calls", MaxRequests, RedisCall)
	svc.StartTime = time.Now()
	for idx := 0; idx < MaxRequests; {
		var wg sync.WaitGroup
		for limit := ConcurrencyLimit; limit > 0 && idx < MaxRequests; limit-- {
			wg.Add(1)
			go RedisCalls[RedisCall](idx, &wg)
			idx++
		}
		wg.Wait()
	}
	svc.EndTime = time.Now()
	gollog.Infof("done %s for %d entries\n", RedisCall, MaxRequests)
	fmt.Printf(`
Started: %s
Finished: %s

Total Requests: %d
Total Errors:   %d
`, svc.StartTime.Format(time.RFC3339),
		svc.EndTime.Format(time.RFC3339),
		MaxRequests,
		svc.ErrorCount)
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
	help()
}

func (svc *RedisService) check_error(index int, err error) {
	if err == nil {
		return
	}
	svc.Lock()
	svc.ErrorCount++
	svc.Unlock()
	gollog.Warnf("failed at index %d for: %s", index, err)
}

func (svc *RedisService) set(index int, wg *sync.WaitGroup) {
	key := fmt.Sprintf("%s-%d", RedisKey, index)
	val := fmt.Sprintf("%s-%d", RedisValPrefix, index)
	result := svc.Client.Set(key, val, RedisKeyExpiry)
	svc.check_error(index, result.Err())
	wg.Done()
	gollog.Infof("done %s for index: %d", RedisCall, index)
}

func (svc *RedisService) get(index int, wg *sync.WaitGroup) {
	key := fmt.Sprintf("%s-%d", RedisKey, index)
	result := svc.Client.Get(key)
	svc.check_error(index, result.Err())
	wg.Done()
	gollog.Infof("done %s for index: %d", RedisCall, index)
}

func (svc *RedisService) del(index int, wg *sync.WaitGroup) {
	key := fmt.Sprintf("%s-%d", RedisKey, index)
	result := svc.Client.Del(key)
	svc.check_error(index, result.Err())
	wg.Done()
	gollog.Infof("done %s for index: %d", RedisCall, index)
}

func (svc *RedisService) sadd(index int, wg *sync.WaitGroup) {
	val := fmt.Sprintf("%s-%d", RedisValPrefix, index)
	result := svc.Client.SAdd(RedisKey, val)
	svc.check_error(index, result.Err())
	wg.Done()
	gollog.Infof("done %s for index: %d", RedisCall, index)
}

func (svc *RedisService) smembers(index int, wg *sync.WaitGroup) {
	result := svc.Client.SMembers(RedisKey)
	svc.check_error(index, result.Err())
	wg.Done()
	gollog.Infof("done %s for index: %d", RedisCall, index)
}
