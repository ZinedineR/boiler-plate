package service_test

import (
	"boiler-plate/internal/base/app"
	"boiler-plate/internal/submissions/domain"
	"boiler-plate/internal/submissions/mocks"
	"boiler-plate/internal/submissions/service"
	"boiler-plate/pkg/db"
	"boiler-plate/pkg/exception"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
	"testing"
)

func setupSQLMock(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	// Setup SQL mock
	db, mockSql, err := sqlmock.New()
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
	return mockSql, gormDB
}

func TestNewService(t *testing.T) {
	t.Run("test create NewService service", func(t *testing.T) {
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, gormDB, validate)
		assert.NotNil(t, mockUsersService)
	})
}

func TestCreateSubmission(t *testing.T) {
	mockAppCtx := &app.Context{}
	//invalidArgumentErrorTemplate := exception.InvalidArgument(errors.New("validation error"))
	//internalMarshallingErrorTemplate := exception.Internal("error marshalling old value", errors.New("json error"))
	internalInsertErrorTemplate := exception.Internal("error inserting submissions", errors.New("insert error"))
	commitErrorTemplate := exception.Internal("commit transaction", errors.New("commit error"))

	validRequest := &domain.SubmissionRequest{
		UserId: 1,
		Answers: []struct {
			QuestionId int    `json:"question_id"`
			Answer     string `json:"answer"`
		}{
			{QuestionId: 1, Answer: "A"},
			{QuestionId: 2, Answer: "B"},
		},
	}

	t.Run("CreateSubmission Invalid Request", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()

		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		invalidRequest := &domain.SubmissionRequest{}
		err := validate.Struct(invalidRequest)
		errorValidator := exception.InvalidArgument(err)
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Create(mockAppCtx, invalidRequest)
		assert.NotNil(t, errService)
		assert.Equal(t, errorValidator, errService)
	})

	//t.Run("CreateSubmission Marshalling Error", func(t *testing.T) {
	//	mockSql, gormDB := setupSQLMock(t)
	//	mockRepository := new(mocks.SubmissionsRepository)
	//	validate := validator.New()
	//
	//	request := validRequest
	//	mockRepository.On("Create", mockAppCtx, mock.Anything, mock.Anything).Return(nil)
	//	mockService := service.NewService(nil, mockRepository, gormDB, validate)
	//	mockSql.ExpectBegin()
	//	mockSql.ExpectRollback()
	//	// Simulate marshalling error by setting an invalid value in Answers
	//	request.Answers = []struct {
	//		QuestionId int    `json:"question_id"`
	//		Answer     string `json:"answer"`
	//	}{{QuestionId: 1, Answer: string([]byte{0xff})}} // Invalid UTF-8 sequence
	//	err := mockService.Create(mockAppCtx, request)
	//	assert.NotNil(t, err)
	//	assert.Equal(t, internalMarshallingErrorTemplate, err)
	//})

	t.Run("CreateSubmission Insert Error", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		mockSql.ExpectBegin()
		mockRepository.On("Create", mockAppCtx, mock.Anything, mock.Anything).Return(errors.New("insert error"))
		mockSql.ExpectRollback()
		err := mockService.Create(mockAppCtx, validRequest)
		assert.NotNil(t, err)
		assert.Equal(t, internalInsertErrorTemplate, err)
	})

	t.Run("CreateSubmission Commit Error", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		mockSql.ExpectBegin()
		mockRepository.On("Create", mockAppCtx, mock.Anything, mock.Anything).Return(nil)
		mockSql.ExpectCommit().WillReturnError(errors.New("commit error"))
		err := mockService.Create(mockAppCtx, validRequest)
		assert.NotNil(t, err)
		assert.Equal(t, commitErrorTemplate, err)
	})

	t.Run("CreateSubmission Valid Input", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		mockSql.ExpectBegin()
		mockRepository.On("Create", mockAppCtx, mock.Anything, mock.Anything).Return(nil)
		mockSql.ExpectCommit()
		err := mockService.Create(mockAppCtx, validRequest)
		assert.Nil(t, err)
	})
}

