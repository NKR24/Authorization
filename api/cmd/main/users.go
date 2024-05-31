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
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

type UserCridentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	tempPasswordHashKey := fmt.Sprintf("TempPasswordHash:%s", user.Email)
	_, err = rh.JSONSet(tempPasswordHashKey, ".", hashedPassword)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	userCridentials := UserCridentials{
		Email:    user.Email,
		Password: user.Password,
	}

	userNameKey := fmt.Sprintf("Cridentials:%s", user.Email)
	existingUserInterface, err := rh.JSONGet(userNameKey, ".")
	if err != nil {
		if errors.Is(err, redis.Nil) {
			_, err = rh.JSONSet(userNameKey, ".", userCridentials)
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
		var existingUser UserCridentials
		if err := json.Unmarshal(existingUserBytes, &existingUser); err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusConflict, "This user is already registered")
	}
}

func login(c echo.Context) error {
	loginData := new(UserCridentials)
	if err := c.Bind(loginData); err != nil {
		return err
	}
	userKey := fmt.Sprintf("Cridentials:%s", loginData.Email)
	userInterface, err := rh.JSONGet(userKey, ".")
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return c.JSON(http.StatusUnauthorized, "Invalid email or password")
		} else {
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	userBytes, ok := userInterface.([]byte)
	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}
	var user User
	if err := json.Unmarshal(userBytes, &user); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	fmt.Println(CheckPasswordHash(loginData.Password, user.Password))

	if !CheckPasswordHash(loginData.Password, user.Password) {
		return c.JSON(http.StatusInternalServerError, "Password is wrong")
	}

	sessionToken := uuid.New().String()
	sessionTokenKey := fmt.Sprintf("SessionTokens:%s", sessionToken)
	_, err = rh.JSONSet(sessionTokenKey, ".", user.ID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	c.SetCookie(&http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		HttpOnly: true,
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "Logged in successfully"})
}
