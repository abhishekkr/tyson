
## HTTP Service

* Show Help

```
$ tyson -service http -mode help

Configurable Http env-vars for Tyson are:
* TYSON_HTTP_HOST:           default("http://127.0.01:8080")
* TYSON_HTTP_AUTH_TOKEN:     default("") ## shall be value for HTTP Auth Header
* TYSON_HTTP_PATH:           default("/") ## whatever required as url path
* TYSON_HTTP_METHOD:         default("GET")
* TYSON_HTTP_REQ_BODY_FILE:  default("") ## path to file to be used as request body
* TYSON_HTTP_PARAMS:         default("") ## url params as CSV
* TYSON_HTTP_HEADERS:        default("X-REQUEST-FOR:perf,X-REQUEST-FROM:tyson") ## http headers as CSV
* TYSON_SKIP_SSL_VERIFY:     default("true")

Configurable common env-vars for Tyson are:
* TYSON_MAX_REQUESTS:       default("5000000")  ## count of max requests made
* TYSON_CONCURRENCY_LIMIT:  default("1000") ## count of max concurrent requests made

* GOLLOG_LOG_LEVEL:         default("5") # set to 0, 1 or 2 for minimal logs, as Logrus
```


* Ping, check availability

```
$ tyson -service http -mode ping

PONG
```


* Execute, default mode is `execute`

```
## http get at 127.0.0.1:8080, will make default 5000000 request with 1000 at a time
$ tyson -service redis -mode execute

## total 100 request with 10 at a time
$ TYSON_CONCURRENCY_LIMIT=10 TYSON_MAX_REQUESTS=100 tyson -service http

## http get at https://myservice.com
$ TYSON_HTTP_HOST="https://myservice.com" tyson -service http

## http get at https://myservice.com/some/path
$ TYSON_HTTP_HOST="https://myservice.com" TYSON_HTTP_PATH="/some/path" tyson -service http

## http DELETE at https://myservice.com/some/path
$ TYSON_HTTP_METHOD="DELETE" TYSON_HTTP_HOST="https://myservice.com" TYSON_HTTP_PATH="/some/path" tyson -service http

## http POST at https://myservice.com/some/path with content from /tmp/post-body
$ TYSON_HTTP_REQ_BODY_FILE="/tmp/post-body" TYSON_HTTP_METHOD="POST" TYSON_HTTP_HOST="https://myservice.com" TYSON_HTTP_PATH="/some/path" tyson -service http
```


---
