package main

import (
	"gopkg.in/gorp.v1"
	"encoding/json"
	"log"
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha1"
	"encoding/base64"
	"crypto/rand"
	"io"
)

type User struct {
	Id   int
	Email string `json:"email"`
	Username string	`json:"username"`
	Password string `json:"password"`
	Salt []byte `json:"-"`
}

func initUser() *User {
	return &User{}
}

func (user *User) save(DbMap *gorp.DbMap) error {

	// Generate random salt
	salt := make([]byte, 10)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return err
	}

	user.Salt = salt

	hash := pbkdf2.Key([]byte(user.Password), user.Salt, 4096, 32, sha1.New)
	user.Password = base64.StdEncoding.EncodeToString(hash)

	err = DbMap.Insert(user)

	if err != nil {
		return err
	}

	return nil
}

func (user *User) checkPassword(password string) bool {
	hash := pbkdf2.Key([]byte(password), user.Salt, 4096, 32, sha1.New)
	base64 := base64.StdEncoding.EncodeToString(hash)

	if (base64 == user.Password) {
		return true
	} else {
		return false
	}
}

func (user *User) JSON() (string) {
	json, err := json.Marshal(user)
	if (err != nil) {
		log.Println(err)
	}
	return string(json)
}
