package context

import (
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
)

// LogConfig represents the log configuration
type LogConfig struct {
	LogLevel string
}

// MongoConfig represents the MongoDB configuration
type MongoConfig struct {
	Addresses      []string
	DatabaseName   string
	Timeout        time.Duration
	Username       string
	Password       string
	ReplicaSetName string
}

// APIConfig represents the API configuration
type APIConfig struct {
	LogConfig     *LogConfig
	MongoDBConfig *MongoConfig
	ProxyURL      *url.URL
}

var apiConfig *APIConfig
var onceConfig sync.Once

// GetAPIConfig return the instance of the APIConfig
func GetAPIConfig() *APIConfig {
	onceConfig.Do(func() {
		apiConfig = &APIConfig{
			LogConfig:     getLogConfig(),
			MongoDBConfig: getMongoConfig(),
			ProxyURL:      getProxyURL(),
		}
	})

	return apiConfig
}

func getLogConfig() *LogConfig {
	return &LogConfig{
		LogLevel: os.Getenv("LOG_LEVEL"),
	}
}

func getMongoConfig() *MongoConfig {
	mongoURL := os.Getenv("DBAAS_MONGODB_ENDPOINT")
	pattern := regexp.MustCompile("(\\w+)://(.*):(.*)@(.*)/(.*)\\?replicaSet=(.*)$")
	match := pattern.FindAllStringSubmatch(mongoURL, 3)

	mongoConfig := MongoConfig{}
	if len(match) > 0 {
		if len(match[0]) > 2 {
			mongoConfig.Username = match[0][2]
		}

		if len(match[0]) > 3 {
			mongoConfig.Password = match[0][3]
		}

		mongoURLList := []string{}
		if len(match[0]) > 4 {
			mongoURLs := match[0][4]
			mongoURLList = strings.Split(mongoURLs, ",")
		}

		mongoConfig.Addresses = mongoURLList

		if len(match[0]) > 5 {
			mongoConfig.DatabaseName = match[0][5]
		}

		if len(match[0]) > 6 {
			mongoConfig.ReplicaSetName = match[0][6]
		}
	}

	return &mongoConfig
}

func getMongoPort() int {
	mongoPort, err := strconv.Atoi(os.Getenv("MONGO_PORT"))
	if err != nil {
		return 27017
	}

	return mongoPort
}

func getMongoTimeout() time.Duration {
	mongoTimeout, err := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
	if err != nil {
		return time.Duration(60) * time.Second
	}

	return time.Duration(mongoTimeout) * time.Second
}

func getProxyURL() *url.URL {
	proxyURL, err := url.Parse(os.Getenv("PROXY_URL"))
	if err != nil {
		log.Error("Main", "createNewRelicApp", "", "", "", "Parse proxy url from env var", err.Error(), "Error on parse proxy url")
		return nil
	}

	return proxyURL
}
