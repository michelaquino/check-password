package context

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// LogConfig represents the log configuration
type LogConfig struct {
	LogLevel string
}

// MongoConfig represents the MongoDB configuration
type MongoConfig struct {
	Address      string
	DatabaseName string
	Timeout      time.Duration
	Username     string
	Password     string
}

// APIConfig represents the API configuration
type APIConfig struct {
	LogConfig     *LogConfig
	MongoDBConfig *MongoConfig
}

var apiConfig *APIConfig
var onceConfig sync.Once

// GetAPIConfig return the instance of the APIConfig
func GetAPIConfig() *APIConfig {
	onceConfig.Do(func() {
		apiConfig = &APIConfig{
			LogConfig:     getLogConfig(),
			MongoDBConfig: getMongoConfig(),
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

	pattern := regexp.MustCompile("(\\w+)://(.*):(.*)@(.*):([0-9]+)/(.*)?$")
	match := pattern.FindAllStringSubmatch(mongoURL, 3)

	mongoConfig := MongoConfig{}
	if len(match) > 0 {
		fmt.Println("match: ", match)
		if len(match[0]) > 2 {
			mongoConfig.Username = match[0][2]
		}

		if len(match[0]) > 3 {
			mongoConfig.Password = match[0][3]
		}

		mongoURL := ""
		if len(match[0]) > 4 {
			mongoURL = match[0][4]
		}

		mongoPort := 27017
		if len(match[0]) > 5 {
			port, err := strconv.Atoi(match[0][5])
			if err == nil {
				mongoPort = port
			}
		}
		mongoConfig.Address = fmt.Sprintf("%s:%d", mongoURL, mongoPort)

		if len(match[0]) > 6 {
			mongoConfig.DatabaseName = match[0][6]
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
