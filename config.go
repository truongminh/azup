package main

import (
	"os"
	"strings"

	"github.com/google/logger"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	AzureStorageAccount   string
	AzureStorageAccessKey string
	AzureBlobContainer    string
	AzureBlobPrefix       string
	LocalDir              string
}

func NewConfig() Config {
	// From the Azure portal, get your storage account name and key and set environment variables.
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT")
	if accountName == "" {
		logger.Fatal("Missing env AZURE_STORAGE_ACCOUNT")
	}
	accountKey := os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	if accountKey == "" {
		logger.Fatal("Missing env AZURE_STORAGE_ACCESS_KEY")
	}
	containerName := os.Getenv("AZURE_BLOB_CONTAINER")
	if containerName == "" {
		logger.Fatal("Missing env AZURE_BLOB_CONTAINER")
	}
	blobPrefix := os.Getenv("AZURE_BLOB_PREFIX")
	if blobPrefix == "" {
		logger.Fatal("Missing env AZURE_BLOB_PREFIX")
	}
	if !strings.HasSuffix(blobPrefix, "/") {
		blobPrefix += "/"
	}
	if len(os.Args) < 2 {
		logger.Fatal("Usage: azup <dir>")
	}
	localDir := os.Args[1]
	c := Config{
		AzureStorageAccount:   accountName,
		AzureStorageAccessKey: accountKey,
		AzureBlobContainer:    containerName,
		AzureBlobPrefix:       blobPrefix,
		LocalDir:              localDir,
	}
	return c
}
