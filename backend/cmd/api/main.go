package main

import (
	"context"
	"github.com/vanta/vanta/backend/internal/api"
	"github.com/vanta/vanta/backend/internal/config"
	"github.com/vanta/vanta/backend/internal/database"
	"github.com/vanta/vanta/backend/internal/findings"
	"log"
	"net/http"
	"time"
)

func main() {
	c := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, e := database.Open(ctx, c.DatabaseURL)
	if e != nil {
		log.Fatal(e)
	}
	defer db.Close()
	if e = database.Migrate(ctx, db); e != nil {
		log.Fatal(e)
	}
	a := api.API{DB: db, Store: findings.Store{DB: db}}
	log.Printf("VANTA API listening on :%s", c.Port)
	log.Fatal(http.ListenAndServe(":"+c.Port, a.Router(c.FrontendOrigin)))
}
