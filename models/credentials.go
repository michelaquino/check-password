package models

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"

	"gitlab.globoi.com/michel.aquino/check-password/context"

	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Email              string `form:"email" bson:"email"`
	Password           string `form:"password"`
	EmailBreached      bool   `bson:"emailBreached"`
	PasswordPwned      bool   `bson:"passwordPwned"`
	PasswordMD5Hash    string `bson:"passwordMD5Hash"`
	PasswordSha1Hash   string `bson:"passwordSha1Hash"`
	PasswordSha256Hash string `bson:"passwordSha256Hash"`
	PasswordSha512Hash string `bson:"passwordSha512Hash"`
	PasswordBcryptHash string `bson:"passwordBcryptHash"`
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

	go c.checkEmailBreached(&waitGroup)
	go c.checkPasswordPwned(&waitGroup)

	waitGroup.Wait()
}

func (c *Credentials) checkEmailBreached(waitGroup *sync.WaitGroup) {
	log := context.GetLogger()
	defer waitGroup.Done()

	url := fmt.Sprintf("https://haveibeenpwned.com/api/v2/breachedaccount/%s", c.Email)
	responseCode, err := makeRequestToHaveibeenpwned(url)
	if err != nil {
		log.Error("Error on make request to haveibeenpwned", "Error", err.Error())
		return
	}

	if responseCode == http.StatusOK {
		c.EmailBreached = true
	}
}

func (c *Credentials) checkPasswordPwned(waitGroup *sync.WaitGroup) {
	log := context.GetLogger()
	defer waitGroup.Done()

	url := fmt.Sprintf("https://haveibeenpwned.com/api/v2/pwnedpassword/%s", c.PasswordSha1Hash)
	responseCode, err := makeRequestToHaveibeenpwned(url)
	if err != nil {
		log.Error("Error on make request to haveibeenpwned", "Error", err.Error())
		return
	}

	if responseCode == http.StatusOK {
		c.PasswordPwned = true
	}
}

func makeRequestToHaveibeenpwned(url string) (int, error) {
	log := context.GetLogger()

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("Build request object", "Error", err.Error())
		return 0, err
	}

	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		log.Error("Make request to haveibeenpwned", "Error", err.Error())
		return 0, err
	}

	defer response.Body.Close()
	return response.StatusCode, nil

	// if response.StatusCode != http.StatusOK {
	// 	log.Info("Verify response code", "Error", fmt.Sprintf("Response code: %d", response.StatusCode))
	// 	return
	// }

	// responseBody, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	log.Error("Read response body", "Error", err.Error())
	// 	return
	// }

	// fmt.Println("Response body: ", string(responseBody))
}
