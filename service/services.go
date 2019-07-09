package service

import (
	"fmt"
	"net"
	"time"

	"github.com/abhishekkr/gol/golconv"
	"github.com/abhishekkr/gol/golenv"
)

var (
	ConcurrencyLimit = golconv.StringToInt(golenv.OverrideIfEnv("TYSON_CONCURRENCY_LIMIT", ""), 1000)
	MaxRequests      = golconv.StringToInt(golenv.OverrideIfEnv("TYSON_MAX_REQUESTS", ""), 5000000)
)

type Service interface {
	Ping() error
	Execute()
	Help()
}

/*
ServiceEngines acts as map for all available service-engines.
*/
var ServiceEngines = make(map[string]Service)

/*
RegisterService gets used by adapters to register themselves.
*/
func RegisterService(name string, service Service) {
	ServiceEngines[name] = service
}

/*
GetService gets used by client to fetch a required db-engine.
*/
func GetService(name string) Service {
	return ServiceEngines[name]
}

/*
Generic TCP availability for Pong, if needed
*/
func ping(hostPort string) error {
	timeout := time.Duration(3 * time.Second)
	conn, err := net.DialTimeout("tcp", hostPort, timeout)
	defer conn.Close()
	if err == nil {
		fmt.Println("PONG")
	} else {
		fmt.Println(err.Error())
	}
	return err
}

func help() {
	fmt.Println(`Configurable common env-vars for Tyson are:
* TYSON_MAX_REQUESTS:       default("5000000")  ## count of max requests made
* TYSON_CONCURRENCY_LIMIT:  default("1000") ## count of max concurrent requests made

* GOLLOG_LOG_LEVEL:         default("5") # set to 0, 1 or 2 for minimal logs, as Logrus
`)
}
