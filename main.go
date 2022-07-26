package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/jaredwarren/plex_ytdl/config"
	"github.com/jaredwarren/plex_ytdl/server"
	"github.com/jaredwarren/rpi_music/log"
	"github.com/spf13/viper"
)

const (
	DBPath = "my.db"
)

func main() {
	logger := log.NewStdLogger(log.Debug)

	// Init Config
	config.InitConfig(logger)
	logger.SetLevel(log.Level(viper.GetInt64("log.level")))

	// Init Server
	htmlServer := server.StartHTTPServer(&server.Config{
		Host:         viper.GetString("host"),
		ReadTimeout:  35 * time.Second,
		WriteTimeout: 35 * time.Second,
		Logger:       logger,
	})
	defer htmlServer.StopHTTPServer()

	// Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	logger.Info("main :shutting down")
}
