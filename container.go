package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/google/logger"
	"github.com/pkg/errors"
)

type AzblobContainer struct {
	az azblob.ContainerURL
}

func NewContainer(accountName string, accountKey string, containerName string) AzblobContainer {
	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		logger.Fatal("Invalid credentials with error: " + err.Error())
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName),
	)
	containerURL := azblob.NewContainerURL(*URL, p)
	return AzblobContainer{containerURL}
}

func (c *AzblobContainer) Upload(ctx context.Context, fileName string, blobName string) error {
	// Here's how to upload a blob.
	blobURL := c.az.NewBlockBlobURL(blobName)
	file, err := os.Open(fileName)
	if err != nil {
		return errors.Wrap(err, "upload: open file: "+fileName)
	}
	logger.Infof("Uploading the file with blob name: %s\n", fileName)
	_, err = azblob.UploadFileToBlockBlob(
		ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
			BlockSize:   4 * 1024 * 1024,
			Parallelism: 16,
		},
	)
	return errors.Wrapf(err, "upload file to blob [%s]", blobName)
}

func (c *AzblobContainer) UploadDir(ctx context.Context, dir string, blobPrefix string) error {
	err := filepath.Walk(dir, func(filename string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(dir, filename)
		if err != nil {
			return err
		}
		blobName := blobPrefix + relative
		return c.Upload(ctx, filename, blobName)
	})
	return errors.Wrapf(err, "upload dir %s to blob prefix %s", dir, blobPrefix)
}

func (c *AzblobContainer) Delete(ctx context.Context, blobName string) error {
	blobURL := c.az.NewBlockBlobURL(blobName)
	_, err := blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
	return errors.Wrap(err, "delete "+blobName)
}

func (c *AzblobContainer) ListByPrefix(ctx context.Context, blobPrefix string) ([]azblob.BlobItem, error) {
	marker := (azblob.Marker{})
	items := []azblob.BlobItem{}
	for marker.NotDone() {
		// Get a result segment starting with the blob indicated by the current Marker.
		listBlob, err := c.az.ListBlobsFlatSegment(
			ctx, marker, azblob.ListBlobsSegmentOptions{
				Prefix: blobPrefix,
			},
		)
		if err != nil {
			return nil, errors.Wrapf(err, "list blob by prefix [%s] error", blobPrefix)
		}

		// ListBlobs returns the start of the next segment; you MUST use this to get
		// the next segment (after processing the current result segment).
		marker = listBlob.NextMarker
		items = append(items, listBlob.Segment.BlobItems...)
	}

	return items, nil
}

func (c *AzblobContainer) DeleteDir(ctx context.Context, blobPrefix string) error {
	items, err := c.ListByPrefix(ctx, blobPrefix)
	if err != nil {
		return errors.Wrap(err, "delete dir")
	}
	for _, item := range items {
		err := c.Delete(ctx, item.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
