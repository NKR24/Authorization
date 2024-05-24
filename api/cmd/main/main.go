package main

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nitishm/go-rejson/v4"
	"github.com/redis/go-redis/v9"
)

func main() {
	fmt.Println("http://localhost:5000")

	ctx := context.Background()
	rh := rejson.NewReJSONHandler()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	rh.SetGoRedisClientWithContext(ctx, rdb)

	e := echo.New()
	e.Use(middleware.CORS())

	e.Logger.Fatal(e.Start(":5000"))
}
