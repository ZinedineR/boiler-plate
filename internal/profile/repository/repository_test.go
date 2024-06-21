package repository_test

import (
	"boiler-plate/internal/profile/domain"
	"boiler-plate/internal/profile/repository"
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
	"time"
)

func setupSQLMock(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	// Setup SQL mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// Setup GORM with the mock DB
	gormDB, gormDBErr := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if gormDBErr != nil {
		t.Fatalf("failed to open GORM connection: %v", gormDBErr)
	}
	return mock, gormDB
}

func TestProfileRepository_CreateUser(t *testing.T) {
	// Setup SQL mock
	mock, gormDB := setupSQLMock(t)

	// Initialize ProfileRepository with mocked GORM connection
	//userRepo := postgres_gorm.NewProfileRepository(gormDB)
	profileRepo := repository.NewRepository(gormDB, nil)

	expectedQueryString := regexp.QuoteMeta(`INSERT INTO "profile" ("profile","password","created_at") VALUES ($1,$2,$3) RETURNING "id"`)
	now := time.Now()
	t.Run("Positive Case", func(t *testing.T) {
		// Expected user data to insert

		profile := &domain.Profile{
			Profile:   "Zinedine",
			Password:  "test",
			CreatedAt: &now,
		}

		// Set mock expectations for the transaction
		mock.ExpectQuery(expectedQueryString).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(1)) // Mock the result of the INSERT operation

		// Call the CreateUser method
		err := profileRepo.Create(context.Background(), gormDB, profile)

		// Assert the result
		require.NoError(t, err)
		require.NotNil(t, profile.ID)
		//require.Equal(t, user.Name, createdUser.Name)
		//require.Equal(t, user.Email, createdUser.Email)
	})

	t.Run("Negative Case", func(t *testing.T) {
		// Expected user data to insert
		profile := &domain.Profile{
			Profile:   "Zinedine",
			Password:  "test",
			CreatedAt: &now,
		}

		// Set mock expectations for the transaction
		mock.ExpectQuery(expectedQueryString).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		// Call the CreateUser method
		err := profileRepo.Create(context.Background(), gormDB, profile)

		// Assert the result
		require.Error(t, err)
	})
}

func TestProfileRepository_DeleteUser(t *testing.T) {
	mock, gormDB := setupSQLMock(t)
	profileRepo := repository.NewRepository(gormDB, nil)

	expectedQueryString := regexp.QuoteMeta(`DELETE FROM "profile" WHERE "profile"."id" = $1`)

	t.Run("Positive Case", func(t *testing.T) {
		// Expected profile ID to delete
		profileID := 1

		// Set mock expectations for the transaction
		mock.ExpectExec(expectedQueryString).
			WithArgs(profileID).
			WillReturnResult(sqlmock.NewResult(1, 1)) // Mock the result of the DELETE operation

		// Call the Delete method
		err := profileRepo.Delete(context.Background(), gormDB, profileID)

		// Assert the result
		require.NoError(t, err)
	})

	t.Run("Negative Case", func(t *testing.T) {
		// Expected profile ID to delete
		profileID := 1

		// Set mock expectations for the transaction
		mock.ExpectExec(expectedQueryString).
			WithArgs(profileID).
			WillReturnError(errors.New("db error"))

		// Call the Delete method
		err := profileRepo.Delete(context.Background(), gormDB, profileID)

		// Assert the result
		require.Error(t, err)
		require.EqualError(t, err, "db error")
	})
}

func TestProfileRepository_UpdateUser(t *testing.T) {
	mock, gormDB := setupSQLMock(t)

	profileRepo := repository.NewRepository(gormDB, nil)

	expectedQueryString := regexp.QuoteMeta(`UPDATE "profile" SET "profile"=$1,"password"=$2,"created_at"=$3 WHERE "id" = $4`)

	now := time.Now()

	t.Run("Positive Case", func(t *testing.T) {
		profileId := 1
		profile := &domain.Profile{
			//ID:        1,
			Profile:   "Zinedine",
			Password:  "updated_password",
			CreatedAt: &now,
		}
		fmt.Println(profile.TableName())
		mock.ExpectExec(expectedQueryString).
			WithArgs(profile.Profile, profile.Password, profile.CreatedAt, profileId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := profileRepo.Update(context.Background(), gormDB, profileId, profile)

		require.NoError(t, err)
	})

	t.Run("Negative Case", func(t *testing.T) {
		profileID := 1
		profile := &domain.Profile{
			Profile:   "Zinedine",
			Password:  "updated_password",
			CreatedAt: &now,
		}

		mock.ExpectExec(expectedQueryString).
			WithArgs(profile.Profile, profile.Password, profile.CreatedAt, profileID).
			WillReturnError(errors.New("db error"))

		err := profileRepo.Update(context.Background(), gormDB, profileID, profile)

		require.Error(t, err)
		require.EqualError(t, err, "db error")
	})
}

func TestProfileRepository_FindUsers(t *testing.T) {
	mock, gormDB := setupSQLMock(t)
	profileRepo := repository.NewRepository(gormDB, nil)

	expectedQueryString := regexp.QuoteMeta(`SELECT * FROM "profile"`)

	t.Run("Positive Case", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "profile", "password", "created_at"}).
			AddRow(1, "Zinedine", "password1", time.Now()).
			AddRow(2, "Ronaldo", "password2", time.Now())

		mock.ExpectQuery(expectedQueryString).
			WillReturnRows(rows)

		profiles, err := profileRepo.Find(context.Background(), gormDB)

		require.NoError(t, err)
		require.NotNil(t, profiles)
		require.Len(t, *profiles, 2)
	})

	t.Run("Negative Case - DB Error", func(t *testing.T) {
		mock.ExpectQuery(expectedQueryString).
			WillReturnError(errors.New("db error"))

		profiles, err := profileRepo.Find(context.Background(), gormDB)

		require.Error(t, err)
		require.Nil(t, profiles)
		require.EqualError(t, err, "db error")
	})

	t.Run("Negative Case - Record Not Found", func(t *testing.T) {
		mock.ExpectQuery(expectedQueryString).
			WillReturnRows(sqlmock.NewRows([]string{"id", "profile", "password", "created_at"})).
			WillReturnError(gorm.ErrRecordNotFound)

		profiles, err := profileRepo.Find(context.Background(), gormDB)
		require.NoError(t, err)
		require.NotNil(t, profiles)
		require.Len(t, *profiles, 0)
	})
}

