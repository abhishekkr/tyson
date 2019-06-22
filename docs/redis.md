
## Redis Service

* Show Help

```
$ GO111MODULE=on go run tyson.go --mode help

Configurable Redis env-vars for Tyson are:
* TYSON_REDIS_HOST:         default(127.0.01:6379)
* TYSON_REDIS_PASSWORD:     default("")
* TYSON_REDIS_DB:           default(0), allowed values: 0-15
* TYSON_REDIS_CALL:         default("set"), set/get/del/sadd/smembers
* TYSON_REDIS_KEY:          default("tyson")
* TYSON_REDIS_KEY_EXPIRY:   default(0), whatever Duration wanna give for set key expiration
* TYSON_REDIS_VALUE_PREFIX: default("peek-a-boo")
* TYSON_REDIS_VALUE_COUNT:  default(5000000)
```


* Ping, check availability

```
$ GO111MODULE=on go run tyson.go --mode ping

PONG
```


* Execute, default mode is `execute`

```
## set
$ TYSON_REDIS_Calls=set GO111MODULE=on go run tyson.go

## set 100 keys only
$ TYSON_REDIS_Calls=set TYSON_REDIS_VALUE_COUNT=100 GO111MODULE=on go run tyson.go

## set 1000 keys, with prefix 'order-number'
$ TYSON_REDIS_Calls=set TYSON_REDIS_VALUE_COUNT=100 TYSON_REDIS_KEY='order-number' GO111MODULE=on go run tyson.go

## get
$ TYSON_REDIS_Calls=get GO111MODULE=on go run tyson.go

## del
$ TYSON_REDIS_Calls=del GO111MODULE=on go run tyson.go

## sadd
$ TYSON_REDIS_Calls=sadd GO111MODULE=on go run tyson.go

## smembers
$ TYSON_REDIS_Calls=smembers GO111MODULE=on go run tyson.go

```


---
