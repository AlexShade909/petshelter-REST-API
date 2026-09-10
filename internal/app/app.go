package app

import (
	"PetShelter/internal/controller"
	"PetShelter/internal/service"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

type App struct{ server *http.Server }

func New() *App {
	store := service.NewStore()
	handler := controller.NewRouter(
		controller.NewDog(service.NewDog(store)),
		controller.NewClinic(service.NewClinic(store)),
		controller.NewShelter(service.NewShelter(store)),
	)
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = "localhost:8080"
	}
	return &App{server: &http.Server{
		Addr: addr, Handler: handler,
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second,
	}}
}
func (a *App) RunPetshelter(ctx context.Context) error {
	if ctx.Err() != nil {
		return nil
	}
	serverErrors := make(chan error, 1)
	go func() { serverErrors <- a.server.ListenAndServe() }()
	log.Printf("HTTP server starting on %s", a.server.Addr)
	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := a.server.Shutdown(shutdownCtx); err != nil {
			_ = a.server.Close()
			return err
		}
		err := <-serverErrors
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}
