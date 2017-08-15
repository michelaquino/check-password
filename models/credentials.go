package models

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Email              string `form:"email"`
	Password           string `form:"password"`
	PasswordMD5Hash    string
	PasswordSha1Hash   string
	PasswordSha256Hash string
	PasswordSha512Hash string
	PasswordBcryptHash string
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
