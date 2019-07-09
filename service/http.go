package service

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/abhishekkr/gol/golconv"
	"github.com/abhishekkr/gol/golenv"
	"github.com/abhishekkr/gol/golfilesystem"
	"github.com/abhishekkr/gol/golhttpclient"
	"github.com/abhishekkr/gol/gollog"
)

var (
	HttpHost          = golenv.OverrideIfEnv("TYSON_HTTP_HOST", "http://127.0.0.1:8080")
	HttpAuthToken     = golenv.OverrideIfEnv("TYSON_HTTP_AUTH_TOKEN", "")
	HttpPath          = golenv.OverrideIfEnv("TYSON_HTTP_PATH", "/")
	HttpMethod        = golenv.OverrideIfEnv("TYSON_HTTP_METHOD", "GET")
	HttpReqBodyFile   = golenv.OverrideIfEnv("TYSON_HTTP_REQ_BODY_FILE", "")
	HttpParams        = golenv.OverrideIfEnv("TYSON_HTTP_PARAMS", "")
	HttpHeaders       = golenv.OverrideIfEnv("TYSON_HTTP_HEADERS", "X-REQUEST-FOR:perf,X-REQUEST-FROM:tyson")
	HttpSkipSSLVerify = golconv.StringToBool(golenv.OverrideIfEnv("TYSON_SKIP_SSL_VERIFY", ""), true)
)

/*
init registers Http to ServiceEngines.
*/
func init() {
	golhttpclient.SkipSSLVerify = HttpSkipSSLVerify
	if _, _, err := net.SplitHostPort(HttpHost); err == nil {
		HttpHost = fmt.Sprintf("http://%s", HttpHost)
	}

	client := golhttpclient.HTTPRequest{
		Method:      HttpMethod,
		Url:         HttpHost,
		GetParams:   nil,
		HTTPHeaders: map[string]string{"User-Agent": "tyson/v0.2+ perf or data-ops"},
	}
	if HttpPath != "" && HttpPath != "/" {
		if strings.HasSuffix(client.Url, "/") {
			client.Url += "/"
		}
		client.Url += HttpPath
	}
	if HttpAuthToken != "" {
		client.HTTPHeaders["Authorization"] = HttpAuthToken
	}
	if HttpParams != "" {
		client.GetParams = map[string]string{}
		for _, param := range strings.Split(HttpParams, ",") {
			paramKeyVal := strings.Split(param, "=")
			client.GetParams[paramKeyVal[0]] = strings.Join(paramKeyVal[1:], "=")
		}
	}
	if HttpHeaders != "" {
		for _, header := range strings.Split(HttpHeaders, ",") {
			headerKeyVal := strings.Split(header, ":")
			client.HTTPHeaders[headerKeyVal[0]] = strings.Join(headerKeyVal[1:], ":")
		}
	}
	if HttpReqBodyFile != "" {
		buffer, err := golfilesystem.FileBuffer(HttpReqBodyFile)
		if err == nil {
			client.Body = buffer
		} else {
			gollog.Errf("error reding request body file: %s", HttpReqBodyFile)
		}
	}
	httpService := HttpService{Client: &client}
	RegisterService("http", &httpService)
}

type HttpService struct {
	Client     *golhttpclient.HTTPRequest
	StartTime  time.Time
	EndTime    time.Time
	ErrorCount uint64
	sync.Mutex
}

func (svc *HttpService) check_error(index int, err error) {
	if err == nil {
		return
	}
	svc.Lock()
	svc.ErrorCount++
	svc.Unlock()
	gollog.Warnf("failed at index %d for: %s", index, err)
}

func (svc *HttpService) Ping() error {
	u, err := url.Parse(HttpHost)
	if err != nil {
		fmt.Println("failed url parse for http host:", HttpHost)
	}
	hostPort := u.Host
	if _, _, err := net.SplitHostPort(u.Host); err != nil {
		if u.Scheme == "https" {
			hostPort = fmt.Sprintf("%s:443", hostPort)
		} else if u.Scheme == "http" {
			hostPort = fmt.Sprintf("%s:80", hostPort)
		} else {
			panic("unamange url scheme, supported http/https")
		}
	}
	return ping(hostPort)
}

func (svc *HttpService) Execute() {
	gollog.Infof("starting to run %d count of [%s](%s) calls",
		MaxRequests, svc.Client.Method, svc.Client.Url)
	svc.StartTime = time.Now()
	for idx := 0; idx < MaxRequests; {
		var wg sync.WaitGroup
		for limit := ConcurrencyLimit; limit > 0 && idx < MaxRequests; limit-- {
			wg.Add(1)
			go svc.request(idx, &wg)
			idx++
		}
		wg.Wait()
	}
	svc.EndTime = time.Now()
	gollog.Infof("done [%s](%s) for %d entries\n",
		svc.Client.Method, svc.Client.Url, MaxRequests)
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

func (svc *HttpService) Help() {
	fmt.Println(`Configurable Http env-vars for Tyson are:
* TYSON_HTTP_HOST:           default("http://127.0.01:8080")
* TYSON_HTTP_AUTH_TOKEN:     default("") ## shall be value for HTTP Auth Header
* TYSON_HTTP_PATH:           default("/") ## whatever required as url path
* TYSON_HTTP_METHOD:         default("GET")
* TYSON_HTTP_REQ_BODY_FILE:  default("") ## path to file to be used as request body
* TYSON_HTTP_PARAMS:         default("") ## url params as CSV
* TYSON_HTTP_HEADERS:        default("X-REQUEST-FOR:perf,X-REQUEST-FROM:tyson") ## http headers as CSV
* TYSON_SKIP_SSL_VERIFY:     default("true")
`)
	help()
}

func (svc *HttpService) request(index int, wg *sync.WaitGroup) {
	response, err := svc.Client.Response()
	if err != nil {
		svc.check_error(index, err)
		gollog.Warnf("failed at index %d for: %s", index, err)
	} else if response.StatusCode > 399 {
		svc.check_error(index, fmt.Errorf("HTTP status: %d", response.StatusCode))
		gollog.Warnf("failed at index %d for: HTTP status %d", index, response.StatusCode)
	}
	wg.Done()
	gollog.Infof("done at index %d for: %s | %s", index, HttpMethod, svc.Client.Url)
}
