package database

import (
	"context"
	"fmt"
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

func (s *StorageDatabase) SaveBatch(batch []storage.BatchRequest, user auth.User, ctx context.Context) ([]storage.ItemRepository, error) {
	var itemRepository storage.ItemRepository
	items := make([]storage.ItemRepository, 0, len(batch))
	for _, item := range batch {
		itemRepository = storage.ItemRepository{OriginalURL: item.OriginalURL, CorrelationID: item.CorrelationID}
		itemRepository.User = user
		itemRepository.ShortURL = utils.GenerateRandomString(16)
		items = append(items, itemRepository)
	}
	conn, err := pgx.Connect(ctx, s.dbDSN)
	if err != nil {
		fmt.Println(err.Error())

		return nil, err
	}
	defer conn.Close(ctx)
	if err = insertItems(ctx, conn, items); err != nil {
		fmt.Println(err.Error())

		return nil, err
	}

	return items, nil
}

func insertItems(ctx context.Context, conn *pgx.Conn, items []storage.ItemRepository) error {
	b := &pgx.Batch{}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	for _, item := range items {
		sqlStatement := `INSERT INTO public.urls (original_url, short_url, user_id, correlation_id) VALUES ($1, $2, $3, $4)`
		b.Queue(sqlStatement, item.OriginalURL, item.ShortURL, item.User.ID, item.CorrelationID)
	}

	batchResults := tx.SendBatch(ctx, b)
	_, err = batchResults.Exec()
	if err != nil {
		return err
	}
	batchResults.Close()

	return tx.Commit(ctx)
}
