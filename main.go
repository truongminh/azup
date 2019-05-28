package main

import (
	"context"
	"io/ioutil"

	"github.com/google/logger"
)

func main() {
	logger.Init("azup", true, false, ioutil.Discard)
	c := NewConfig()
	container := NewContainer(
		c.AzureStorageAccount,
		c.AzureStorageAccessKey,
		c.AzureBlobContainer,
	)
	dir := c.LocalDir
	prefix := c.AzureBlobPrefix
	ctx := context.Background()
	err := container.UploadDir(ctx, dir, prefix)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Uploaded dir %s to blob prefix %s", dir, prefix)
}
