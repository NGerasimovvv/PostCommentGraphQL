package main

import (
	"github.com/NGerasimovvv/GraphQL/internal/config"
	"github.com/NGerasimovvv/GraphQL/internal/storage"
	"github.com/NGerasimovvv/GraphQL/server"
)

func main() {
	cfg := config.LoadConfig()
	storageType := storage.StorageType(cfg)
	defer func() {
		if postgresStorage, ok := storageType.(*storage.PostgresStorage); ok {
			postgresStorage.ClosePostgres()
		}
	}()
	server.InitServer(storageType)
}
