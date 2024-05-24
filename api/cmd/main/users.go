package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/google/uuid"

	"github.com/redis/go-redis/v9"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

func register(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return err
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	user.Password = hashedPassword

	user.ID = uuid.New()

	userIdKey := fmt.Sprintf("Users:%s", user.ID)
	_, err = rh.JSONSet(userIdKey, ".", user)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	tempPasswordHashKey := fmt.Sprintf("TempPasswordHash:%s", user.Username)
	_, err = rh.JSONSet(tempPasswordHashKey, ".", hashedPassword)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	userNameKey := fmt.Sprintf("UsersByName:%s", user.Username)
	existingUserInterface, err := rh.JSONGet(userNameKey, ".")
	if err != nil {
		if errors.Is(err, redis.Nil) {
			_, err = rh.JSONSet(userNameKey, ".", user)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}
			_, err = rh.JSONDel(tempPasswordHashKey, ".")
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}
			return c.JSON(http.StatusCreated, user.ID)
		} else {
			return c.NoContent(http.StatusInternalServerError)
		}
	} else {
		existingUserBytes, ok := existingUserInterface.([]byte)
		if !ok {
			return c.NoContent(http.StatusInternalServerError)
		}
		var existingUser User
		if err := json.Unmarshal(existingUserBytes, &existingUser); err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusConflict, nil)
	}
}
