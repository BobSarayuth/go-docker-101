package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func main() {

	e := echo.New()

	e.GET("/api/ping", ping)

	go startServer(e)

	waitForGracefulShutdown(e)
}

func ping(c echo.Context) error {
	ping := os.Getenv("PING")
	if len(ping) == 0 {
		return c.String(http.StatusOK, "OK")
	} else {
		return c.String(http.StatusOK, ping)
	}
}

func startServer(e *echo.Echo) {
	if err := e.Start(":8080"); err != nil {
		log.Info("shutting down the server")
	}
}

func waitForGracefulShutdown(e *echo.Echo) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}