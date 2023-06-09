package jwk_storage

import (
	"cloud.google.com/go/storage"
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"github.com/a-novel/agora-backend/framework/validation"
	"google.golang.org/api/iterator"
	"io"
	"os"
	"path"
	"sort"
	"strings"
)

type Repository interface {
	// Write creates a new entry.
	Write(ctx context.Context, key ed25519.PrivateKey, name string) (*Model, error)
	// Read a key from the specified name.
	Read(ctx context.Context, name string) (*Model, error)
	// List all entries.
	List(ctx context.Context) ([]*Model, error)
	// Delete the specified entry.
	Delete(ctx context.Context, name string) error
}

type fileSystemRepositoryImpl struct {
	basePath string
	prefix   string
}

func NewFileSystemRepository(basePath, prefix string) Repository {
	return &fileSystemRepositoryImpl{basePath: basePath, prefix: prefix}
}

func (repository *fileSystemRepositoryImpl) getPath(name string) string {
	if strings.HasPrefix(name, repository.prefix+"-") {
		return path.Join(repository.basePath, name)
	}

	return path.Join(repository.basePath, fmt.Sprintf("%s-%s", repository.prefix, name))
}

func (repository *fileSystemRepositoryImpl) removePrefix(name string) string {
	return strings.TrimPrefix(name, repository.prefix+"-")
}

func (repository *fileSystemRepositoryImpl) Write(_ context.Context, key ed25519.PrivateKey, name string) (*Model, error) {
	fileWriter, err := os.Create(repository.getPath(name))
	if err != nil {
		return nil, fmt.Errorf("failed to create file %q: %w", name, err)
	}

	defer fileWriter.Close()

	if err = writeKeyToOutput(fileWriter, key); err != nil {
		return nil, fmt.Errorf("failed to write key to file %q: %w", name, err)
	}

	stat, err := fileWriter.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to read stats of file %q: %w", name, err)
	}

	return &Model{
		Key:  key,
		Date: stat.ModTime(),
		Name: repository.removePrefix(stat.Name()),
	}, nil
}

func (repository *fileSystemRepositoryImpl) Read(_ context.Context, name string) (*Model, error) {
	fileReader, err := os.Open(repository.getPath(name))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, validation.ErrNotFound
		}

		return nil, fmt.Errorf("failed to open file %q: %w", repository.getPath(name), err)
	}

	fileData, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read content of file %q: %w", repository.getPath(name), err)
	}

	key, err := unmarshalPrivateKey(fileData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file %q: %w", repository.getPath(name), err)
	}

	stat, err := fileReader.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to read stats of file %q: %w", name, err)
	}

	return &Model{
		Key:  key,
		Date: stat.ModTime(),
		Name: repository.removePrefix(stat.Name()),
	}, nil
}

func (repository *fileSystemRepositoryImpl) List(ctx context.Context) ([]*Model, error) {
	entries, err := os.ReadDir(repository.basePath)
	if err != nil {
		return nil, err
	}

	var records []*Model

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasPrefix(entry.Name(), repository.prefix) {
			record, err := repository.Read(ctx, entry.Name())
			if err != nil {
				return nil, err
			}
			if record != nil {
				records = append(records, record)
			}
		}
	}

	sort.SliceStable(records, func(i, j int) bool {
		return records[i].Date.After(records[j].Date)
	})

	return records, nil
}

func (repository *fileSystemRepositoryImpl) Delete(_ context.Context, name string) error {
	if err := os.Remove(repository.getPath(name)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return validation.ErrNotFound
		}

		return fmt.Errorf("failed to delete file %q: %w", name, err)
	}

	return nil
}

type googleDatastoreRepositoryImpl struct {
	bucket *storage.BucketHandle
}

func NewGoogleDatastoreRepository(bucket *storage.BucketHandle) Repository {
	return &googleDatastoreRepositoryImpl{bucket: bucket}
}

func (repository *googleDatastoreRepositoryImpl) Write(ctx context.Context, key ed25519.PrivateKey, name string) (*Model, error) {
	fileWriter := repository.bucket.Object(name).NewWriter(ctx)
	defer fileWriter.Close()

	if err := writeKeyToOutput(fileWriter, key); err != nil {
		return nil, fmt.Errorf("failed to write key to file %q: %w", name, err)
	}

	return &Model{
		Key:  key,
		Date: fileWriter.Updated,
		Name: fileWriter.Name,
	}, nil
}

func (repository *googleDatastoreRepositoryImpl) Read(ctx context.Context, name string) (*Model, error) {
	fileReader, err := repository.bucket.Object(name).NewReader(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil, validation.ErrNotFound
		}

		return nil, fmt.Errorf("failed to acquire reader for file %q: %w", name, err)
	}
	defer fileReader.Close()

	data, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", name, err)
	}

	key, err := unmarshalPrivateKey(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file %q: %w", name, err)
	}
	return &Model{
		Key:  key,
		Date: fileReader.Attrs.LastModified,
		Name: name,
	}, nil
}

func (repository *googleDatastoreRepositoryImpl) List(ctx context.Context) ([]*Model, error) {
	entries := repository.bucket.Objects(ctx, nil)

	var records []*Model
	for {
		entry, err := entries.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read next entry: %w", err)
		}

		record, err := repository.Read(ctx, entry.Name)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	sort.SliceStable(records, func(i, j int) bool {
		return records[i].Date.After(records[j].Date)
	})

	return records, nil
}

func (repository *googleDatastoreRepositoryImpl) Delete(ctx context.Context, name string) error {
	if err := repository.bucket.Object(name).Delete(ctx); err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return validation.ErrNotFound
		}

		return fmt.Errorf("failed to acquire reader for file %q: %w", name, err)
	}

	return nil
}
