package credentials_storage

import (
	"context"
	"github.com/a-novel/agora-backend/framework/test"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"testing"
	"time"
)

var (
	baseTime   = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
	updateTime = time.Date(2020, time.May, 4, 9, 0, 0, 0, time.UTC)
)

var Fixtures = []*Model{
	// Standard user.
	{
		ID:        test_utils.NumberUUID(1000),
		CreatedAt: baseTime,
		UpdatedAt: &baseTime,
		Core: Core{
			Email: models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			Password: models.Password{
				Hashed: "foobarqux",
			},
		},
	},
	// Password pending reset.
	{
		ID:        test_utils.NumberUUID(1001),
		CreatedAt: baseTime,
		UpdatedAt: &baseTime,
		Core: Core{
			Email: models.Email{
				User:   "bill.cook",
				Domain: "amazon.com",
			},
			Password: models.Password{
				Validation: "youshallpass",
				Hashed:     "foobarqux",
			},
		},
	},
	// Main email pending validation.
	{
		ID:        test_utils.NumberUUID(1002),
		CreatedAt: baseTime,
		UpdatedAt: &baseTime,
		Core: Core{
			Email: models.Email{
				Validation: "youshallpass",
				User:       "joe.doe",
				Domain:     "poe.co",
			},
			Password: models.Password{
				Hashed: "foobarqux",
			},
		},
	},
	// New email pending validation.
	{
		ID:        test_utils.NumberUUID(1003),
		CreatedAt: baseTime,
		UpdatedAt: &baseTime,
		Core: Core{
			Email: models.Email{
				User:   "isaac.asimov",
				Domain: "terminus.gal",
			},
			NewEmail: models.Email{
				Validation: "youshallpass",
				User:       "letter.number",
				Domain:     "alphabet.xyz",
			},
			Password: models.Password{
				Hashed: "foobarqux",
			},
		},
	},
	// New email pending validation, but taken.
	{
		ID:        test_utils.NumberUUID(1004),
		CreatedAt: baseTime,
		UpdatedAt: &baseTime,
		Core: Core{
			Email: models.Email{
				User:   "java",
				Domain: "satan.hell",
			},
			NewEmail: models.Email{
				Validation: "youshallpass",
				User:       "elon.bezos",
				Domain:     "gmail.com",
			},
			Password: models.Password{
				Hashed: "foobarqux",
			},
		},
	},
	// Main and new emails pending validation.
	{
		ID:        test_utils.NumberUUID(1005),
		CreatedAt: baseTime,
		UpdatedAt: &baseTime,
		Core: Core{
			Email: models.Email{
				Validation: "youshallpass",
				User:       "potato",
				Domain:     "food.fr",
			},
			NewEmail: models.Email{
				Validation: "youshallpass",
				User:       "strawberry",
				Domain:     "food.fr",
			},
			Password: models.Password{
				Hashed: "foobarqux",
			},
		},
	},
}

