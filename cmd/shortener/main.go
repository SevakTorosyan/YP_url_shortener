package main

import (
	"context"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/config"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/handlers"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage/database"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage/file"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage/memory"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/worker"
	"github.com/jackc/pgx/v4"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	cfg := config.InitConfig()

	s := initStorage(cfg)
	defer s.Close()
	deletingWorkerCh := make(chan worker.ItemsDeleter)
	handler := handlers.NewHandler(s, cfg, deletingWorkerCh)
	go worker.DeletingWorker(deletingWorkerCh, s)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, handler))
}

func initDatabase(databaseURL string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer cancel()
	defer conn.Close(ctx)

	sqlFile, err := ioutil.ReadFile("db/init.sql")
	if err != nil {
		log.Fatal(err)
	}

	requests := strings.Split(string(sqlFile), ";")

	for _, request := range requests {
		_, err := conn.Exec(ctx, request)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func initStorage(cfg *config.Config) storage.Storage {
	if cfg.DatabaseDSN != "" {
		initDatabase(cfg.DatabaseDSN)
		stg, err := database.NewStorageDatabase(time.Duration(cfg.DatabaseTimeout)*time.Second, cfg.DatabaseDSN)
		if err != nil {
			log.Fatal(err)
		}

		return stg
	}
	if cfg.FileStoragePath != "" {
		stg, err := file.NewStorageFile(cfg.FileStoragePath)
		if err != nil {
			log.Fatal(err)
		}

		return stg
	}

	return memory.NewStorageMap()
}
