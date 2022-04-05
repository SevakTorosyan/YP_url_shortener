package worker

import (
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/auth"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage"
	"github.com/labstack/gommon/log"
)

type ItemsDeleter struct {
	itemsID []string
	user    auth.User
}

func NewItemsDeleter(itemsID []string, user auth.User) ItemsDeleter {
	return ItemsDeleter{itemsID: itemsID, user: user}
}

func DeletingWorker(deleteWorkerCh chan ItemsDeleter, storage storage.Storage) {
	for itemDeleter := range deleteWorkerCh {
		if err := storage.DeleteByIds(itemDeleter.itemsID, itemDeleter.user); err != nil {
			log.Error(err)
		}
	}
}
