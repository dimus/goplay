package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Run starts HTTP/1 service for scientific names verification.
func Run() {
	log.Printf("Starting the HTTP API server on port %d.", 8888)
	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/ping", wait())
	e.POST("/ping", waitPost())

	addr := fmt.Sprintf(":%d", 8888)
	s := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
}

func main() {
	Run()
}

func wait() func(echo.Context) error {
	return func(c echo.Context) error {
		chErr := make(chan error)
		fmt.Println("start")
		ctx := c.Request().Context()
		// 2 -- cancel via remote or local timeout.
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		go func() {
			time.Sleep(100 * time.Second)
			chErr <- c.String(http.StatusOK, "pong")
		}()
		select {
		case <-ctx.Done():
			fmt.Println("context end")
			return ctx.Err()
		case err := <-chErr:
			fmt.Println("end")
			return err
		}
	}
}

func waitPost() func(echo.Context) error {
	return func(c echo.Context) error {
		return nil
	}
}
