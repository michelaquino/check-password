package models

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.globoi.com/michel.aquino/check-password/context"

	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Email              string `form:"email" bson:"email,omitempty"`
	EmailLeaked        bool   `bson:"emailLeaked,omitempty"`
	Password           string `form:"password"`
	PasswordMD5Hash    string `bson:"passwordMD5Hash,omitempty"`
	PasswordSha1Hash   string `bson:"passwordSha1Hash,omitempty"`
	PasswordSha256Hash string `bson:"passwordSha256Hash,omitempty"`
	PasswordSha512Hash string `bson:"passwordSha512Hash,omitempty"`
	PasswordBcryptHash string `bson:"passwordBcryptHash,omitempty"`
}

func (c *Credentials) SetPasswordHash() {
	plainPassword := []byte(c.Password)

	md5Hash := md5.Sum(plainPassword)
	sha1Hash := sha1.Sum(plainPassword)
	sha256Hash := sha256.Sum256(plainPassword)
	sha512Hash := sha512.Sum512(plainPassword)
	bcryptHash, _ := bcrypt.GenerateFromPassword(plainPassword, bcrypt.DefaultCost)

	c.PasswordMD5Hash = hex.EncodeToString(md5Hash[:])
	c.PasswordSha1Hash = hex.EncodeToString(sha1Hash[:])
	c.PasswordSha256Hash = hex.EncodeToString(sha256Hash[:])
	c.PasswordSha512Hash = hex.EncodeToString(sha512Hash[:])
	c.PasswordBcryptHash = string(bcryptHash)
}

func (c *Credentials) CheckEmailLeak() {
	log := context.GetLogger()

	url := fmt.Sprintf("https://haveibeenpwned.com/api/v2/breachedaccount/%s", c.Email)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("Build request object", "Error", err.Error())
		return
	}

	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		log.Error("Make request to haveibeenpwned", "Error", err.Error())
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Info("Verify response code", "Error", fmt.Sprintf("Response code: %d", response.StatusCode))
		return
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("Read response body", "Error", err.Error())
		return
	}

	fmt.Println("Response body: ", string(responseBody))
}
