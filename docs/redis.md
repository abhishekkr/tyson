
## Redis Service

* Show Help

```
$ tyson -service redis -mode help

Configurable Redis env-vars for Tyson are:
* TYSON_REDIS_HOST:         default(127.0.01:6379)
* TYSON_REDIS_PASSWORD:     default("")
* TYSON_REDIS_DB:           default(0), allowed values: 0-15
* TYSON_REDIS_CALL:         default("set"), set/get/del/sadd/smembers
* TYSON_REDIS_KEY:          default("tyson")
* TYSON_REDIS_KEY_EXPIRY:   default(0), whatever Duration wanna give for set key expiration
* TYSON_REDIS_VALUE_PREFIX: default("peek-a-boo")

Configurable common env-vars for Tyson are:
* TYSON_MAX_REQUESTS:       default(5000000), maximum requests to be made
* TYSON_CONCURRENCY_LIMIT:  default(1000), concurrent request made to achieve max-requests

* GOLLOG_LOG_LEVEL:         default("5") # set to 0, 1 or 2 for minimal logs, as Logrus
```


* Ping, check availability

```
$ tyson -service redis -mode ping

PONG
```


* Execute, default mode is `execute`

```
## set
$ TYSON_REDIS_Calls=set tyson -service redis -mode execute

## set 100 keys only
$ TYSON_REDIS_Calls=set TYSON_REDIS_VALUE_COUNT=100 tyson -service redis -mode execute

## set 1000 keys, with prefix 'order-number'
$ TYSON_REDIS_Calls=set TYSON_REDIS_VALUE_COUNT=100 TYSON_REDIS_KEY='order-number' tyson -service redis -mode execute

## get
$ TYSON_REDIS_Calls=get GO111MODULE=on tyson -service redis -mode execute

## del
$ TYSON_REDIS_Calls=del tyson -service redis -mode execute

## sadd
$ TYSON_REDIS_Calls=sadd tyson -service redis -mode execute

## smembers
$ TYSON_REDIS_Calls=smembers tyson -service redis -mode execute
```


---
