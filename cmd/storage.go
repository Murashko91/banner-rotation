package main

import (
	"github.com/otus-murashko/banners-rotation/internal/config"
	"github.com/otus-murashko/banners-rotation/internal/storage"
	"github.com/otus-murashko/banners-rotation/internal/storage/psql"
)

func getStorage(config config.DBConfig) storage.Storage {
	/*	if config.InMemory {
		return memorystorage.New()
	}*/

	return psql.New(psql.StorageInfo{
		Host:     config.Host,
		Port:     config.Port,
		User:     config.User,
		Password: config.Password,
		DBName:   config.DBName,
	})
}
