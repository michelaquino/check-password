package context

import (
	"fmt"
	"sync"

	mgo "gopkg.in/mgo.v2"
)

var mongoSession *mgo.Session
var onceDatabase sync.Once

// GetMongoSession return a copy of mongodb session
func GetMongoSession() *mgo.Session {
	onceDatabase.Do(func() {
		var err error

		mongoSession, err = getNewMongoSession()
		if err != nil {
			errorMsg := fmt.Sprintf("Error on start database: %s", err.Error())
			panic(errorMsg)
		}
	})

	return mongoSession.Copy()
}

func getNewMongoSession() (*mgo.Session, error) {
	apiConfig := GetAPIConfig()
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:          apiConfig.MongoDBConfig.Addresses,
		ReplicaSetName: apiConfig.MongoDBConfig.ReplicaSetName,
		Username:       apiConfig.MongoDBConfig.Username,
		Password:       apiConfig.MongoDBConfig.Password,
		Database:       apiConfig.MongoDBConfig.DatabaseName,
		Timeout:        apiConfig.MongoDBConfig.Timeout,
	}

	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	return session, err
}
