package main

import (
	"flag"
	"log"

	"github.com/abhishekkr/tyson/service"
)

func main() {
	serviceName := flag.String("service", "redis", "tyson run mode: redis")
	mode := flag.String("mode", "execute", "service run mode: execute/ping/help")
	flag.Parse()
	svc := service.ServiceEngines[*serviceName]
	if svc == nil {
		log.Fatalf("TYSON_MODE env var only allows 'redis' as of now")
	}
	switch *mode {
	case "ping":
		svc.Ping()
	case "execute":
		svc.Execute()
	default:
		svc.Help()
	}
}
