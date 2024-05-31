package main

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true
	} else {
		return false
	}
}

func GetUserHashPassword(email string) (string, error) {
	key := fmt.Sprintf("Emails:%s", email)
	userJSON, err := redis.Bytes(rh.JSONGet(key, "."))
	if err != nil {
		return "not found", err
	}
	user := new(UserCridentials)
	json.Unmarshal(userJSON, &user)
	return user.Password, err
}
