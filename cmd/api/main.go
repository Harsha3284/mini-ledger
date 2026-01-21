package main

import (
	"context"
	"log"

	"mini-ledger/internal/config"
	"mini-ledger/internal/db"
	apphttp "mini-ledger/internal/http"
	"mini-ledger/internal/http/handlers"
	"mini-ledger/internal/service"
	"mini-ledger/internal/store"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DBURL)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer pool.Close()

	r := apphttp.NewRouter(pool)

	v1 := r.Group("/v1")

	// accounts feature (store -> service -> handler)
	accountStore := store.NewAccountStore(pool)
	accountService := service.NewAccountService(accountStore)
	accountHandler := handlers.NewAccountHandler(accountService)

	ledgerStore := store.NewLedgerStore(pool)
	ledgerService := service.NewLedgerService(ledgerStore, accountStore)
	ledgerHandler := handlers.NewLedgerHandler(ledgerService)

	accountHandler.Register(v1)
	ledgerHandler.Register(v1)

	log.Printf("mini-ledger listening on :%s", cfg.HTTPPort)
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal(err)
	}
}
