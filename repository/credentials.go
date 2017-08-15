package repository

import "gitlab.globoi.com/michel.aquino/check-password/models"

var credentialList = []*models.Credentials{}

func SaveCredentials(credentials *models.Credentials) {
	credentialList = append(credentialList, credentials)
}

func ListCredentials() []*models.Credentials {
	return credentialList
}
