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
	log := GetLogger()
	log.Debug("GetMongoSession", "", "Start to get MongoDB session")

	onceDatabase.Do(func() {
		var err error

		mongoSession, err = getNewMongoSession()
		if err != nil {
			errorMsg := fmt.Sprintf("Error on start database: %s", err.Error())
			log.Error("GetMongoSession", "Error", errorMsg)
			panic(errorMsg)
		}
	})

	log.Debug("GetMongoSession", "Success", "MongoDB session getted with success")
	return mongoSession.Copy()
}

func getNewMongoSession() (*mgo.Session, error) {
	log := GetLogger()

	apiConfig := GetAPIConfig()
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:          apiConfig.MongoDBConfig.Addresses,
		ReplicaSetName: apiConfig.MongoDBConfig.ReplicaSetName,
		Username:       apiConfig.MongoDBConfig.Username,
		Password:       apiConfig.MongoDBConfig.Password,
		Database:       apiConfig.MongoDBConfig.DatabaseName,
		Timeout:        apiConfig.MongoDBConfig.Timeout,
	}

	log.Debug("getNewMongoSession", "", "Start to dial to MongoDB")
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	return session, err
}
