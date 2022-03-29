package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/auth"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/utils"
	"github.com/jackc/pgx/v4"
	"sync"
	"time"
)

var ErrItemAlreadyExists = errors.New("item already exists")

type StorageDatabase struct {
	conn        *pgx.Conn
	connTimeout time.Duration
}

func NewStorageDatabase(timeout time.Duration, dbDSN string) (*StorageDatabase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		return nil, err
	}
	return &StorageDatabase{conn: conn, connTimeout: timeout}, nil
}

func (s StorageDatabase) GetItem(shortURL string) (storage.ItemRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.connTimeout)
	defer cancel()
	row := s.conn.QueryRow(ctx, "SELECT u.short_url, u.original_url, u.user_id, u.is_deleted FROM urls u WHERE short_url = $1", shortURL)
	item := storage.ItemRepository{}
	err := row.Scan(&item.ShortURL, &item.OriginalURL, &item.User.ID, &item.IsDeleted)
	if err != nil {
		return storage.ItemRepository{}, err
	}

	return item, nil
}

func (s StorageDatabase) SaveItem(originalURL string, user auth.User) (storage.ItemRepository, error) {
	shortURL := utils.GenerateRandomString(15)
	item := storage.ItemRepository{ShortURL: shortURL, OriginalURL: originalURL, User: user}
	ctx, cancel := context.WithTimeout(context.Background(), s.connTimeout)
	defer cancel()
	cmdTag, err := s.conn.Exec(ctx, "INSERT INTO public.urls (original_url, short_url, user_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING", item.OriginalURL, item.ShortURL, item.User.ID)
	if err != nil {
		return storage.ItemRepository{}, err
	}
	if cmdTag.RowsAffected() == 0 {
		item, err = s.getByOriginalURL(originalURL)
		if err != nil {
			return storage.ItemRepository{}, err
		}

		return item, ErrItemAlreadyExists
	}

	return item, nil
}

func (s StorageDatabase) GetItemsByUserID(serverAddress string, user auth.User) ([]storage.ItemRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.connTimeout)
	defer cancel()
	rows, err := s.conn.Query(ctx, "SELECT u.short_url, u.original_url, u.user_id FROM urls u WHERE user_id = $1", user.ID)
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

func (s *StorageDatabase) SaveBatch(ctx context.Context, batch []storage.BatchRequest, user auth.User) ([]storage.ItemRepository, error) {
	var itemRepository storage.ItemRepository
	items := make([]storage.ItemRepository, 0, len(batch))
	for _, item := range batch {
		itemRepository = storage.ItemRepository{OriginalURL: item.OriginalURL, CorrelationID: item.CorrelationID}
		itemRepository.User = user
		itemRepository.ShortURL = utils.GenerateRandomString(16)
		items = append(items, itemRepository)
	}
	if err := insertItems(ctx, s.conn, items); err != nil {
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

func (s StorageDatabase) getByOriginalURL(originalURL string) (storage.ItemRepository, error) {
	row := s.conn.QueryRow(context.Background(), "SELECT u.short_url, u.original_url, u.user_id FROM urls u WHERE original_url = $1", originalURL)
	item := storage.ItemRepository{}
	err := row.Scan(&item.ShortURL, &item.OriginalURL, &item.User.ID)
	if err != nil {
		return storage.ItemRepository{}, err
	}

	return item, nil
}

func (s StorageDatabase) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.connTimeout)
	defer cancel()

	return s.conn.Ping(ctx)
}

func (s StorageDatabase) Close() error {
	if s.conn.IsClosed() {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.connTimeout)
	defer cancel()

	return s.conn.Close(ctx)
}

func (s StorageDatabase) DeleteByIds(ctx context.Context, batchItems []string, user auth.User) error {
	workerChs := make([]chan string, 0, len(batchItems))
	for _, itemID := range batchItems {
		w := s.createSQLForDeleteItem(itemID)
		workerChs = append(workerChs, w)
	}

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)
	batch := <-s.createBatchForDeleteItems(user, workerChs...)

	batchResults := tx.SendBatch(ctx, batch)
	_, err = batchResults.Exec()
	if err != nil {
		return err
	}
	batchResults.Close()

	return tx.Commit(ctx)
}

func (s StorageDatabase) createSQLForDeleteItem(itemID string) chan string {
	sqlRow := make(chan string)

	go func() {
		sqlStatement := fmt.Sprintf("UPDATE urls SET is_deleted = true WHERE user_id = $1 AND short_url = '%s'", itemID)
		sqlRow <- sqlStatement

		close(sqlRow)
	}()

	return sqlRow
}

func (s StorageDatabase) createBatchForDeleteItems(user auth.User, rows ...chan string) chan *pgx.Batch {
	batch := &pgx.Batch{}
	sqlBatch := make(chan *pgx.Batch)

	go func() {
		wg := &sync.WaitGroup{}
		for _, row := range rows {
			wg.Add(1)

			go func(inputCh chan string) {
				defer wg.Done()
				for item := range inputCh {
					batch.Queue(item, user.ID)
				}
			}(row)
		}

		wg.Wait()
		sqlBatch <- batch

		close(sqlBatch)
	}()

	return sqlBatch
}
