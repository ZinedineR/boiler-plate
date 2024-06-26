package repository_test

import (
	"boiler-plate/internal/submissions/domain"
	"boiler-plate/internal/submissions/repository"
	userDomain "boiler-plate/internal/users/domain"
	"context"
	"encoding/json"
	"errors"
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

func TestSubmissionsRepository_CreateSubmission(t *testing.T) {
	mock, gormDB := setupSQLMock(t)
	submissionsRepo := repository.NewRepository(gormDB, nil)

	expectedQueryString := regexp.QuoteMeta(`INSERT INTO "submissions" ("user_id","answers","risk_score","risk_category","risk_definition","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)

	t.Run("Positive Case", func(t *testing.T) {
		answers := json.RawMessage(`[{"answer": "Pertumbuhan kekayaan untuk jangka panjang", "user_id": 1, "question_id": 1}, {"answer": "≥ 10 tahun", "user_id": 1, "question_id": 2}]`)
		submission := &domain.Submissions{
			UserId:         1,
			Answers:        answers,
			RiskScore:      5,
			RiskCategory:   "Medium",
			RiskDefinition: "Moderate risk",
		}

		mock.ExpectQuery(expectedQueryString).
			WithArgs(submission.UserId, submission.Answers, submission.RiskScore, submission.RiskCategory, submission.RiskDefinition, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err := submissionsRepo.Create(context.Background(), gormDB, submission)

		require.NoError(t, err)
		require.NotNil(t, submission.ID)
	})

	t.Run("Negative Case", func(t *testing.T) {
		answers := json.RawMessage(`[{"answer": "Pertumbuhan kekayaan untuk jangka panjang", "user_id": 1, "question_id": 1}, {"answer": "≥ 10 tahun", "user_id": 1, "question_id": 2}]`)
		submission := &domain.Submissions{
			UserId:         1,
			Answers:        answers,
			RiskScore:      5,
			RiskCategory:   "Medium",
			RiskDefinition: "Moderate risk",
		}

		mock.ExpectQuery(expectedQueryString).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := submissionsRepo.Create(context.Background(), gormDB, submission)

		require.Error(t, err)
	})
}

func TestSubmissionsRepository_UpdateSubmission(t *testing.T) {
	mock, gormDB := setupSQLMock(t)
	submissionsRepo := repository.NewRepository(gormDB, nil)

	expectedQueryString := regexp.QuoteMeta(`UPDATE "submissions" SET "id"=$1,"user_id"=$2,"answers"=$3,"risk_score"=$4,"risk_category"=$5,"risk_definition"=$6,"updated_at"=$7 WHERE "id" = $8`)
	idSubmissions := 1
	t.Run("Positive Case", func(t *testing.T) {
		answers := json.RawMessage(`[{"answer": "Pertumbuhan kekayaan untuk jangka panjang", "user_id": 1, "question_id": 1}, {"answer": "≥ 10 tahun", "user_id": 1, "question_id": 2}]`)
		submission := &domain.Submissions{
			ID:             idSubmissions,
			UserId:         1,
			Answers:        answers,
			RiskScore:      5,
			RiskCategory:   "Medium",
			RiskDefinition: "Moderate risk",
		}

		mock.ExpectExec(expectedQueryString).
			WithArgs(submission.ID, submission.UserId, submission.Answers, submission.RiskScore, submission.RiskCategory, submission.RiskDefinition, sqlmock.AnyArg(), idSubmissions).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := submissionsRepo.Update(context.Background(), gormDB, submission.ID, submission)

		require.NoError(t, err)
	})

	t.Run("Negative Case", func(t *testing.T) {
		answers := json.RawMessage(`[{"answer": "Pertumbuhan kekayaan untuk jangka panjang", "user_id": 1, "question_id": 1}, {"answer": "≥ 10 tahun", "user_id": 1, "question_id": 2}]`)
		submission := &domain.Submissions{
			ID:             1,
			UserId:         1,
			Answers:        answers,
			RiskScore:      5,
			RiskCategory:   "Medium",
			RiskDefinition: "Moderate risk",
		}

		mock.ExpectExec(expectedQueryString).
			WillReturnError(errors.New("db error"))

		err := submissionsRepo.Update(context.Background(), gormDB, submission.ID, submission)

		require.Error(t, err)
	})
}

func TestSubmissionsRepository_DetailByUser(t *testing.T) {
	mock, gormDB := setupSQLMock(t)
	submissionsRepo := repository.NewRepository(gormDB, nil)

	expectedQueryString := regexp.QuoteMeta(`SELECT * FROM "submissions" WHERE user_id = $1 ORDER BY created_at DESC,"submissions"."id" LIMIT 1`)

	t.Run("Existing Record", func(t *testing.T) {
		userId := 1
		expectedSubmission := &domain.Submissions{
			ID:             1,
			UserId:         1,
			Answers:        json.RawMessage(`[{"answer": "Test answer"}]`),
			RiskScore:      5,
			RiskCategory:   "Low",
			RiskDefinition: "Low risk",
		}

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "answers", "risk_score", "risk_category", "risk_definition", "created_at", "updated_at",
		}).
			AddRow(expectedSubmission.ID, expectedSubmission.UserId, expectedSubmission.Answers, expectedSubmission.RiskScore, expectedSubmission.RiskCategory, expectedSubmission.RiskDefinition, time.Now(), time.Now())

		mock.ExpectQuery(expectedQueryString).
			WithArgs(userId).
			WillReturnRows(rows)

		result, err := submissionsRepo.DetailByUser(context.Background(), gormDB, 1)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, expectedSubmission.ID, result.ID)
		require.Equal(t, expectedSubmission.UserId, result.UserId)
		// Add more assertions for other fields as needed
	})
	t.Run("Negative Case - DB Error", func(t *testing.T) {
		mock.ExpectQuery(expectedQueryString).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "user_id", "answers", "risk_score", "risk_category", "risk_definition", "created_at",
				"updated_at",
			})).
			WillReturnError(errors.New("db error"))

		submissions, err := submissionsRepo.DetailByUser(context.Background(), gormDB, 1)

		require.Error(t, err)
		require.Nil(t, submissions)
		require.EqualError(t, err, "db error")
	})
	t.Run("Non-existing Record", func(t *testing.T) {
		mock.ExpectQuery(expectedQueryString).
			WithArgs(2).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := submissionsRepo.DetailByUser(context.Background(), gormDB, 2)

		require.NoError(t, err)
		require.Nil(t, result)
	})
}

func TestSubmissionsRepository_FindSubmissions(t *testing.T) {
	mock, gormDB := setupSQLMock(t)
	submissionsRepo := repository.NewRepository(gormDB, nil)
	countQuery := regexp.QuoteMeta(`SELECT count(*) FROM "submissions"`)
	expectedQueryString := regexp.QuoteMeta(`SELECT * FROM "submissions" LIMIT 1 OFFSET 9`)
	limit := 1
	page := 10

	t.Run("Positive Case", func(t *testing.T) {
		mock.ExpectQuery(countQuery).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "answers", "risk_score", "risk_category", "risk_definition", "created_at", "updated_at",
		}).
			AddRow(1, 1, json.RawMessage(`[{"answer": "Test answer"}]`), 5, "Low", "Low risk", time.Now(), time.Now()).
			AddRow(2, 2, json.RawMessage(`[{"answer": "Test answer"}]`), 3, "Medium", "Medium risk", time.Now(), time.Now())

		mock.ExpectQuery(expectedQueryString).
			WillReturnRows(rows)

		submissions, pagination, err := submissionsRepo.Find(context.Background(), gormDB, limit, page)

		require.NoError(t, err)
		require.NotNil(t, submissions)
		require.NotNil(t, pagination)
		require.Len(t, *submissions, 2)
	})

	t.Run("Negative Case - DB Error", func(t *testing.T) {
		mock.ExpectQuery(countQuery).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
		mock.ExpectQuery(expectedQueryString).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "user_id", "answers", "risk_score", "risk_category", "risk_definition", "created_at",
				"updated_at",
			})).
			WillReturnError(errors.New("db error"))

		submissions, pagination, err := submissionsRepo.Find(context.Background(), gormDB, limit, page)

		require.Error(t, err)
		require.Nil(t, submissions)
		require.Nil(t, pagination)
		require.EqualError(t, err, "db error")
	})

	t.Run("Negative Case - Record Not Found", func(t *testing.T) {
		mock.ExpectQuery(countQuery).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
		mock.ExpectQuery(expectedQueryString).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "user_id", "answers", "risk_score", "risk_category", "risk_definition", "created_at",
				"updated_at",
			})).
			WillReturnError(gorm.ErrRecordNotFound)

		_, pagination, err := submissionsRepo.Find(context.Background(), gormDB, limit, page)

		require.NoError(t, err)
		require.NotNil(t, pagination)
	})
}

func TestSubmissionsRepository_Delete(t *testing.T) {
	mock, gormDB := setupSQLMock(t)
	submissionsRepo := repository.NewRepository(gormDB, nil)
	expectedQueryString := regexp.QuoteMeta(`DELETE FROM "submissions" WHERE "submissions"."id" = $1`)

	t.Run("Positive Case", func(t *testing.T) {
		// Mock the deletion query and its expected behavior
		mock.ExpectExec(expectedQueryString).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Call the Delete method
		err := submissionsRepo.Delete(context.Background(), gormDB, 1)

		// Assert no error occurred
		require.NoError(t, err)
	})

	t.Run("Negative Case - DB Error", func(t *testing.T) {
		// Mock the deletion query to return an error
		mock.ExpectExec(expectedQueryString).
			WithArgs(1).
			WillReturnError(errors.New("db error"))

		// Call the Delete method
		err := submissionsRepo.Delete(context.Background(), gormDB, 1)

		// Assert an error occurred
		require.Error(t, err)
		require.EqualError(t, err, "db error")
	})
}

func TestSubmissionsRepository_Detail(t *testing.T) {
	mock, gormDB := setupSQLMock(t)
	submissionsRepo := repository.NewRepository(gormDB, nil)

	expectedQueryString := regexp.QuoteMeta(`SELECT * FROM "submissions" WHERE "submissions"."id" = $1 ORDER BY "submissions"."id" LIMIT 1`)
	expectedPreloadQueryString := regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)

	t.Run("Positive Case", func(t *testing.T) {
		// Mock the main query to return a submission
		submissionID := 1
		expectedSubmission := &domain.Submissions{
			ID:             submissionID,
			UserId:         1,
			Answers:        json.RawMessage(`[{"answer": "Test answer"}]`),
			RiskScore:      5,
			RiskCategory:   "Low",
			RiskDefinition: "Low risk",
			CreatedAt:      &time.Time{},
			UpdatedAt:      &time.Time{},
		}

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "answers", "risk_score", "risk_category", "risk_definition", "created_at", "updated_at",
		}).
			AddRow(expectedSubmission.ID, expectedSubmission.UserId, expectedSubmission.Answers, expectedSubmission.RiskScore, expectedSubmission.RiskCategory, expectedSubmission.RiskDefinition, expectedSubmission.CreatedAt, expectedSubmission.UpdatedAt)

		mock.ExpectQuery(expectedQueryString).
			WithArgs(submissionID).
			WillReturnRows(rows)

		// Mock the preload query to return a user
		expectedUser := &userDomain.Users{
			ID:    1,
			Email: "test@example.com",
		}

		userRows := sqlmock.NewRows([]string{"id", "email"}).
			AddRow(expectedUser.ID, expectedUser.Email)

		mock.ExpectQuery(expectedPreloadQueryString).
			WithArgs(expectedUser.ID).
			WillReturnRows(userRows)

		// Call the Detail method
		result, err := submissionsRepo.Detail(context.Background(), gormDB, submissionID)

		// Assert no error occurred and the returned submission is as expected
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, expectedSubmission.ID, result.ID)
		require.Equal(t, expectedSubmission.UserId, result.UserId)
		require.Equal(t, expectedUser.Email, result.User.Email)
		// Add more assertions for other fields as needed
	})

	t.Run("Negative Case - DB Error", func(t *testing.T) {
		submissionID := 1

		// Mock the main query to return a database error
		mock.ExpectQuery(expectedQueryString).
			WithArgs(submissionID).
			WillReturnError(errors.New("db error"))

		// Call the Detail method
		result, err := submissionsRepo.Detail(context.Background(), gormDB, submissionID)

		// Assert an error occurred and no submission is returned
		require.Error(t, err)
		require.Nil(t, result)
		require.EqualError(t, err, "db error")
	})

	t.Run("Negative Case - Record Not Found", func(t *testing.T) {
		submissionID := 1

		// Mock the main query to return a record not found error
		mock.ExpectQuery(expectedQueryString).
			WithArgs(submissionID).
			WillReturnError(gorm.ErrRecordNotFound)

		// Call the Detail method
		result, err := submissionsRepo.Detail(context.Background(), gormDB, submissionID)

		// Assert no error occurred but no submission is returned
		require.NoError(t, err)
		require.Nil(t, result)
	})
}

func TestSubmissionsRepository_FindByUser(t *testing.T) {
	mock, gormDB := setupSQLMock(t)
	submissionsRepo := repository.NewRepository(gormDB, nil)

	countQuery := regexp.QuoteMeta(`SELECT count(*) FROM "submissions" WHERE user_id = $1`)
	expectedQueryString := regexp.QuoteMeta(`SELECT * FROM "submissions" WHERE user_id = $1 LIMIT 1 OFFSET 9`)
	limit := 1
	page := 10
	userId := 1

	t.Run("Positive Case", func(t *testing.T) {
		// Mock the count query
		mock.ExpectQuery(countQuery).
			WithArgs(userId).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		// Mock the main query to return submissions
		rows := sqlmock.NewRows([]string{
			"id", "user_id", "answers", "risk_score", "risk_category", "risk_definition", "created_at", "updated_at",
		}).
			AddRow(1, userId, json.RawMessage(`[{"answer": "Test answer"}]`), 5, "Low", "Low risk", time.Now(), time.Now()).
			AddRow(2, userId, json.RawMessage(`[{"answer": "Test answer"}]`), 6, "Medium", "Medium risk", time.Now(), time.Now())

		mock.ExpectQuery(expectedQueryString).
			WithArgs(userId).
			WillReturnRows(rows)

		// Call the FindByUser method
		result, pagination, err := submissionsRepo.FindByUser(context.Background(), gormDB, limit, page, userId)

		// Assert no error occurred and the returned submissions are as expected
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, pagination)
		require.Len(t, *result, 2)
	})

	t.Run("Negative Case - DB Error", func(t *testing.T) {
		// Mock the count query
		mock.ExpectQuery(countQuery).
			WithArgs(userId).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		// Mock the main query to return a database error
		mock.ExpectQuery(expectedQueryString).
			WithArgs(userId).
			WillReturnError(errors.New("db error"))

		// Call the FindByUser method
		result, pagination, err := submissionsRepo.FindByUser(context.Background(), gormDB, limit, page, userId)

		// Assert an error occurred and no submissions are returned
		require.Error(t, err)
		require.Nil(t, result)
		require.Nil(t, pagination)
		require.EqualError(t, err, "db error")
	})

	t.Run("Negative Case - Record Not Found", func(t *testing.T) {
		// Mock the count query
		mock.ExpectQuery(countQuery).
			WithArgs(userId).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		// Mock the main query to return a record not found error
		mock.ExpectQuery(expectedQueryString).
			WithArgs(userId).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "user_id", "answers", "risk_score", "risk_category", "risk_definition", "created_at",
				"updated_at",
			})).
			WillReturnError(gorm.ErrRecordNotFound)

		// Call the FindByUser method
		_, pagination, err := submissionsRepo.FindByUser(context.Background(), gormDB, limit, page, userId)

		// Assert no error occurred but no submissions are returned
		require.NoError(t, err)
		require.NotNil(t, pagination)
	})
}
