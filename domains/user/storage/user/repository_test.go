package user_storage

import (
	"context"
	"fmt"
	"github.com/a-novel/agora-backend/domains/user/storage/credentials"
	"github.com/a-novel/agora-backend/domains/user/storage/identity"
	"github.com/a-novel/agora-backend/domains/user/storage/profile"
	"github.com/a-novel/agora-backend/framework/test"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"sort"
	"testing"
	"time"
)

var (
	baseTime   = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
	updateTime = time.Date(2020, time.May, 4, 9, 0, 0, 0, time.UTC)
)

func concatFixtures(s ...[]interface{}) []interface{} {
	var fixtures []interface{}
	for _, fixture := range s {
		fixtures = append(fixtures, fixture...)
	}

	return fixtures
}

// generateSearchFixture allows to simplify fixture creation by getting rid of some parameters.
func generateSearchFixture(
	username, slug,
	firstName, lastName string,
	sex models.Sex,
	createdAt time.Time, index int,
) []interface{} {
	return []interface{}{
		&credentials_storage.Model{
			ID:        test_utils.NumberUUID(index),
			CreatedAt: createdAt,
			UpdatedAt: &createdAt,
			Core: credentials_storage.Core{
				Email: models.Email{
					User:   fmt.Sprintf("user+%d", index),
					Domain: "gmail.com",
				},
				Password: models.Password{
					Hashed: "foobarqux",
				},
			},
		},
		&identity_storage.Model{
			ID:        test_utils.NumberUUID(index),
			CreatedAt: createdAt,
			UpdatedAt: &createdAt,
			Core: identity_storage.Core{
				FirstName: firstName,
				LastName:  lastName,
				Birthday:  time.Date(2000, time.June, 28, 12, 0, 0, 0, time.UTC),
				Sex:       sex,
			},
		},
		&profile_storage.Model{
			ID:        test_utils.NumberUUID(index),
			CreatedAt: createdAt,
			UpdatedAt: &createdAt,
			Core: profile_storage.Core{
				Username: username,
				Slug:     slug,
			},
		},
	}
}

var Fixtures = []interface{}{
	// Standard user.
	&credentials_storage.Model{
		ID:        test_utils.NumberUUID(1000),
		CreatedAt: baseTime,
		UpdatedAt: &baseTime,
		Core: credentials_storage.Core{
			Email: models.Email{
				User:   "anna.banana",
				Domain: "coco.nut",
			},
			Password: models.Password{
				Hashed: "foobarqux",
			},
		},
	},
	&identity_storage.Model{
		ID:        test_utils.NumberUUID(1000),
		CreatedAt: baseTime,
		UpdatedAt: &baseTime,
		Core: identity_storage.Core{
			FirstName: "Anna",
			LastName:  "Banana",
			Birthday:  time.Date(2000, time.June, 28, 12, 0, 0, 0, time.UTC),
			Sex:       models.SexFemale,
		},
	},
	&profile_storage.Model{
		ID:        test_utils.NumberUUID(1000),
		CreatedAt: baseTime,
		UpdatedAt: &baseTime,
		Core: profile_storage.Core{
			Username: "BanAnna",
			Slug:     "fruit-basket",
		},
	},
}

