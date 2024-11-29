package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/heltirj/image_previewer/internal/app"
	"github.com/heltirj/image_previewer/internal/cache"
	"github.com/heltirj/image_previewer/internal/config"
	"github.com/heltirj/image_previewer/internal/logger"
	"github.com/heltirj/image_previewer/internal/server/http"
)

const configPath = "./configs/config.yaml"

func main() {
	conf, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load conf: %v", err)
	}

	logg := logger.New(logger.LogLevelInfo)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	a := app.New(logg, cache.NewLruImageCache(conf.LRUSize, conf.StoragePath))

	err = a.Cache.Load()
	if err != nil {
		log.Fatalf("error loading cache: %s", err)
	}

	server := http.NewServer(logg, a, conf.Port)

	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("image_previewer is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