func TestProfileRepository_DetailUser(t *testing.T) {
	mock, gormDB := setupSQLMock(t)
	profileRepo := repository.NewRepository(gormDB, nil)

	expectedQueryString := regexp.QuoteMeta(`SELECT * FROM "profile" WHERE "profile"."id" = $1 ORDER BY "profile"."id" LIMIT 1`)

	t.Run("Positive Case", func(t *testing.T) {
		profileID := 1
		rows := sqlmock.NewRows([]string{"id", "profile", "password", "created_at"}).
			AddRow(profileID, "Zinedine", "password", time.Now())

		mock.ExpectQuery(expectedQueryString).
			WithArgs(profileID).
			WillReturnRows(rows)

		profile, err := profileRepo.Detail(context.Background(), gormDB, profileID)

		require.NoError(t, err)
		require.NotNil(t, profile)
		require.Equal(t, profileID, profile.ID)
		require.Equal(t, "Zinedine", profile.Profile)
	})

	t.Run("Negative Case - DB Error", func(t *testing.T) {
		profileID := 1

		mock.ExpectQuery(expectedQueryString).
			WithArgs(profileID).
			WillReturnError(errors.New("db error"))

		profile, err := profileRepo.Detail(context.Background(), gormDB, profileID)

		require.Error(t, err)
		require.Nil(t, profile)
		require.EqualError(t, err, "db error")
	})

	t.Run("Negative Case - Record Not Found", func(t *testing.T) {
		profileID := 1

		mock.ExpectQuery(expectedQueryString).
			WithArgs(profileID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "profile", "password", "created_at"}))

		profile, err := profileRepo.Detail(context.Background(), gormDB, profileID)

		require.NoError(t, err)
		require.Nil(t, profile)
	})
}

func TestProfileRepository_AuthUser(t *testing.T) {
	mock, gormDB := setupSQLMock(t)
	profileRepo := repository.NewRepository(gormDB, nil)

	expectedQueryString := regexp.QuoteMeta(`SELECT * FROM "profile" WHERE profile = $1 AND password = $2 ORDER BY "profile"."id" LIMIT 1`)

	t.Run("Positive Case", func(t *testing.T) {
		profile := "Zinedine"
		password := "password"
		rows := sqlmock.NewRows([]string{"id", "profile", "password", "created_at"}).
			AddRow(1, profile, password, time.Now())

		mock.ExpectQuery(expectedQueryString).
			WithArgs(profile, password).
			WillReturnRows(rows)

		authProfile, err := profileRepo.Auth(context.Background(), gormDB, profile, password)

		require.NoError(t, err)
		require.NotNil(t, authProfile)
		require.Equal(t, profile, authProfile.Profile)
		require.Equal(t, password, authProfile.Password)
	})

	t.Run("Negative Case - DB Error", func(t *testing.T) {
		profile := "Zinedine"
		password := "password"

		mock.ExpectQuery(expectedQueryString).
			WithArgs(profile, password).
			WillReturnError(errors.New("db error"))

		authProfile, err := profileRepo.Auth(context.Background(), gormDB, profile, password)

		require.Error(t, err)
		require.Nil(t, authProfile)
		require.EqualError(t, err, "db error")
	})

	t.Run("Negative Case - Record Not Found", func(t *testing.T) {
		profile := "Zinedine"
		password := "password"

		mock.ExpectQuery(expectedQueryString).
			WithArgs(profile, password).
			WillReturnRows(sqlmock.NewRows([]string{"id", "profile", "password", "created_at"}))

		authProfile, err := profileRepo.Auth(context.Background(), gormDB, profile, password)

		require.NoError(t, err)
		require.Nil(t, authProfile)
	})
}
