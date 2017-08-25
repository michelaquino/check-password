package repository

import (
	"github.com/michelaquino/check-password/context"
	"github.com/michelaquino/check-password/models"
	"gopkg.in/mgo.v2/bson"
)

func SaveCredentials(credentials *models.Credentials) error {
	log := context.GetLogger()
	config := context.GetAPIConfig()

	dbSession := context.GetMongoSession()
	defer dbSession.Close()

	connection := dbSession.DB(config.MongoDBConfig.DatabaseName).C("credentials")

	if err := connection.Insert(&credentials); err != nil {
		log.Error("Save credentials on DB", "Error", err.Error())
		return err
	}

	log.Info("Save credentials on DB", "Success", "")
	return nil
}

func ListCredentials(onlyHackedCredentials bool) ([]models.Credentials, error) {
	log := context.GetLogger()
	config := context.GetAPIConfig()

	dbSession := context.GetMongoSession()
	defer dbSession.Close()

	connection := dbSession.DB(config.MongoDBConfig.DatabaseName).C("credentials")

	query := bson.M{}
	if onlyHackedCredentials {
		query = bson.M{
			"$or": []bson.M{
				bson.M{"emailBreached": true},
				bson.M{"passwordPwned": true},
			},
		}
	}

	credentialsList := []models.Credentials{}
	if err := connection.Find(query).All(&credentialsList); err != nil {
		log.Error("Get all credentials on DB", "Error", err.Error())
		return nil, err
	}

	log.Info("Get all credentials on DB", "Success", "")
	return credentialsList, nil
}

func GetUnprocessedCredentials() ([]models.Credentials, error) {
	log := context.GetLogger()
	config := context.GetAPIConfig()

	log.Debug("Start to get unprocess credentials on DB", "", "")

	dbSession := context.GetMongoSession()
	defer dbSession.Close()

	connection := dbSession.DB(config.MongoDBConfig.DatabaseName).C("credentials")
	query := bson.M{"passwordProcessed": false}

	credentialsList := []models.Credentials{}
	if err := connection.Find(query).All(&credentialsList); err != nil {
		log.Error("Get unprocess credentials on DB", "Error", err.Error())
		return nil, err
	}

	log.Info("Get unprocess credentials on DB", "Success", "")
	return credentialsList, nil
}

func UpdateCredentialsProcessed(md5Processed, md5Hacked []string) error {
	log := context.GetLogger()
	config := context.GetAPIConfig()

	log.Debug("Start to update processed credentials on DB", "", "")

	dbSession := context.GetMongoSession()
	defer dbSession.Close()
	connection := dbSession.DB(config.MongoDBConfig.DatabaseName).C("credentials")

	seletorMd5Processed := bson.M{"passwordMD5Hash": bson.M{"$in": md5Processed}}
	updateMd5Processed := bson.M{"$set": bson.M{"passwordProcessed": true}}

	if _, err := connection.UpdateAll(seletorMd5Processed, updateMd5Processed); err != nil {
		log.Error("Update credentials processed", "Error", err.Error())
		return err
	}
	log.Info("Update credentials processed", "Success", "")

	seletorMd5Hacked := bson.M{"passwordMD5Hash": bson.M{"$in": md5Hacked}}
	updateMd5Hacked := bson.M{"$set": bson.M{"passwordMD5HashHacked": true}}

	if _, err := connection.UpdateAll(seletorMd5Hacked, updateMd5Hacked); err != nil {
		log.Error("Update credentials hacked", "Error", err.Error())
		return err
	}
	log.Info("Update credentials hacked", "Success", "")

	return nil
}
