package main

import (
	"gopkg.in/gorp.v1"
	"encoding/json"
	"log"
)

type User struct {
	Id   int
	Email string `json:"email"`
	Username string	`json:"username"`
	Password string	`json:"password"`
	Salt string
}

func initUser() *User {
	return &User{}
}

func (user *User) save(DbMap *gorp.DbMap) error {
	// TODO: Hash password with salt
	err := DbMap.Insert(user)

	if err != nil {
		return err
	}

	return nil
}

func (user *User) JSON() (string) {
	json, err := json.Marshal(user)
	if (err != nil) {
		log.Println(err)
	}
	return string(json)
}
