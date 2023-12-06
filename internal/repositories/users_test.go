package repositories

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/qthuy2k1/product-management/internal/models"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UserRepository_CreateUser(t *testing.T) {
	// test cases
	testCases := map[string]struct {
		input User
		err   error
	}{
		"success": {
			input: User{
				Name:     "Quang Thuy",
				Email:    "qthuy@gmail.ok",
				Password: "213123",
				Role:     "user",
				Status:   "activated",
			},
		},
		"create user with name longer than 255 characters": {
			input: User{
				Name:     "Quang Thuy" + strings.Repeat("a", 251),
				Email:    "qthuy@gmail.com",
				Password: "123123",
			},
			err: errors.New("models: unable to insert into users: pq: value too long for type character varying(255)"),
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			// Given
			db, err := Initialize(os.Getenv("DB_URL"))
			assert.NoError(t, err)
			boil.SetDB(db)
			defer db.Exec("DELETE FROM users;")

			redis := RedisInitialize(os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASS"))

			repo := NewRepository(db, redis)

			// When
			err = repo.CreateUser(context.Background(), tc.input)

			// Then
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			}
		})
	}
}

func Test_UserRepository_GetUser(t *testing.T) {
	// test cases
	testCases := map[string]struct {
		input     int
		expResult models.User
		err       error
	}{
		"success": {
			input: 1001,
			expResult: models.User{
				ID:       1001,
				Name:     "John Doe",
				Email:    "doe@example.com",
				Password: "123123",
				Role:     "user",
				Status:   "activated",
			},
		},
		"user not found": {
			input: -1,
			err:   ErrUserNotFound,
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			// Given
			db, dbErr := Initialize(os.Getenv("DB_URL"))
			assert.NoError(t, dbErr)
			boil.SetDB(db)

			err := runSQLTest(db, "./datatest/users/insert_user.sql")
			assert.NoError(t, err)

			defer runSQLTest(db, "./datatest/users/rollback_insert_user.sql")

			redis := RedisInitialize(os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASS"))

			repo := NewRepository(db, redis)

			// When
			result, err := repo.GetUser(context.Background(), tc.input)

			// Then
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)

				result.CreatedAt = time.Time{}
				result.UpdatedAt = time.Time{}

				require.Equal(t, tc.expResult, result)
			}
		})
	}
}