func TestUserRepository_Create(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		data *Core
		id   uuid.UUID
		now  time.Time

		expect    *Model
		expectErr error
	}{
		{
			name: "Success",
			data: &Core{
				Credentials: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
				Identity: identity_storage.Core{
					FirstName: "Elon",
					LastName:  "Bezos",
					Birthday:  time.Date(1964, time.June, 28, 19, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: profile_storage.Core{
					Slug: "square-the-circle",
				},
			},
			id:  test_utils.NumberUUID(1),
			now: baseTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Core: Core{
					Credentials: credentials_storage.Core{
						Email: models.Email{
							User:   "elon.bezos",
							Domain: "gmail.com",
						},
						Password: models.Password{Hashed: "foobarqux"},
					},
					Identity: identity_storage.Core{
						FirstName: "Elon",
						LastName:  "Bezos",
						Birthday:  time.Date(1964, time.June, 28, 19, 0, 0, 0, time.UTC),
						Sex:       models.SexMale,
					},
					Profile: profile_storage.Core{
						Slug: "square-the-circle",
					},
				},
			},
		},
		{
			name: "Success/WithEmailValidation",
			data: &Core{
				Credentials: credentials_storage.Core{
					Email: models.Email{
						Validation: "youshallpass",
						User:       "elon.bezos",
						Domain:     "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
				Identity: identity_storage.Core{
					FirstName: "Elon",
					LastName:  "Bezos",
					Birthday:  time.Date(1964, time.June, 28, 19, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: profile_storage.Core{
					Slug: "square-the-circle",
				},
			},
			id:  test_utils.NumberUUID(1),
			now: baseTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Core: Core{
					Credentials: credentials_storage.Core{
						Email: models.Email{
							Validation: "youshallpass",
							User:       "elon.bezos",
							Domain:     "gmail.com",
						},
						Password: models.Password{Hashed: "foobarqux"},
					},
					Identity: identity_storage.Core{
						FirstName: "Elon",
						LastName:  "Bezos",
						Birthday:  time.Date(1964, time.June, 28, 19, 0, 0, 0, time.UTC),
						Sex:       models.SexMale,
					},
					Profile: profile_storage.Core{
						Slug: "square-the-circle",
					},
				},
			},
		},
		{
			name: "Success/WithPasswordReset",
			data: &Core{
				Credentials: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					Password: models.Password{Validation: "youshallpass"},
				},
				Identity: identity_storage.Core{
					FirstName: "Elon",
					LastName:  "Bezos",
					Birthday:  time.Date(1964, time.June, 28, 19, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: profile_storage.Core{
					Slug: "square-the-circle",
				},
			},
			id:  test_utils.NumberUUID(1),
			now: baseTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Core: Core{
					Credentials: credentials_storage.Core{
						Email: models.Email{
							User:   "elon.bezos",
							Domain: "gmail.com",
						},
						Password: models.Password{Validation: "youshallpass"},
					},
					Identity: identity_storage.Core{
						FirstName: "Elon",
						LastName:  "Bezos",
						Birthday:  time.Date(1964, time.June, 28, 19, 0, 0, 0, time.UTC),
						Sex:       models.SexMale,
					},
					Profile: profile_storage.Core{
						Slug: "square-the-circle",
					},
				},
			},
		},
		{
			name: "Success/WithUsername",
			data: &Core{
				Credentials: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
				Identity: identity_storage.Core{
					FirstName: "Elon",
					LastName:  "Bezos",
					Birthday:  time.Date(1964, time.June, 28, 19, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: profile_storage.Core{
					Username: "SpaceCatGirl",
					Slug:     "square-the-circle",
				},
			},
			id:  test_utils.NumberUUID(1),
			now: baseTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Core: Core{
					Credentials: credentials_storage.Core{
						Email: models.Email{
							User:   "elon.bezos",
							Domain: "gmail.com",
						},
						Password: models.Password{Hashed: "foobarqux"},
					},
					Identity: identity_storage.Core{
						FirstName: "Elon",
						LastName:  "Bezos",
						Birthday:  time.Date(1964, time.June, 28, 19, 0, 0, 0, time.UTC),
						Sex:       models.SexMale,
					},
					Profile: profile_storage.Core{
						Username: "SpaceCatGirl",
						Slug:     "square-the-circle",
					},
				},
			},
		},
		{
			name: "Error/EmailTaken",
			data: &Core{
				Credentials: credentials_storage.Core{
					Email: models.Email{
						User:   "anna.banana",
						Domain: "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
				Identity: identity_storage.Core{
					FirstName: "Elon",
					LastName:  "Bezos",
					Birthday:  time.Date(1964, time.June, 28, 19, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: profile_storage.Core{
					Slug: "square-the-circle",
				},
			},
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrUniqConstraintViolation,
		},
		{
			name: "Error/EmailUserMissing",
			data: &Core{
				Credentials: credentials_storage.Core{
					Email: models.Email{
						Domain: "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
				Identity: identity_storage.Core{
					FirstName: "Elon",
					LastName:  "Bezos",
					Birthday:  time.Date(1964, time.June, 28, 19, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: profile_storage.Core{
					Slug: "square-the-circle",
				},
			},
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrConstraintViolation,
		},
		{
			name: "Error/EmailDomainMissing",
			data: &Core{
				Credentials: credentials_storage.Core{
					Email: models.Email{
						User: "elon.bezos",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
				Identity: identity_storage.Core{
					FirstName: "Elon",
					LastName:  "Bezos",
					Birthday:  time.Date(1964, time.June, 28, 19, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: profile_storage.Core{
					Slug: "square-the-circle",
				},
			},
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrConstraintViolation,
		},
		{
			name: "Error/SlugTaken",
			data: &Core{
				Credentials: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
				Identity: identity_storage.Core{
					FirstName: "Elon",
					LastName:  "Bezos",
					Birthday:  time.Date(1964, time.June, 28, 19, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: profile_storage.Core{
					Slug: "fruit-basket",
				},
			},
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrUniqConstraintViolation,
		},
		{
			name: "Error/SlugMissing",
			data: &Core{
				Credentials: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
				Identity: identity_storage.Core{
					FirstName: "Elon",
					LastName:  "Bezos",
					Birthday:  time.Date(1964, time.June, 28, 19, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: profile_storage.Core{},
			},
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrConstraintViolation,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.BeginTx(ctx, nil)
				require.NoError(st, err)
				defer stx.Rollback()

				res, err := NewRepository(stx).Create(ctx, d.data, d.id, d.now)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestUserRepository_Delete(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		id  uuid.UUID
		now time.Time

		expect    *Model
		expectErr error
	}{
		{
			name: "Success",
			id:   test_utils.NumberUUID(1000),
			now:  updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &baseTime,
				DeletedAt: &updateTime,
				Core: Core{Identity: identity_storage.Core{
					FirstName: "",
					LastName:  "",
					Birthday:  time.Date(2000, time.June, 28, 12, 0, 0, 0, time.UTC),
					Sex:       models.SexFemale,
				}},
			},
		},
		{
			name:      "Error/NotFound",
			id:        test_utils.NumberUUID(1),
			now:       updateTime,
			expectErr: validation.ErrNotFound,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.BeginTx(ctx, nil)
				require.NoError(st, err)
				defer stx.Rollback()

				res, err := NewRepository(stx).Delete(ctx, d.id, d.now)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestUserRepository_Search(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		fixtures []interface{}

		query  string
		offset int
		limit  int

		expect      []*PublicPreview
		expectCount int64
		expectErr   error
	}{
		{
			name: "Success",
			fixtures: concatFixtures(
				// Most relevant result, older.
				generateSearchFixture("", "blue-x", "Elon", "Bezos", models.SexMale, baseTime, 1),
				// Most relevant result, newer.
				generateSearchFixture("", "spaceorigin", "Elon", "Bezos", models.SexMale, baseTime.Add(time.Hour), 2),
				// Slightly relevant result.
				generateSearchFixture("", "notyetwritten", "Eleonore", "Payet", models.SexFemale, baseTime, 4),
				// Irrelevant result.
				generateSearchFixture("", "fruit-basket", "Anna", "Banana", models.SexFemale, baseTime, 5),
			),
			query:       "Ele Be",
			limit:       10,
			expectCount: 3,
			expect: []*PublicPreview{
				{
					Slug:      "spaceorigin",
					FirstName: "Elon",
					LastName:  "Bezos",
					CreatedAt: baseTime.Add(time.Hour),
				},
				{
					Slug:      "blue-x",
					FirstName: "Elon",
					LastName:  "Bezos",
					CreatedAt: baseTime,
				},
				{
					Slug:      "notyetwritten",
					FirstName: "Eleonore",
					LastName:  "Payet",
					CreatedAt: baseTime,
				},
			},
		},
		{
			name: "Success/PreciseQuery",
			fixtures: concatFixtures(
				// Most relevant result, older.
				generateSearchFixture("", "blue-x", "Elon", "Bezos", models.SexMale, baseTime, 1),
				// Most relevant result, newer.
				generateSearchFixture("", "spaceorigin", "Elon", "Bezos", models.SexMale, baseTime.Add(time.Hour), 2),
				// Slightly relevant result.
				generateSearchFixture("", "notyetwritten", "Eleonore", "Payet", models.SexFemale, baseTime, 4),
				// Irrelevant result.
				generateSearchFixture("", "fruit-basket", "Anna", "Banana", models.SexFemale, baseTime, 5),
			),
			query:       "Elon Bezos",
			limit:       10,
			expectCount: 2,
			expect: []*PublicPreview{
				{
					Slug:      "spaceorigin",
					FirstName: "Elon",
					LastName:  "Bezos",
					CreatedAt: baseTime.Add(time.Hour),
				},
				{
					Slug:      "blue-x",
					FirstName: "Elon",
					LastName:  "Bezos",
					CreatedAt: baseTime,
				},
			},
		},
		{
			name: "Success/PreciseQueryReverseOrder",
			fixtures: concatFixtures(
				// Most relevant result, older.
				generateSearchFixture("", "blue-x", "Elon", "Bezos", models.SexMale, baseTime, 1),
				// Most relevant result, newer.
				generateSearchFixture("", "spaceorigin", "Elon", "Bezos", models.SexMale, baseTime.Add(time.Hour), 2),
				// Slightly relevant result.
				generateSearchFixture("", "notyetwritten", "Eleonore", "Payet", models.SexFemale, baseTime, 4),
				// Irrelevant result.
				generateSearchFixture("", "fruit-basket", "Anna", "Banana", models.SexFemale, baseTime, 5),
			),
			query:       "Bezos Elon",
			limit:       10,
			expectCount: 2,
			expect: []*PublicPreview{
				{
					Slug:      "spaceorigin",
					FirstName: "Elon",
					LastName:  "Bezos",
					CreatedAt: baseTime.Add(time.Hour),
				},
				{
					Slug:      "blue-x",
					FirstName: "Elon",
					LastName:  "Bezos",
					CreatedAt: baseTime,
				},
			},
		},
		{
			name: "Success/NoQuery",
			fixtures: concatFixtures(
				// older.
				generateSearchFixture("", "blue-x", "Elon", "Bezos", models.SexMale, baseTime, 1),
				// Most relevant result, newer.
				generateSearchFixture("", "spaceorigin", "Elon", "Bezos", models.SexMale, baseTime.Add(time.Hour), 2),
				// standard.
				generateSearchFixture("", "notyetwritten", "Eleonore", "Payet", models.SexFemale, baseTime.Add(30*time.Minute), 4),
			),
			query:       "",
			limit:       10,
			expectCount: 3,
			expect: []*PublicPreview{
				{
					Slug:      "spaceorigin",
					FirstName: "Elon",
					LastName:  "Bezos",
					CreatedAt: baseTime.Add(time.Hour),
				},
				{
					Slug:      "notyetwritten",
					FirstName: "Eleonore",
					LastName:  "Payet",
					CreatedAt: baseTime.Add(30 * time.Minute),
				},
				{
					Slug:      "blue-x",
					FirstName: "Elon",
					LastName:  "Bezos",
					CreatedAt: baseTime,
				},
			},
		},
		{
			name: "Success/LookUsernameOnlyIfGiven",
			fixtures: concatFixtures(
				// Slightly relevant result.
				generateSearchFixture("eleonore", "blue-x", "Elon", "Bezos", models.SexMale, baseTime, 1),
				// Most relevant result, newer.
				generateSearchFixture("", "spaceorigin", "Elon", "Bezos", models.SexMale, baseTime.Add(time.Hour), 2),
				// Most relevant result, older.
				generateSearchFixture("Elon Bezos", "notyetwritten", "Eleonore", "Payet", models.SexFemale, baseTime, 4),
				// Irrelevant result.
				generateSearchFixture("", "fruit-basket", "Anna", "Banana", models.SexFemale, baseTime, 5),
			),
			query:       "Ele Be",
			limit:       10,
			expectCount: 3,
			expect: []*PublicPreview{
				{
					Slug:      "spaceorigin",
					FirstName: "Elon",
					LastName:  "Bezos",
					CreatedAt: baseTime.Add(time.Hour),
				},
				{
					Username:  "Elon Bezos",
					Slug:      "notyetwritten",
					FirstName: "Eleonore",
					LastName:  "Payet",
					CreatedAt: baseTime,
				},
				{
					Username:  "eleonore",
					Slug:      "blue-x",
					FirstName: "Elon",
					LastName:  "Bezos",
					CreatedAt: baseTime,
				},
			},
		},
		{
			name: "Success/IgnoreAccents",
			fixtures: concatFixtures(
				// Most relevant result, older.
				generateSearchFixture("", "blue-x", "Élon", "Bezos", models.SexMale, baseTime, 1),
				// Most relevant result, newer.
				generateSearchFixture("", "spaceorigin", "Elon", "Bezos", models.SexMale, baseTime.Add(time.Hour), 2),
				// Slightly relevant result.
				generateSearchFixture("", "notyetwritten", "ElËonore", "Payet", models.SexFemale, baseTime, 4),
				// Irrelevant result.
				generateSearchFixture("", "fruit-basket", "Anna", "Banana", models.SexFemale, baseTime, 5),
			),
			query:       "Elè Bé",
			limit:       10,
			expectCount: 3,
			expect: []*PublicPreview{
				{
					Slug:      "spaceorigin",
					FirstName: "Elon",
					LastName:  "Bezos",
					CreatedAt: baseTime.Add(time.Hour),
				},
				{
					Slug:      "blue-x",
					FirstName: "Élon",
					LastName:  "Bezos",
					CreatedAt: baseTime,
				},
				{
					Slug:      "notyetwritten",
					FirstName: "ElËonore",
					LastName:  "Payet",
					CreatedAt: baseTime,
				},
			},
		},
		{
			name: "Success/NoRelevantResults",
			fixtures: concatFixtures(
				// Most relevant result, older.
				generateSearchFixture("", "blue-x", "Elon", "Bezos", models.SexMale, baseTime, 1),
				// Most relevant result, newer.
				generateSearchFixture("", "spaceorigin", "Elon", "Bezos", models.SexMale, baseTime.Add(time.Hour), 2),
				// Slightly relevant result.
				generateSearchFixture("", "notyetwritten", "Eleonore", "Payet", models.SexFemale, baseTime, 4),
				// Irrelevant result.
				generateSearchFixture("", "fruit-basket", "Anna", "Banana", models.SexFemale, baseTime, 5),
			),
			query:  "Abracadabra",
			limit:  10,
			expect: []*PublicPreview(nil),
		},
		{
			name: "Success/PaginationLimit",
			fixtures: concatFixtures(
				// Most relevant result, older.
				generateSearchFixture("", "blue-x", "Elon", "Bezos", models.SexMale, baseTime, 1),
				// Most relevant result, newer.
				generateSearchFixture("", "spaceorigin", "Elon", "Bezos", models.SexMale, baseTime.Add(time.Hour), 2),
				// Slightly relevant result.
				generateSearchFixture("", "notyetwritten", "Eleonore", "Payet", models.SexFemale, baseTime, 4),
				// Irrelevant result.
				generateSearchFixture("", "fruit-basket", "Anna", "Banana", models.SexFemale, baseTime, 5),
			),
			query:       "Ele Be",
			limit:       2,
			expectCount: 3,
			expect: []*PublicPreview{
				{
					Slug:      "spaceorigin",
					FirstName: "Elon",
					LastName:  "Bezos",
					CreatedAt: baseTime.Add(time.Hour),
				},
				{
					Slug:      "blue-x",
					FirstName: "Elon",
					LastName:  "Bezos",
					CreatedAt: baseTime,
				},
			},
		},
		{
			name: "Success/PaginationOffset",
			fixtures: concatFixtures(
				// Most relevant result, older.
				generateSearchFixture("", "blue-x", "Elon", "Bezos", models.SexMale, baseTime, 1),
				// Most relevant result, newer.
				generateSearchFixture("", "spaceorigin", "Elon", "Bezos", models.SexMale, baseTime.Add(time.Hour), 2),
				// Slightly relevant result.
				generateSearchFixture("", "notyetwritten", "Eleonore", "Payet", models.SexFemale, baseTime, 4),
				// Irrelevant result.
				generateSearchFixture("", "fruit-basket", "Anna", "Banana", models.SexFemale, baseTime, 5),
			),
			query:       "Ele Be",
			offset:      2,
			limit:       10,
			expectCount: 3,
			expect: []*PublicPreview{
				{
					Slug:      "notyetwritten",
					FirstName: "Eleonore",
					LastName:  "Payet",
					CreatedAt: baseTime,
				},
			},
		},
		{
			name: "Success/OnSlug",
			fixtures: concatFixtures(
				// Most relevant result, older.
				generateSearchFixture("", "blue-x", "Elon", "Bezos", models.SexMale, baseTime, 1),
				// Most relevant result, newer.
				generateSearchFixture("", "spaceorigin", "Elon", "Bezos", models.SexMale, baseTime.Add(time.Hour), 2),
				// Slightly relevant result.
				generateSearchFixture("", "notyetwritten", "Eleonore", "Payet", models.SexFemale, baseTime, 4),
				// Irrelevant result.
				generateSearchFixture("", "fruit-basket", "Anna", "Banana", models.SexFemale, baseTime, 5),
			),
			query:       "fruitbasket",
			limit:       10,
			expectCount: 1,
			expect: []*PublicPreview{
				{
					Slug:      "fruit-basket",
					FirstName: "Anna",
					LastName:  "Banana",
					CreatedAt: baseTime,
				},
			},
		},
	}

	for _, d := range data {
		err := test_utils.RunTransactionalTest(db, d.fixtures, func(ctx context.Context, tx bun.Tx) {
			t.Run(d.name, func(st *testing.T) {
				res, count, err := NewRepository(tx).Search(ctx, d.query, d.limit, d.offset)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
				require.Equal(t, d.expectCount, count)
			})
		})
		require.NoError(t, err)
	}
}

func TestUserRepository_GetPreview(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name      string
		id        uuid.UUID
		expect    *Preview
		expectErr error
	}{
		{
			name: "Success",
			id:   test_utils.NumberUUID(1000),
			expect: &Preview{
				ID:        test_utils.NumberUUID(1000),
				Username:  "BanAnna",
				FirstName: "Anna",
				LastName:  "Banana",
				Email: models.Email{
					User:   "anna.banana",
					Domain: "coco.nut",
				},
				Slug: "fruit-basket",
				Sex:  models.SexFemale,
			},
		},
		{
			name:      "Error/NotFound",
			id:        test_utils.NumberUUID(1001),
			expectErr: validation.ErrNotFound,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := NewRepository(tx).GetPreview(ctx, d.id)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestUserRepository_GetPublic(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name      string
		slug      string
		expect    *Public
		expectErr error
	}{
		{
			name: "Success",
			slug: "fruit-basket",
			expect: &Public{
				ID:        test_utils.NumberUUID(1000),
				Username:  "BanAnna",
				FirstName: "Anna",
				LastName:  "Banana",
				CreatedAt: baseTime,
				Sex:       models.SexFemale,
			},
		},
		{
			name:      "Error/NotFound",
			slug:      "not-found",
			expectErr: validation.ErrNotFound,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := NewRepository(tx).GetPublic(ctx, d.slug)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestUserRepository_GetPublicPreviews(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name      string
		ids       []uuid.UUID
		fixtures  []interface{}
		expect    []*PublicPreview
		expectErr error
	}{
		{
			name: "Success",
			fixtures: concatFixtures(
				// Most relevant result, older.
				generateSearchFixture("BigBrother", "blue-x", "Elon", "Bezos", models.SexMale, baseTime, 1),
				// Most relevant result, newer.
				generateSearchFixture("", "spaceorigin", "Elon", "Bezos", models.SexMale, baseTime.Add(time.Hour), 2),
				// Slightly relevant result.
				generateSearchFixture("3", "notyetwritten", "Eleonore", "Payet", models.SexFemale, baseTime, 4),
				// Irrelevant result.
				generateSearchFixture("Banananana", "fruit-basket", "Anna", "Banana", models.SexFemale, baseTime, 5),
			),
			ids: []uuid.UUID{
				test_utils.NumberUUID(1),
				test_utils.NumberUUID(2),
				test_utils.NumberUUID(5),
				// Don't exist.
				test_utils.NumberUUID(15),
			},
			expect: []*PublicPreview{
				{
					ID:        test_utils.NumberUUID(1),
					Slug:      "blue-x",
					Username:  "BigBrother",
					FirstName: "Elon",
					LastName:  "Bezos",
					CreatedAt: baseTime,
				},
				{
					ID:        test_utils.NumberUUID(2),
					Slug:      "spaceorigin",
					FirstName: "Elon",
					LastName:  "Bezos",
					CreatedAt: baseTime.Add(time.Hour),
				},
				{
					ID:        test_utils.NumberUUID(5),
					Slug:      "fruit-basket",
					Username:  "Banananana",
					FirstName: "Anna",
					LastName:  "Banana",
					CreatedAt: baseTime,
				},
			},
		},
		{
			name: "Success/NoResults",
			fixtures: concatFixtures(
				// Most relevant result, older.
				generateSearchFixture("BigBrother", "blue-x", "Elon", "Bezos", models.SexMale, baseTime, 1),
				// Most relevant result, newer.
				generateSearchFixture("", "spaceorigin", "Elon", "Bezos", models.SexMale, baseTime.Add(time.Hour), 2),
				// Slightly relevant result.
				generateSearchFixture("3", "notyetwritten", "Eleonore", "Payet", models.SexFemale, baseTime, 4),
				// Irrelevant result.
				generateSearchFixture("Banananana", "fruit-basket", "Anna", "Banana", models.SexFemale, baseTime, 5),
			),
			ids: []uuid.UUID{
				// Don't exist.
				test_utils.NumberUUID(15),
			},
			expect: []*PublicPreview{},
		},
	}

	for _, d := range data {
		err := test_utils.RunTransactionalTest(db, d.fixtures, func(ctx context.Context, tx bun.Tx) {
			t.Run(d.name, func(st *testing.T) {
				res, err := NewRepository(tx).GetPublicPreviews(ctx, d.ids)
				test_utils.RequireError(t, d.expectErr, err)
				sort.Slice(res, func(i, j int) bool {
					return res[i].Slug < res[j].Slug
				})
				sort.Slice(d.expect, func(i, j int) bool {
					return d.expect[i].Slug < d.expect[j].Slug
				})
				require.Equal(t, d.expect, res)
			})
		})
		require.NoError(t, err)
	}
}
