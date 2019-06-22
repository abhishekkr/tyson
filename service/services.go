package service

import (
	"github.com/abhishekkr/gol/golconv"
	"github.com/abhishekkr/gol/golenv"
)

var (
	ConcurrencyLimit = golconv.StringToInt(golenv.OverrideIfEnv("TYSON_CONCURRENCY_LIMIT", ""), 1000)
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
