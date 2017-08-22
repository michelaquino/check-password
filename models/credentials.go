package models

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"gitlab.globoi.com/michel.aquino/check-password/context"

	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Email             string `form:"email" bson:"email"`
	Password          string `form:"password" bson:"-"`
	EmailPwned        bool   `bson:"emailPwned"`
	EmailLeakList     []leak `bson:"emailLeakList"`
	PasswordPwned     bool   `bson:"passwordPwned"`
	PasswordProcessed bool   `bson:"passwordProcessed"`

	PasswordMD5Hash          string `bson:"passwordMD5Hash"`
	PasswordMD5HashHacked    bool   `bson:"passwordMD5HashHacked"`
	PasswordSha1Hash         string `bson:"passwordSha1Hash"`
	PasswordSha1HashHacked   bool   `bson:"passwordSha1HashHacked"`
	PasswordSha256Hash       string `bson:"passwordSha256Hash"`
	PasswordSha256HashHacked bool   `bson:"passwordSha256HashHacked"`
	PasswordSha512Hash       string `bson:"passwordSha512Hash"`
	PasswordSha512HashHacked bool   `bson:"passwordSha512HashHacked"`
	PasswordBcryptHash       string `bson:"passwordBcryptHash"`
	PasswordBcryptHashHacked bool   `bson:"passwordBcryptHashHacked"`
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

	c.CheckLeaks()
}

func (c *Credentials) CheckLeaks() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go c.checkEmailPwned(&waitGroup)
	go c.checkPasswordPwned(&waitGroup)

	waitGroup.Wait()
}

func (c *Credentials) checkEmailPwned(waitGroup *sync.WaitGroup) {
	log := context.GetLogger()
	defer waitGroup.Done()

	url := fmt.Sprintf("https://haveibeenpwned.com/api/v2/breachedaccount/%s", c.Email)
	response, err := makeRequestToHaveibeenpwned(url)
	if err != nil {
		log.Error("checkEmailPwned - Error on make request to haveibeenpwned", "Error", err.Error())
		return
	}

	if response.StatusCode != http.StatusOK {
		return
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		log.Error("Read response body", "Error", err.Error())
		return
	}

	leakList := []leak{}
	if response.StatusCode == http.StatusOK {
		err = json.Unmarshal(responseBody, &leakList)
		if err != nil {
			log.Error("Parse response body to object", "Error", err.Error())
			return
		}
	}

	c.EmailPwned = true
	c.EmailLeakList = leakList
	log.Info("Check if email is pwned", "Success", "")
}

func (c *Credentials) checkPasswordPwned(waitGroup *sync.WaitGroup) {
	log := context.GetLogger()
	defer waitGroup.Done()

	url := fmt.Sprintf("https://haveibeenpwned.com/api/v2/pwnedpassword/%s", c.PasswordSha1Hash)
	response, err := makeRequestToHaveibeenpwned(url)
	if err != nil {
		log.Error("checkPasswordPwned - Error on make request to haveibeenpwned", "Error", err.Error())
		return
	}

	if response.StatusCode == http.StatusOK {
		c.PasswordPwned = true
	}

	log.Info("Check if password is pwned", "Success", "")
}

func makeRequestToHaveibeenpwned(url string) (*http.Response, error) {
	log := context.GetLogger()

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("Build request object", "Error", err.Error())
		return nil, err
	}

	httpClient := http.Client{
		Timeout: time.Duration(3 * time.Second),
	}
	response, err := httpClient.Do(request)
	if err != nil {
		log.Error("Make request to haveibeenpwned", "Error", err.Error())
		return nil, err
	}

	return response, nil
}

type leak struct {
	Title string `json:"Title"`
}
