package repository_test

import (
	"boiler-plate/internal/users/domain"
	"boiler-plate/internal/users/repository"
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

func TestUsersRepository_CreateUser(t *testing.T) {
	// Setup SQL mock
	mock, gormDB := setupSQLMock(t)

	// Initialize UsersRepository with mocked GORM connection
	//userRepo := postgres_gorm.NewUsersRepository(gormDB)
	usersRepo := repository.NewRepository(gormDB, nil)

	expectedQueryString := regexp.QuoteMeta(`INSERT INTO "users" ("email","password","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`)
	t.Run("Positive Case", func(t *testing.T) {
		// Expected user data to insert

		users := &domain.Users{
			Email:    "Zinedine",
			Password: "test",
		}

		// Set mock expectations for the transaction
		mock.ExpectQuery(expectedQueryString).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(1)) // Mock the result of the INSERT operation

		// Call the CreateUser method
		err := usersRepo.Create(context.Background(), gormDB, users)

		// Assert the result
		require.NoError(t, err)
		require.NotNil(t, users.ID)
	})

	t.Run("Negative Case", func(t *testing.T) {
		// Expected user data to insert
		users := &domain.Users{
			Email:    "Zinedine",
			Password: "test",
		}

		// Set mock expectations for the transaction
		mock.ExpectQuery(expectedQueryString).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		// Call the CreateUser method
		err := usersRepo.Create(context.Background(), gormDB, users)

		// Assert the result
		require.Error(t, err)
	})
}

func TestUsersRepository_DeleteUser(t *testing.T) {
	mock, gormDB := setupSQLMock(t)
	usersRepo := repository.NewRepository(gormDB, nil)

	expectedQueryString := regexp.QuoteMeta(`DELETE FROM "users" WHERE "users"."id" = $1`)

	t.Run("Positive Case", func(t *testing.T) {
		// Expected users ID to delete
		usersID := 1

		// Set mock expectations for the transaction
		mock.ExpectExec(expectedQueryString).
			WithArgs(usersID).
			WillReturnResult(sqlmock.NewResult(1, 1)) // Mock the result of the DELETE operation

		// Call the Delete method
		err := usersRepo.Delete(context.Background(), gormDB, usersID)

		// Assert the result
		require.NoError(t, err)
	})

	t.Run("Negative Case", func(t *testing.T) {
		// Expected users ID to delete
		usersID := 1

		// Set mock expectations for the transaction
		mock.ExpectExec(expectedQueryString).
			WithArgs(usersID).
			WillReturnError(errors.New("db error"))

		// Call the Delete method
		err := usersRepo.Delete(context.Background(), gormDB, usersID)

		// Assert the result
		require.Error(t, err)
		require.EqualError(t, err, "db error")
	})
}

