package repository

import (
	"gitlab.globoi.com/michel.aquino/check-password/context"
	"gitlab.globoi.com/michel.aquino/check-password/models"
	"gopkg.in/mgo.v2/bson"
)

func SaveCredentials(credentials *models.Credentials) error {
	log := context.GetLogger()

	dbSession := context.GetMongoSession()
	defer dbSession.Close()

	connection := dbSession.DB("check-password").C("credentials")

	if err := connection.Insert(&credentials); err != nil {
		log.Error("Save credentials on DB", "Error", err.Error())
		return err
	}

	log.Info("Save credentials on DB", "Success", "")
	return nil
}

func ListCredentials(onlyHackedCredentials bool) ([]models.Credentials, error) {
	log := context.GetLogger()

	dbSession := context.GetMongoSession()
	defer dbSession.Close()

	connection := dbSession.DB("check-password").C("credentials")

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