func TestDeleteSubmission(t *testing.T) {
	mockAppCtx := &app.Context{}
	permissionDeniedErrorTemplate := exception.PermissionDenied("Input of id must be integer")
	internalErrorTemplate := exception.Internal("error deleting submissions", errors.New("test error"))
	commitErrorTemplate := exception.Internal("commit transaction", errors.New("commit error"))
	id := "1"

	t.Run("DeleteSubmission Invalid ID", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result := mockService.Delete(mockAppCtx, "invalid_id")
		assert.NotNil(t, result)
		assert.Equal(t, permissionDeniedErrorTemplate, result)
	})

	t.Run("DeleteSubmission Error Deleting Submission", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		idInt := 1
		mockRepository.On("Delete", mockAppCtx, mock.Anything, idInt).Return(errors.New("test error"))
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result := mockService.Delete(mockAppCtx, id)
		assert.NotNil(t, result)
		assert.Equal(t, internalErrorTemplate, result)
	})

	t.Run("DeleteSubmission Commit Error", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		idInt := 1
		mockRepository.On("Delete", mockAppCtx, mock.Anything, idInt).Return(nil)
		mockSql.ExpectBegin()
		mockSql.ExpectCommit().WillReturnError(errors.New("commit error"))
		result := mockService.Delete(mockAppCtx, id)
		assert.NotNil(t, result)
		assert.Equal(t, commitErrorTemplate, result)
	})

	t.Run("DeleteSubmission Valid Input", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		idInt := 1
		mockRepository.On("Delete", mockAppCtx, mock.Anything, idInt).Return(nil)
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		result := mockService.Delete(mockAppCtx, id)
		assert.Nil(t, result)
	})
}

