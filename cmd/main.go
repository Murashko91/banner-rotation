package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/otus-murashko/banners-rotation/internal/app"
	"github.com/otus-murashko/banners-rotation/internal/config"
	internalhttp "github.com/otus-murashko/banners-rotation/internal/server/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "conf", "./../configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := config.GetBannersConfig(configFile)
	storage := getStorage(config.Database)
	if err := storage.Connect(); err != nil {
		log.Println(err.Error())
	}
	bannerApp := app.New(storage)
	server := internalhttp.NewServer(bannerApp, config.Server)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Printf("failed to stop http server: %s \n", err.Error())
		}
	}()

	log.Println("banner server is running...")
	if err := server.Start(ctx); err != nil {
		log.Printf("failed to start http server: %s \n", err.Error())
	}
}
