package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/otus-murashko/banners-rotation/internal/app"
	"github.com/otus-murashko/banners-rotation/internal/config"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Server struct {
	server *http.Server
	app    app.Application
}

type ServerConf struct {
	Host string
	Port int
}

func NewServer(app app.Application, conf config.Server) *Server {
	bannerRouter := http.NewServeMux()

	appHandler := Handler{
		app: app,
	}

	swagHandler := httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s:%d/swagger/doc.json", conf.Host, conf.Port)))

	bannerRouter.Handle("/banner-rotation", loggingMiddleware(http.HandlerFunc(appHandler.bannerRotationHandler)))
	bannerRouter.Handle("/banner", loggingMiddleware(http.HandlerFunc(appHandler.bannerHandler)))
	bannerRouter.Handle("/slot", loggingMiddleware(http.HandlerFunc(appHandler.slotHandler)))
	bannerRouter.Handle("/group", loggingMiddleware(http.HandlerFunc(appHandler.groupHandler)))
	bannerRouter.Handle("/stat", loggingMiddleware(http.HandlerFunc(appHandler.statHandler)))
	bannerRouter.Handle("/swagger/", swagHandler)

	httpServer := &http.Server{
		ReadHeaderTimeout: 3 * time.Second,
		Addr:              fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Handler:           bannerRouter,
	}
	return &Server{
		server: httpServer,
		app:    app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.server.ListenAndServe()
	ctx.Done()
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	<-ctx.Done()
	return err
}