func TestUsersRepository_UpdateUser(t *testing.T) {
	mock, gormDB := setupSQLMock(t)

	usersRepo := repository.NewRepository(gormDB, nil)

	expectedQueryString := regexp.QuoteMeta(`UPDATE "users" SET "email"=$1,"password"=$2,"updated_at"=$3 WHERE "id" = $4`)

	t.Run("Positive Case", func(t *testing.T) {
		usersId := 1
		users := &domain.Users{
			//ID:        1,
			Email:    "Zinedine",
			Password: "updated_password",
		}
		fmt.Println(users.TableName())
		mock.ExpectExec(expectedQueryString).
			WithArgs(users.Email, users.Password, sqlmock.AnyArg(), usersId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := usersRepo.Update(context.Background(), gormDB, usersId, users)

		require.NoError(t, err)
	})

	t.Run("Negative Case", func(t *testing.T) {
		usersID := 1
		users := &domain.Users{
			Email:    "Zinedine",
			Password: "updated_password",
		}

		mock.ExpectExec(expectedQueryString).
			WithArgs(users.Email, users.Password, sqlmock.AnyArg(), usersID).
			WillReturnError(errors.New("db error"))

		err := usersRepo.Update(context.Background(), gormDB, usersID, users)

		require.Error(t, err)
		require.EqualError(t, err, "db error")
	})
}

func TestUsersRepository_FindUsers(t *testing.T) {
	mock, gormDB := setupSQLMock(t)
	usersRepo := repository.NewRepository(gormDB, nil)
	countQuery := regexp.QuoteMeta(`SELECT count(*) FROM "users"`)
	expectedQueryString := regexp.QuoteMeta(`SELECT "id","email","password","created_at","updated_at" FROM "users" LIMIT 1 OFFSET 9`)
	limit := 1
	page := 10
	t.Run("Positive Case", func(t *testing.T) {

		mock.ExpectQuery(countQuery).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		rows := sqlmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at"}).
			AddRow(1, "Zinedine@email.com", "password1", time.Now(), time.Now()).
			AddRow(2, "Ronaldo@email.com", "password2", time.Now(), time.Now())

		mock.ExpectQuery(expectedQueryString).
			WillReturnRows(rows)

		users, pagination, err := usersRepo.Find(context.Background(), gormDB, limit, page)

		require.NoError(t, err)
		require.NotNil(t, users)
		require.NotNil(t, pagination)
		require.Len(t, *users, 2)
	})

	t.Run("Negative Case - DB Error", func(t *testing.T) {
		mock.ExpectQuery(countQuery).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
		mock.ExpectQuery(expectedQueryString).
			WillReturnError(errors.New("db error"))

		users, pagination, err := usersRepo.Find(context.Background(), gormDB, limit, page)

		require.Error(t, err)
		require.Nil(t, users)
		require.Nil(t, pagination)
		require.EqualError(t, err, "db error")
	})

	t.Run("Negative Case - Record Not Found", func(t *testing.T) {
		mock.ExpectQuery(countQuery).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
		mock.ExpectQuery(expectedQueryString).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at"})).
			WillReturnError(gorm.ErrRecordNotFound)

		_, pagination, err := usersRepo.Find(context.Background(), gormDB, limit, page)
		require.NoError(t, err)
		//require.Nil(t, users)
		require.NotNil(t, pagination)
		//require.Len(t, *users, 0)
	})
}

func TestUsersRepository_DetailUser(t *testing.T) {
	mock, gormDB := setupSQLMock(t)
	usersRepo := repository.NewRepository(gormDB, nil)

	expectedQueryString := regexp.QuoteMeta(`SELECT "id","email","password","created_at","updated_at" FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT 1`)

	t.Run("Positive Case", func(t *testing.T) {
		usersID := 1
		rows := sqlmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at"}).
			AddRow(usersID, "Zinedine@email.com", "password1", time.Now(), time.Now())
		mock.ExpectQuery(expectedQueryString).
			WithArgs(usersID).
			WillReturnRows(rows)

		users, err := usersRepo.Detail(context.Background(), gormDB, usersID)

		require.NoError(t, err)
		require.NotNil(t, users)
		require.Equal(t, usersID, users.ID)
		require.Equal(t, "Zinedine@email.com", users.Email)
	})

	t.Run("Negative Case - DB Error", func(t *testing.T) {
		usersID := 1

		mock.ExpectQuery(expectedQueryString).
			WithArgs(usersID).
			WillReturnError(errors.New("db error"))

		users, err := usersRepo.Detail(context.Background(), gormDB, usersID)

		require.Error(t, err)
		require.Nil(t, users)
		require.EqualError(t, err, "db error")
	})

	t.Run("Negative Case - Record Not Found", func(t *testing.T) {
		usersID := 1

		mock.ExpectQuery(expectedQueryString).
			WithArgs(usersID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at"}))

		users, err := usersRepo.Detail(context.Background(), gormDB, usersID)

		require.NoError(t, err)
		require.Nil(t, users)
	})
}

func TestUsersRepository_AuthUser(t *testing.T) {
	mock, gormDB := setupSQLMock(t)
	usersRepo := repository.NewRepository(gormDB, nil)

	expectedQueryString := regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND password = $2 ORDER BY "users"."id" LIMIT 1`)

	t.Run("Positive Case", func(t *testing.T) {
		users := "Zinedine@email.com"
		password := "password1"
		rows := sqlmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at"}).
			AddRow(1, "Zinedine@email.com", "password1", time.Now(), time.Now())

		mock.ExpectQuery(expectedQueryString).
			WithArgs(users, password).
			WillReturnRows(rows)

		authUsers, err := usersRepo.Auth(context.Background(), gormDB, users, password)

		require.NoError(t, err)
		require.NotNil(t, authUsers)
		require.Equal(t, users, authUsers.Email)
		require.Equal(t, password, authUsers.Password)
	})

	t.Run("Negative Case - DB Error", func(t *testing.T) {
		email := "Zinedine@email.com"
		password := "password1"

		mock.ExpectQuery(expectedQueryString).
			WithArgs(email, password).
			WillReturnError(errors.New("db error"))

		authUsers, err := usersRepo.Auth(context.Background(), gormDB, email, password)

		require.Error(t, err)
		require.Nil(t, authUsers)
		require.EqualError(t, err, "db error")
	})

	t.Run("Negative Case - Record Not Found", func(t *testing.T) {
		email := "Zinedine@email.com"
		password := "password1"

		mock.ExpectQuery(expectedQueryString).
			WithArgs(email, password).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at"}))

		authUsers, err := usersRepo.Auth(context.Background(), gormDB, email, password)

		require.NoError(t, err)
		require.Nil(t, authUsers)
	})
}
