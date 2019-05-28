package main

import (
	"context"
	"testing"

	_ "github.com/joho/godotenv/autoload"
)

func TestUpload(t *testing.T) {
	c := NewConfig()
	container := NewContainer(
		c.AzureStorageAccount,
		c.AzureStorageAccessKey,
		c.AzureBlobContainer,
	)
	filename := "test/hello.txt"
	blobName := c.AzureBlobPrefix + filename
	ctx := context.Background()
	err := container.Upload(ctx, filename, blobName)
	if err != nil {
		t.Fatal(err)
	}
	err2 := container.Delete(ctx, blobName)
	if err2 != nil {
		t.Fatal(err2)
	}
}

func TestUploadDir(t *testing.T) {
	c := NewConfig()
	container := NewContainer(
		c.AzureStorageAccount,
		c.AzureStorageAccessKey,
		c.AzureBlobContainer,
	)
	dir := "test"
	prefix := c.AzureBlobPrefix
	ctx := context.Background()
	err := container.UploadDir(ctx, dir, prefix)
	if err != nil {
		t.Fatal(err)
	}
	err2 := container.DeleteDir(ctx, prefix)
	if err2 != nil {
		t.Fatal(err2)
	}
}