func TestDetailSubmission(t *testing.T) {
	mockAppCtx := &app.Context{}
	permissionDeniedErrorTemplate := exception.PermissionDenied("Input of id must be integer")
	internalErrorTemplate := exception.Internal("error getting detail submissions", errors.New("test error"))
	commitErrorTemplate := exception.Internal("commit transaction", errors.New("commit error"))
	id := "1"

	t.Run("DetailSubmission Invalid ID", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result, err := mockService.Detail(mockAppCtx, "invalid_id")
		assert.NotNil(t, err)
		assert.Equal(t, permissionDeniedErrorTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("DetailSubmission Error Getting Detail", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		idInt := 1
		mockRepository.On("Detail", mockAppCtx, mock.Anything, idInt).Return(nil, errors.New("test error"))
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result, err := mockService.Detail(mockAppCtx, id)
		assert.NotNil(t, err)
		assert.Equal(t, internalErrorTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("DetailSubmission Commit Error", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		idInt := 1
		submission := &domain.Submissions{ID: idInt}
		mockRepository.On("Detail", mockAppCtx, mock.Anything, idInt).Return(submission, nil)
		mockSql.ExpectBegin()
		mockSql.ExpectCommit().WillReturnError(errors.New("commit error"))
		result, err := mockService.Detail(mockAppCtx, id)
		assert.NotNil(t, err)
		assert.Equal(t, commitErrorTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("DetailSubmission Valid Input", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		idInt := 1
		submission := &domain.Submissions{ID: idInt}
		mockRepository.On("Detail", mockAppCtx, mock.Anything, idInt).Return(submission, nil)
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		result, err := mockService.Detail(mockAppCtx, id)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, idInt, result.ID)
	})
}

func TestFindSubmissions(t *testing.T) {
	mockAppCtx := &app.Context{}
	permissionDeniedErrorTemplate := exception.PermissionDenied("Input of limit must be integer")
	permissionPageDeniedTemplate := exception.PermissionDenied("Input of page must be integer")
	internalErrorTemplate := exception.Internal("error getting submissions", errors.New("test error"))
	commitErrorTemplate := exception.Internal("commit transaction", errors.New("commit error"))
	limit := "10"
	page := "1"

	t.Run("FindSubmissions Invalid Limit", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result, err := mockService.Find(mockAppCtx, "invalid_limit", page)
		assert.NotNil(t, err)
		assert.Equal(t, permissionDeniedErrorTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("FindSubmissions Invalid Page", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result, err := mockService.Find(mockAppCtx, limit, "invalid_page")
		assert.NotNil(t, err)
		assert.Equal(t, permissionPageDeniedTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("FindSubmissions Error Getting Submissions", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		limitInt, _ := strconv.Atoi(limit)
		pageInt, _ := strconv.Atoi(page)
		mockRepository.On("Find", mockAppCtx, mock.Anything, limitInt, pageInt).Return(nil, nil, errors.New("test error"))
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result, err := mockService.Find(mockAppCtx, limit, page)
		assert.NotNil(t, err)
		assert.Equal(t, internalErrorTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("FindSubmissions Commit Error", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		limitInt, _ := strconv.Atoi(limit)
		pageInt, _ := strconv.Atoi(page)
		submissions := []domain.Submissions{{ID: 1}, {ID: 2}}
		pagination := &db.Paginate{Limit: limitInt, Page: pageInt}
		mockRepository.On("Find", mockAppCtx, mock.Anything, limitInt, pageInt).Return(&submissions, pagination, nil)
		mockSql.ExpectBegin()
		mockSql.ExpectCommit().WillReturnError(errors.New("commit error"))
		result, err := mockService.Find(mockAppCtx, limit, page)
		assert.NotNil(t, err)
		assert.Equal(t, commitErrorTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("FindSubmissions Valid Input", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		limitInt, _ := strconv.Atoi(limit)
		pageInt, _ := strconv.Atoi(page)
		submissions := []domain.Submissions{
			{ID: 1}, {ID: 2},
		}
		pagination := &db.Paginate{Limit: limitInt, Page: pageInt}
		mockRepository.On("Find", mockAppCtx, mock.Anything, limitInt, pageInt).Return(&submissions, pagination, nil)
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		result, err := mockService.Find(mockAppCtx, limit, page)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, submissions, result.Data)
		assert.Equal(t, *pagination, result.Pagination)
	})
}

func TestFindByUserSubmissions(t *testing.T) {
	mockAppCtx := &app.Context{}
	permissionDeniedErrorTemplate := exception.PermissionDenied("Input of id must be integer")
	permissionLimitDeniedTemplate := exception.PermissionDenied("Input of limit must be integer")
	permissionPageDeniedTemplate := exception.PermissionDenied("Input of page must be integer")
	internalErrorTemplate := exception.Internal("error getting submissions", errors.New("test error"))
	commitErrorTemplate := exception.Internal("commit transaction", errors.New("commit error"))
	limit := "10"
	page := "1"
	userID := "1"

	t.Run("FindByUser Invalid User ID", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result, err := mockService.FindByUser(mockAppCtx, limit, page, "invalid_userid")
		assert.NotNil(t, err)
		assert.Equal(t, permissionDeniedErrorTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("FindByUser Invalid Limit", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result, err := mockService.FindByUser(mockAppCtx, "invalid_limit", page, userID)
		assert.NotNil(t, err)
		assert.Equal(t, permissionLimitDeniedTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("FindByUser Invalid Page", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result, err := mockService.FindByUser(mockAppCtx, limit, "invalid_page", userID)
		assert.NotNil(t, err)
		assert.Equal(t, permissionPageDeniedTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("FindByUser Error Getting Submissions", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		userIDInt, _ := strconv.Atoi(userID)
		limitInt, _ := strconv.Atoi(limit)
		pageInt, _ := strconv.Atoi(page)
		mockRepository.On("FindByUser", mockAppCtx, mock.Anything, limitInt, pageInt, userIDInt).Return(nil, nil, errors.New("test error"))
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result, err := mockService.FindByUser(mockAppCtx, limit, page, userID)
		assert.NotNil(t, err)
		assert.Equal(t, internalErrorTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("FindByUser Commit Error", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		userIDInt, _ := strconv.Atoi(userID)
		limitInt, _ := strconv.Atoi(limit)
		pageInt, _ := strconv.Atoi(page)
		submissions := []domain.Submissions{{ID: 1}, {ID: 2}}
		pagination := &db.Paginate{Limit: limitInt, Page: pageInt}
		mockRepository.On("FindByUser", mockAppCtx, mock.Anything, limitInt, pageInt, userIDInt).Return(&submissions, pagination, nil)
		mockSql.ExpectBegin()
		mockSql.ExpectCommit().WillReturnError(errors.New("commit error"))
		result, err := mockService.FindByUser(mockAppCtx, limit, page, userID)
		assert.NotNil(t, err)
		assert.Equal(t, commitErrorTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("FindByUser Valid Input", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.SubmissionsRepository)
		validate := validator.New()
		mockService := service.NewService(nil, mockRepository, gormDB, validate)

		userIDInt, _ := strconv.Atoi(userID)
		limitInt, _ := strconv.Atoi(limit)
		pageInt, _ := strconv.Atoi(page)
		submissions := []domain.Submissions{{ID: 1}, {ID: 2}}
		pagination := &db.Paginate{Limit: limitInt, Page: pageInt}
		mockRepository.On("FindByUser", mockAppCtx, mock.Anything, limitInt, pageInt, userIDInt).Return(&submissions, pagination, nil)
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		result, err := mockService.FindByUser(mockAppCtx, limit, page, userID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userIDInt, result.UserId)
		assert.Equal(t, submissions, result.Data)
		assert.Equal(t, *pagination, result.Pagination)
	})
}