func TestCredentialsRepository_Read(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		id uuid.UUID

		expect    *Model
		expectErr error
	}{
		{
			name:   "Success",
			id:     test_utils.NumberUUID(1000),
			expect: Fixtures[0],
		},
		{
			name:      "Error/NotFound",
			id:        test_utils.NumberUUID(1),
			expectErr: validation.ErrNotFound,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.Read(ctx, d.id)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_ReadEmail(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		email models.Email

		expect    *Model
		expectErr error
	}{
		{
			name: "Success",
			email: models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			expect: Fixtures[0],
		},
		{
			name: "Error/NotFound",
			email: models.Email{
				User:   "jeff.musk",
				Domain: "gmail.com",
			},
			expectErr: validation.ErrNotFound,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.ReadEmail(ctx, d.email)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_EmailExists(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		email models.Email

		expect    bool
		expectErr error
	}{
		{
			name: "Success/Exists",
			email: models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			expect: true,
		},
		{
			name: "Success/DoesNotExists",
			email: models.Email{
				User:   "jeff.musk",
				Domain: "gmail.com",
			},
			expect: false,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.EmailExists(ctx, d.email)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_UpdateEmail(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		email models.Email
		code  string
		id    uuid.UUID
		now   time.Time

		expect    *Model
		expectErr error
	}{
		{
			name: "Success",
			email: models.Email{
				Validation: "thisshouldbeignored",
				User:       "123",
				Domain:     "nya.arigatou",
			},
			code: "lyoko",
			id:   test_utils.NumberUUID(1000),
			now:  updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "lyoko",
						User:       "123",
						Domain:     "nya.arigatou",
					},
					Password: models.Password{
						Hashed: "foobarqux",
					},
				},
			},
		},
		{
			name: "Success/WithValidationOnMainEmail",
			email: models.Email{
				Validation: "thisshouldbeignored",
				User:       "123",
				Domain:     "nya.arigatou",
			},
			code: "lyoko",
			id:   test_utils.NumberUUID(1002),
			now:  updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1002),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						Validation: "youshallpass",
						User:       "joe.doe",
						Domain:     "poe.co",
					},
					NewEmail: models.Email{
						Validation: "lyoko",
						User:       "123",
						Domain:     "nya.arigatou",
					},
					Password: models.Password{
						Hashed: "foobarqux",
					},
				},
			},
		},
		{
			name: "Success/WithPreviousNewEmail",
			email: models.Email{
				Validation: "thisshouldbeignored",
				User:       "123",
				Domain:     "nya.arigatou",
			},
			code: "lyoko",
			id:   test_utils.NumberUUID(1003),
			now:  updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1003),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						User:   "isaac.asimov",
						Domain: "terminus.gal",
					},
					NewEmail: models.Email{
						Validation: "lyoko",
						User:       "123",
						Domain:     "nya.arigatou",
					},
					Password: models.Password{
						Hashed: "foobarqux",
					},
				},
			},
		},
		{
			name: "Success/WhenAnotherAccountHasTheSameEmailPendingValidation",
			email: models.Email{
				Validation: "youshallpass",
				User:       "letter.number",
				Domain:     "alphabet.xyz",
			},
			code: "lyoko",
			id:   test_utils.NumberUUID(1000),
			now:  updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "lyoko",
						User:       "letter.number",
						Domain:     "alphabet.xyz",
					},
					Password: models.Password{
						Hashed: "foobarqux",
					},
				},
			},
		},
		// Illegal from a business standpoint, but this layer should not be the one validating the data. There is no
		// data integrity violation as long as the main email is unique for each user. The pending email update has no
		// such obligations.
		{
			name: "Success/TakenByAnotherAccount",
			email: models.Email{
				User:   "bill.cook",
				Domain: "amazon.com",
			},
			code: "lyoko",
			id:   test_utils.NumberUUID(1000),
			now:  updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "lyoko",
						User:       "bill.cook",
						Domain:     "amazon.com",
					},
					Password: models.Password{
						Hashed: "foobarqux",
					},
				},
			},
		},
		{
			name: "Error/NotFound",
			email: models.Email{
				Validation: "thisshouldbeignored",
				User:       "123",
				Domain:     "nya.arigatou",
			},
			code:      "lyoko",
			id:        test_utils.NumberUUID(100),
			now:       updateTime,
			expectErr: validation.ErrNotFound,
		},
		{
			name: "Error/WithoutValidationCode",
			email: models.Email{
				User:   "123",
				Domain: "nya.arigatou",
			},
			id:        test_utils.NumberUUID(1000),
			now:       updateTime,
			expectErr: validation.ErrConstraintViolation,
		},
		{
			name: "Error/WithoutUser",
			email: models.Email{
				Domain: "nya.arigatou",
			},
			code:      "lyoko",
			id:        test_utils.NumberUUID(1000),
			now:       updateTime,
			expectErr: validation.ErrConstraintViolation,
		},
		{
			name: "Error/WithoutDomain",
			email: models.Email{
				User: "123",
			},
			code:      "lyoko",
			id:        test_utils.NumberUUID(1000),
			now:       updateTime,
			expectErr: validation.ErrConstraintViolation,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.BeginTx(ctx, nil)
				require.NoError(st, err)
				defer stx.Rollback()

				res, err := NewRepository(stx).UpdateEmail(ctx, d.email, d.code, d.id, d.now)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_ValidateEmail(t *testing.T) {
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
			id:   test_utils.NumberUUID(1002),
			now:  updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1002),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						User:   "joe.doe",
						Domain: "poe.co",
					},
					Password: models.Password{
						Hashed: "foobarqux",
					},
				},
			},
		},
		{
			name: "Success/WithEmailPendingValidation",
			id:   test_utils.NumberUUID(1005),
			now:  updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1005),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						User:   "potato",
						Domain: "food.fr",
					},
					NewEmail: models.Email{
						Validation: "youshallpass",
						User:       "strawberry",
						Domain:     "food.fr",
					},
					Password: models.Password{
						Hashed: "foobarqux",
					},
				},
			},
		},
		{
			name:      "Error/NoPendingValidation",
			id:        test_utils.NumberUUID(1000),
			now:       updateTime,
			expectErr: validation.ErrNotFound,
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

				res, err := NewRepository(stx).ValidateEmail(ctx, d.id, d.now)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_ValidateNewEmail(t *testing.T) {
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
			id:   test_utils.NumberUUID(1003),
			now:  updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1003),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						User:   "letter.number",
						Domain: "alphabet.xyz",
					},
					Password: models.Password{
						Hashed: "foobarqux",
					},
				},
			},
		},
		{
			name: "Success/WithMainEmailPendingValidation",
			id:   test_utils.NumberUUID(1005),
			now:  updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1005),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						User:   "strawberry",
						Domain: "food.fr",
					},
					Password: models.Password{
						Hashed: "foobarqux",
					},
				},
			},
		},
		{
			name:      "Error/NoPendingValidation",
			id:        test_utils.NumberUUID(1000),
			now:       updateTime,
			expectErr: validation.ErrNotFound,
		},
		{
			name:      "Error/NotFound",
			id:        test_utils.NumberUUID(1),
			now:       updateTime,
			expectErr: validation.ErrNotFound,
		},
		{
			name:      "Error/AlreadyTaken",
			id:        test_utils.NumberUUID(1004),
			now:       updateTime,
			expectErr: validation.ErrUniqConstraintViolation,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.BeginTx(ctx, nil)
				require.NoError(st, err)
				defer stx.Rollback()

				res, err := NewRepository(stx).ValidateNewEmail(ctx, d.id, d.now)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_UpdatePassword(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		newPassword string
		id          uuid.UUID
		now         time.Time

		expect    *Model
		expectErr error
	}{
		{
			name:        "Success",
			newPassword: "hackmeifyoucan",
			id:          test_utils.NumberUUID(1000),
			now:         updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					Password: models.Password{
						Hashed: "hackmeifyoucan",
					},
				},
			},
		},
		{
			name:        "Success/WithPendingReset",
			newPassword: "hackmeifyoucan",
			id:          test_utils.NumberUUID(1001),
			now:         updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1001),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						User:   "bill.cook",
						Domain: "amazon.com",
					},
					Password: models.Password{
						Hashed: "hackmeifyoucan",
					},
				},
			},
		},
		{
			name:        "Error/NotFound",
			newPassword: "hackmeifyoucan",
			id:          test_utils.NumberUUID(100),
			now:         updateTime,
			expectErr:   validation.ErrNotFound,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.BeginTx(ctx, nil)
				require.NoError(st, err)
				defer stx.Rollback()

				res, err := NewRepository(stx).UpdatePassword(ctx, d.newPassword, d.id, d.now)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_ResetPassword(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		code  string
		email models.Email
		now   time.Time

		expect    *Model
		expectErr error
	}{
		{
			name: "Success",
			code: "lyoko",
			email: models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			now: updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					Password: models.Password{
						Validation: "lyoko",
						Hashed:     "foobarqux",
					},
				},
			},
		},
		{
			name: "Success/WithPreviousResetPending",
			code: "lyoko",
			email: models.Email{
				User:   "bill.cook",
				Domain: "amazon.com",
			},
			now: updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1001),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						User:   "bill.cook",
						Domain: "amazon.com",
					},
					Password: models.Password{
						Validation: "lyoko",
						Hashed:     "foobarqux",
					},
				},
			},
		},
		{
			name: "Error/NotFound",
			code: "lyoko",
			email: models.Email{
				User:   "elon.gates",
				Domain: "yahoo.com",
			},
			now:       updateTime,
			expectErr: validation.ErrNotFound,
		},
		{
			name: "Error/NotFoundIfEmailPendingUpdate",
			code: "lyoko",
			email: models.Email{
				Validation: "youshallpass",
				User:       "letter.number",
				Domain:     "alphabet.xyz",
			},
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

				res, err := NewRepository(stx).ResetPassword(ctx, d.code, d.email, d.now)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_UpdateEmailValidation(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		code string
		id   uuid.UUID
		now  time.Time

		expect    *Model
		expectErr error
	}{
		{
			name: "Success",
			code: "lyoko",
			id:   test_utils.NumberUUID(1002),
			now:  updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1002),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						Validation: "lyoko",
						User:       "joe.doe",
						Domain:     "poe.co",
					},
					Password: models.Password{
						Hashed: "foobarqux",
					},
				},
			},
		},
		{
			name: "Success/WithAnEmailPendingUpdate",
			code: "lyoko",
			id:   test_utils.NumberUUID(1005),
			now:  updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1005),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						Validation: "lyoko",
						User:       "potato",
						Domain:     "food.fr",
					},
					NewEmail: models.Email{
						Validation: "youshallpass",
						User:       "strawberry",
						Domain:     "food.fr",
					},
					Password: models.Password{
						Hashed: "foobarqux",
					},
				},
			},
		},
		{
			name:      "Error/NotFound",
			code:      "lyoko",
			id:        test_utils.NumberUUID(100),
			now:       updateTime,
			expectErr: validation.ErrNotFound,
		},
		{
			name:      "Error/AlreadyValidated",
			code:      "lyoko",
			id:        test_utils.NumberUUID(1000),
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

				res, err := NewRepository(stx).UpdateEmailValidation(ctx, d.code, d.id, d.now)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_UpdateNewEmailValidation(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		code string
		id   uuid.UUID
		now  time.Time

		expect    *Model
		expectErr error
	}{
		{
			name: "Success",
			code: "lyoko",
			id:   test_utils.NumberUUID(1003),
			now:  updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1003),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						User:   "isaac.asimov",
						Domain: "terminus.gal",
					},
					NewEmail: models.Email{
						Validation: "lyoko",
						User:       "letter.number",
						Domain:     "alphabet.xyz",
					},
					Password: models.Password{
						Hashed: "foobarqux",
					},
				},
			},
		},
		{
			name: "Success/WithMainEmailPendingValidation",
			code: "lyoko",
			id:   test_utils.NumberUUID(1005),
			now:  updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1005),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						Validation: "youshallpass",
						User:       "potato",
						Domain:     "food.fr",
					},
					NewEmail: models.Email{
						Validation: "lyoko",
						User:       "strawberry",
						Domain:     "food.fr",
					},
					Password: models.Password{
						Hashed: "foobarqux",
					},
				},
			},
		},
		{
			name:      "Error/NoValidationCode",
			id:        test_utils.NumberUUID(1003),
			now:       updateTime,
			expectErr: validation.ErrConstraintViolation,
		},
		{
			name:      "Error/NotFound",
			code:      "lyoko",
			id:        test_utils.NumberUUID(100),
			now:       updateTime,
			expectErr: validation.ErrNotFound,
		},
		{
			name:      "Error/NoPendingUpdate",
			code:      "lyoko",
			id:        test_utils.NumberUUID(1000),
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

				res, err := NewRepository(stx).UpdateNewEmailValidation(ctx, d.code, d.id, d.now)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_CancelNewEmail(t *testing.T) {
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
			id:   test_utils.NumberUUID(1003),
			now:  updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1003),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						User:   "isaac.asimov",
						Domain: "terminus.gal",
					},
					Password: models.Password{
						Hashed: "foobarqux",
					},
				},
			},
		},
		{
			name: "Success/WithMainEmailPendingValidation",
			id:   test_utils.NumberUUID(1005),
			now:  updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1005),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						Validation: "youshallpass",
						User:       "potato",
						Domain:     "food.fr",
					},
					Password: models.Password{
						Hashed: "foobarqux",
					},
				},
			},
		},
		{
			name: "Success/NoPendingUpdate",
			id:   test_utils.NumberUUID(1000),
			now:  updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					Password: models.Password{
						Hashed: "foobarqux",
					},
				},
			},
		},
		{
			name:      "Error/NotFound",
			id:        test_utils.NumberUUID(100),
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

				res, err := NewRepository(stx).CancelNewEmail(ctx, d.id, d.now)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}
