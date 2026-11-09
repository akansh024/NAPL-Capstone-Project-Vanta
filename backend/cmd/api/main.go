package main

import (
	"context"
	"errors"
	"github.com/vanta/vanta/backend/internal/api"
	"github.com/vanta/vanta/backend/internal/config"
	"github.com/vanta/vanta/backend/internal/database"
	"github.com/vanta/vanta/backend/internal/findings"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	c := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := database.Open(ctx, c.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	if err = db.Ping(ctx); err != nil {
		return err
	}
	if err = database.Migrate(ctx, db); err != nil {
		return err
	}
	a := api.API{DB: db, Store: findings.Store{DB: db}}
	server := &http.Server{Addr: ":" + c.Port, Handler: a.Router(c.FrontendOrigin)}
	errCh := make(chan error, 1)
	go func() { errCh <- server.ListenAndServe() }()

	log.Printf("VANTA API listening on :%s", c.Port)
	shutdown, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err = <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-shutdown.Done():
		log.Print("shutting down VANTA API")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err = server.Shutdown(shutdownCtx); err != nil {
			return err
		}
	}
	return nil
}
