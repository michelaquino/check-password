package context

import (
	"fmt"
	"os"
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
	mongoURL := os.Getenv("MONGO_URL")
	mongoPort := getMongoPort()

	mongoAddress := fmt.Sprintf("%s:%d", mongoURL, mongoPort)
	mongoDatabaseName := os.Getenv("MONGO_DATABASE_NAME")
	mongoTimeout := getMongoTimeout()
	mongoUserName := os.Getenv("MONGO_DATABASE_USERNAME")
	mongoPassword := os.Getenv("MONGO_DATABASE_PASSWORD")

	return &MongoConfig{
		Address:      mongoAddress,
		DatabaseName: mongoDatabaseName,
		Timeout:      mongoTimeout,
		Username:     mongoUserName,
		Password:     mongoPassword,
	}
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
