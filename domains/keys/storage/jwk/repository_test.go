package jwk_storage

import (
	"context"
	"crypto/ed25519"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFileRepository_Write(t *testing.T) {
	data := []struct {
		name string

		key     ed25519.PrivateKey
		keyName string

		expect    *Model
		expectErr error
	}{
		{
			name:    "Success",
			key:     MockedKeys[4],
			keyName: "test-4",
			expect: &Model{
				Key:  MockedKeys[4],
				Name: "test-4",
			},
		},
		{
			name:    "Success/Exists",
			key:     MockedKeys[4],
			keyName: "test-2",
			expect: &Model{
				Key:  MockedKeys[4],
				Name: "test-2",
			},
		},
	}

	err := test_utils.RunFileTransactionalTest(t, Fixtures, func(ctx context.Context, basePath string) {
		repository := NewFileSystemRepository(basePath, "foo")

		for _, d := range data {
			t.Run(d.name, func(t *testing.T) {
				res, err := repository.Write(ctx, d.key, d.keyName)
				test_utils.RequireError(t, d.expectErr, err)

				if d.expect != nil {
					require.Equal(t, d.expect.Name, res.Name)
					require.True(t, d.expect.Key.Equal(res.Key))
				} else {
					require.Nil(t, res)
				}
			})
		}
	})
	require.NoError(t, err)
}

func TestFileRepository_Read(t *testing.T) {
	data := []struct {
		name string

		keyName string

		expect    *Model
		expectErr error
	}{
		{
			name:    "Success",
			keyName: "test-2",
			expect: &Model{
				Key:  MockedKeys[1],
				Name: "test-2",
			},
		},
		{
			name:      "Error/NotExists",
			keyName:   "test-4",
			expectErr: validation.ErrNotFound,
		},
	}

	err := test_utils.RunFileTransactionalTest(t, Fixtures, func(ctx context.Context, basePath string) {
		repository := NewFileSystemRepository(basePath, "foo")

		for _, d := range data {
			t.Run(d.name, func(t *testing.T) {
				res, err := repository.Read(ctx, d.keyName)
				test_utils.RequireError(t, d.expectErr, err)

				if d.expect != nil {
					require.Equal(t, d.expect.Name, res.Name)
					require.True(t, d.expect.Key.Equal(res.Key))
				} else {
					require.Nil(t, res)
				}
			})
		}
	})
	require.NoError(t, err)
}

func TestFileRepository_List(t *testing.T) {
	data := []struct {
		name string

		prefix string

		expect    []*Model
		expectErr error
	}{
		{
			name:   "Success",
			prefix: "foo",
			expect: []*Model{
				{
					Key:  MockedKeys[1],
					Name: "test-2",
				},
				{
					Key:  MockedKeys[2],
					Name: "test-3",
				},
				{
					Key:  MockedKeys[0],
					Name: "test-1",
				},
			},
		},
		{
			name:   "Success/NoKeys",
			prefix: "qux",
			expect: []*Model(nil),
		},
	}

	err := test_utils.RunFileTransactionalTest(t, Fixtures, func(ctx context.Context, basePath string) {
		for _, d := range data {
			repository := NewFileSystemRepository(basePath, d.prefix)
			t.Run(d.name, func(t *testing.T) {
				res, err := repository.List(ctx)
				test_utils.RequireError(t, d.expectErr, err)

				if d.expect != nil {
					require.Len(t, res, len(d.expect))
					for i, e := range d.expect {
						require.Equal(t, e.Name, res[i].Name)
						require.True(t, e.Key.Equal(res[i].Key))
					}
				} else {
					require.Nil(t, res)
				}
			})
		}
	})
	require.NoError(t, err)
}

func TestFileRepository_Delete(t *testing.T) {
	data := []struct {
		name string

		keyName string

		expectErr error
	}{
		{
			name:    "Success",
			keyName: "test-2",
		},
		{
			name:      "Error/NotExists",
			keyName:   "test-4",
			expectErr: validation.ErrNotFound,
		},
	}

	err := test_utils.RunFileTransactionalTest(t, Fixtures, func(ctx context.Context, basePath string) {
		repository := NewFileSystemRepository(basePath, "foo")

		for _, d := range data {
			t.Run(d.name, func(t *testing.T) {
				err := repository.Delete(ctx, d.keyName)
				test_utils.RequireError(t, d.expectErr, err)
			})
		}
	})
	require.NoError(t, err)
}
