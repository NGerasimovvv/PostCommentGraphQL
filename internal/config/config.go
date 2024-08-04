package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Postgres *PostgresConfig
	Storage  *StorageTypeConfig
}

type StorageTypeConfig struct {
	StorageType string
}

type PostgresConfig struct {
	PostgresPort     string
	PostgresHost     string
	DatabaseName     string
	PostgresUser     string
	PostgresPassword string
}

func LoadConfig() *Config {
	postgresConfig := loadPostgresConfig()
	storageTypeConfig := loadStorageTypeConfig()
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	logger.Printf("Postgres Port: %s", postgresConfig.PostgresPort)
	logger.Printf("StorageType: %s", storageTypeConfig.StorageType)
	return &Config{
		Postgres: postgresConfig,
		Storage:  storageTypeConfig,
	}
}

func loadPostgresConfig() *PostgresConfig {
	const opt = "LoadPostgresConfig"
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("%s: %v", opt, err)
	}
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	port, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		logger.Fatal("POSTGRES! Can't read PORT")
	}

	host, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		logger.Fatal("POSTGRES! Can't read HOST")
	}

	password, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		logger.Fatal("POSTGRES! Can't read PASSWORD")
	}

	db, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		logger.Fatal("POSTGRES! Can't read DB")
	}

	user, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		logger.Fatal("POSTGRES! Can't read USER")
	}

	return &PostgresConfig{
		PostgresPort:     port,
		PostgresHost:     host,
		DatabaseName:     db,
		PostgresUser:     user,
		PostgresPassword: password,
	}
}

func loadStorageTypeConfig() *StorageTypeConfig {
	err := godotenv.Load(".env")
	const opt = "loadStorageConfig"
	if err != nil {
		log.Fatalf("%s: %v", opt, err)
	}
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	storageType, ok := os.LookupEnv("STORAGE_TYPE")
	if !ok {
		logger.Fatal("Can't read STORAGE_TYPE")
	}
	return &StorageTypeConfig{StorageType: storageType}
}
