package database

import (
	"context"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/auth"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/utils"
	"github.com/jackc/pgx/v4"
	"time"
)

type StorageDatabase struct {
	dbDSN string
}

func NewStorageDatabase(dbDSN string) (*StorageDatabase, error) {
	return &StorageDatabase{dbDSN: dbDSN}, nil
}

func (s StorageDatabase) GetItem(shortURL string) (storage.ItemRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, s.dbDSN)
	if err != nil {
		return storage.ItemRepository{}, err
	}
	defer conn.Close(ctx)

	row := conn.QueryRow(context.Background(), "SELECT u.short_url, u.original_url, u.user_id FROM urls u WHERE short_url = $1", shortURL)
	item := storage.ItemRepository{}
	err = row.Scan(&item.ShortURL, &item.OriginalURL, &item.User.ID)
	if err != nil {
		return storage.ItemRepository{}, err
	}

	return item, nil
}

func (s StorageDatabase) SaveItem(originalURL string, user auth.User) (storage.ItemRepository, error) {
	shortURL := utils.GenerateRandomString(15)
	item := storage.ItemRepository{ShortURL: shortURL, OriginalURL: originalURL, User: user}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, s.dbDSN)
	if err != nil {
		return storage.ItemRepository{}, err
	}
	defer conn.Close(ctx)
	_, err = conn.Exec(ctx, "INSERT INTO public.urls (original_url, short_url, user_id) VALUES ($1, $2, $3)", item.OriginalURL, item.ShortURL, item.User.ID)
	if err != nil {
		return storage.ItemRepository{}, err
	}

	return item, nil
}

func (s StorageDatabase) GetItemsByUserID(serverAddress string, user auth.User) ([]storage.ItemRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, s.dbDSN)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	rows, err := conn.Query(context.Background(), "SELECT u.short_url, u.original_url, u.user_id FROM urls u WHERE user_id = $1", user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]storage.ItemRepository, 0)

	for rows.Next() {
		var item storage.ItemRepository
		err = rows.Scan(&item.ShortURL, &item.OriginalURL, &item.User.ID)
		if err != nil {
			return nil, err
		}
		item.ShortURL = serverAddress + item.ShortURL

		items = append(items, item)
	}

	return items, nil
}
