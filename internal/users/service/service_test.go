package service_test

import (
	"boiler-plate/internal/base/app"
	submissionsDomain "boiler-plate/internal/submissions/domain"
	submissionsMock "boiler-plate/internal/submissions/mocks"
	"boiler-plate/internal/users/domain"
	"boiler-plate/internal/users/mocks"
	"boiler-plate/internal/users/service"
	"boiler-plate/pkg/db"
	"boiler-plate/pkg/exception"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)
		assert.NotNil(t, mockUsersService)
	})
}

func TestCreateUsers(t *testing.T) {
	mockAppCtx := &app.Context{}
	errorTemplate := exception.Internal("error inserting users", errors.New("test error"))
	errorCommitTemplate := exception.Internal("commit transaction", errors.New("commit error"))
	t.Run("CreateUsers Failed Repository", func(t *testing.T) {
		// Variabel
		request := &domain.Users{
			Email:    "Zinedine",
			Password: "test",
		}
		// Mock
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()

		mockRepository.On("Create", mockAppCtx, mock.Anything, request).Return(errors.New("test error"))
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		err := mockUsersService.Create(mockAppCtx, request)
		assert.NotNil(t, err)
		assert.Equal(t, errorTemplate, err)
	})

	t.Run("CreateUsers InvalidInput Validator", func(t *testing.T) {

		// Mock dependencies
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()

		// Set up input
		request := &domain.Users{
			Email:    "Z",
			Password: "test",
		}
		mockRepository.On("Create", mockAppCtx, mock.Anything, request).Return(nil)
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		err := validate.Struct(request)
		errorValidator := exception.InvalidArgument(err)
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockUsersService.Create(mockAppCtx, request)
		// Assert the validation result
		if assert.NotNil(t, err) {
			assert.NotNil(t, errService)
			assert.Equal(t, errService, errorValidator)
		}
	})
	t.Run("CreateUsers Commit Error", func(t *testing.T) {
		// Mock dependencies
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()

		// Mock the Begin method to return the mock transaction

		// Set up input
		request := &domain.Users{
			Email:    "Zinedine",
			Password: "test",
		}

		mockSql.ExpectBegin()
		mockSql.ExpectCommit().WillReturnError(errors.New("commit error"))
		mockRepository.On("Create", mockAppCtx, mock.Anything, request).Return(nil)

		// Initialize service with mockDB
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		// Validate the request
		err := validate.Struct(request)
		if err != nil {
			t.Errorf("Validation failed: %v", err)
		} else {
			assert.Nil(t, err, "Validation should pass")
		}

		// Call the function under test
		errService := mockUsersService.Create(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
		assert.Equal(t, errorCommitTemplate, errService)

		// Assert that the expected methods were called
		mockRepository.AssertCalled(t, "Create", mockAppCtx, mock.Anything, request)
	})
	t.Run("CreateUsers ValidInput Validator", func(t *testing.T) {
		// Mock dependencies
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()

		// Mock the Begin method to return the mock transaction

		// Set up input
		request := &domain.Users{
			Email:    "Zinedine",
			Password: "test",
		}

		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		mockRepository.On("Create", mockAppCtx, mock.Anything, request).Return(nil)

		// Initialize service with mockDB
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		// Validate the request
		err := validate.Struct(request)
		if err != nil {
			t.Errorf("Validation failed: %v", err)
		} else {
			assert.Nil(t, err, "Validation should pass")
		}

		// Call the function under test
		errService := mockUsersService.Create(mockAppCtx, request)

		// Assert the result
		assert.Nil(t, errService)

		// Assert that the expected methods were called
		mockRepository.AssertCalled(t, "Create", mockAppCtx, mock.Anything, request)
	})
}

func TestUpdateUsers(t *testing.T) {
	mockAppCtx := &app.Context{}
	errorTemplate := exception.Internal("error updating users", errors.New("test error"))
	errorCommitTemplate := exception.Internal("commit transaction", errors.New("commit error"))
	invalidIDErrorTemplate := exception.PermissionDenied("Input of id must be integer")
	//validationError := validator.New().Struct(&domain.Users{})

	t.Run("UpdateUsers Failed Repository", func(t *testing.T) {
		request := &domain.Users{
			Email:    "Zinedine",
			Password: "test",
		}
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()

		mockRepository.On("Update", mockAppCtx, mock.Anything, 1, request).Return(errors.New("test error"))
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		err := mockUsersService.Update(mockAppCtx, "1", request)
		assert.NotNil(t, err)
		assert.Equal(t, errorTemplate, err)
	})

	t.Run("UpdateUsers InvalidInput Validator", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()

		request := &domain.Users{
			Email:    "Z",
			Password: "test",
		}
		mockRepository.On("Update", mockAppCtx, mock.Anything, 1, request).Return(nil)
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		err := validate.Struct(request)
		errorValidator := exception.InvalidArgument(err)
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockUsersService.Update(mockAppCtx, "1", request)
		if assert.NotNil(t, err) {
			assert.NotNil(t, errService)
			assert.Equal(t, errService, errorValidator)
		}
	})

	t.Run("UpdateUsers Invalid ID", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()

		request := &domain.Users{
			Email:    "Zinedine",
			Password: "test",
		}
		mockRepository.On("Update", mockAppCtx, mock.Anything, mock.Anything, request).Return(nil)
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockUsersService.Update(mockAppCtx, "invalid_id", request)
		assert.NotNil(t, errService)
		assert.Equal(t, invalidIDErrorTemplate, errService)
	})

	t.Run("UpdateUsers Commit Error", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()

		request := &domain.Users{
			Email:    "Zinedine",
			Password: "test",
		}

		mockSql.ExpectBegin()
		mockSql.ExpectCommit().WillReturnError(errors.New("commit error"))
		mockRepository.On("Update", mockAppCtx, mock.Anything, 1, request).Return(nil)
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		err := validate.Struct(request)
		if err != nil {
			t.Errorf("Validation failed: %v", err)
		} else {
			assert.Nil(t, err, "Validation should pass")
		}

		errService := mockUsersService.Update(mockAppCtx, "1", request)
		assert.NotNil(t, errService)
		assert.Equal(t, errorCommitTemplate, errService)
		mockRepository.AssertCalled(t, "Update", mockAppCtx, mock.Anything, 1, request)
	})

	t.Run("UpdateUsers ValidInput Validator", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()

		request := &domain.Users{
			Email:    "Zinedine",
			Password: "test",
		}

		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		mockRepository.On("Update", mockAppCtx, mock.Anything, 1, request).Return(nil)
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		err := validate.Struct(request)
		if err != nil {
			t.Errorf("Validation failed: %v", err)
		} else {
			assert.Nil(t, err, "Validation should pass")
		}

		errService := mockUsersService.Update(mockAppCtx, "1", request)
		assert.Nil(t, errService)
		mockRepository.AssertCalled(t, "Update", mockAppCtx, mock.Anything, 1, request)
	})
}

func TestAuthUsers(t *testing.T) {
	mockAppCtx := &app.Context{}
	authUserErrorTemplate := exception.Internal("error finding users", errors.New("test error"))
	commitErrorTemplate := exception.Internal("commit transaction", errors.New("commit error"))

	t.Run("AuthUsers Failed Repository", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		mockRepository.On("Auth", mockAppCtx, mock.Anything, "test@example.com", "password").Return(nil, errors.New("test error"))
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		_, err := mockUsersService.Auth(mockAppCtx, "test@example.com", "password")
		assert.NotNil(t, err)
		assert.Equal(t, authUserErrorTemplate, err)
	})

	t.Run("AuthUsers Commit Error", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		user := &domain.Users{Email: "test@example.com"}
		mockRepository.On("Auth", mockAppCtx, mock.Anything, "test@example.com", "password").Return(user, nil)
		mockSql.ExpectBegin()
		mockSql.ExpectCommit().WillReturnError(errors.New("commit error"))
		_, err := mockUsersService.Auth(mockAppCtx, "test@example.com", "password")
		assert.NotNil(t, err)
		assert.Equal(t, commitErrorTemplate, err)
	})

	t.Run("AuthUsers Valid Input", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		user := &domain.Users{Email: "test@example.com"}
		mockRepository.On("Auth", mockAppCtx, mock.Anything, "test@example.com", "password").Return(user, nil)
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()

		result, err := mockUsersService.Auth(mockAppCtx, "test@example.com", "password")
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.Email, result.Email)
	})
}

func TestDeleteUsers(t *testing.T) {
	mockAppCtx := &app.Context{}
	invalidIDErrorTemplate := exception.PermissionDenied("Input of id must be integer")
	deleteUserErrorTemplate := exception.Internal("error deleting users", errors.New("test error"))
	commitErrorTemplate := exception.Internal("commit transaction", errors.New("commit error"))

	t.Run("DeleteUsers Invalid ID", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		err := mockUsersService.Delete(mockAppCtx, "invalid_id")
		assert.NotNil(t, err)
		assert.Equal(t, invalidIDErrorTemplate, err)
	})

	t.Run("DeleteUsers Failed Repository", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		mockRepository.On("Delete", mockAppCtx, mock.Anything, 1).Return(errors.New("test error"))
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		err := mockUsersService.Delete(mockAppCtx, "1")
		assert.NotNil(t, err)
		assert.Equal(t, deleteUserErrorTemplate, err)
	})

	t.Run("DeleteUsers Commit Error", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		mockRepository.On("Delete", mockAppCtx, mock.Anything, 1).Return(nil)
		mockSql.ExpectBegin()
		mockSql.ExpectCommit().WillReturnError(errors.New("commit error"))
		err := mockUsersService.Delete(mockAppCtx, "1")
		assert.NotNil(t, err)
		assert.Equal(t, commitErrorTemplate, err)
	})

	t.Run("DeleteUsers Valid Input", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		mockRepository.On("Delete", mockAppCtx, mock.Anything, 1).Return(nil)
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()

		err := mockUsersService.Delete(mockAppCtx, "1")
		assert.Nil(t, err)
	})
}

func TestDetailUser(t *testing.T) {
	mockAppCtx := &app.Context{}
	notFoundErrorTemplate := exception.NotFound("detail not found")
	internalErrorTemplate := exception.Internal("error getting detail users", errors.New("test error"))
	internalErrorSubmissionsTemplate := exception.Internal("error getting detail submissions", errors.New("test error"))
	commitErrorTemplate := exception.Internal("commit transaction", errors.New("commit error"))
	id := "1"
	userID := 1

	t.Run("DetailUser Invalid ID", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result, err := mockUsersService.Detail(mockAppCtx, "invalid_id")
		assert.NotNil(t, err)
		assert.Equal(t, exception.PermissionDenied("Input of id must be integer"), err)
		assert.Nil(t, result)
	})

	t.Run("DetailUser User Not Found", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		mockRepository.On("Detail", mockAppCtx, mock.Anything, userID).Return(nil, nil)
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result, err := mockUsersService.Detail(mockAppCtx, id)
		assert.NotNil(t, err)
		assert.Equal(t, notFoundErrorTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("DetailUser Error Getting User Detail", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		mockRepository.On("Detail", mockAppCtx, mock.Anything, userID).Return(nil, errors.New("test error"))
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result, err := mockUsersService.Detail(mockAppCtx, id)
		assert.NotNil(t, err)
		assert.Equal(t, internalErrorTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("DetailUser Error Getting Submissions", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		user := &domain.UserResponse{ID: userID}
		mockRepository.On("Detail", mockAppCtx, mock.Anything, userID).Return(user, nil)
		mockSubmitRepository.On("DetailByUser", mockAppCtx, mock.Anything, userID).Return(nil, errors.New("test error"))
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result, err := mockUsersService.Detail(mockAppCtx, id)
		assert.NotNil(t, err)
		assert.Equal(t, internalErrorSubmissionsTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("DetailUser Commit Error", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		user := &domain.UserResponse{ID: userID}
		submissions := &submissionsDomain.Submissions{UserId: userID}
		mockRepository.On("Detail", mockAppCtx, mock.Anything, userID).Return(user, nil)
		mockSubmitRepository.On("DetailByUser", mockAppCtx, mock.Anything, userID).Return(submissions, nil)
		mockSql.ExpectBegin()
		mockSql.ExpectCommit().WillReturnError(errors.New("commit error"))
		result, err := mockUsersService.Detail(mockAppCtx, id)
		assert.NotNil(t, err)
		assert.Equal(t, commitErrorTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("DetailUser Valid Input", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		user := &domain.UserResponse{ID: userID}
		submissions := &submissionsDomain.Submissions{UserId: userID}
		mockRepository.On("Detail", mockAppCtx, mock.Anything, userID).Return(user, nil)
		mockSubmitRepository.On("DetailByUser", mockAppCtx, mock.Anything, userID).Return(submissions, nil)
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		result, err := mockUsersService.Detail(mockAppCtx, id)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID, result.ID)
	})
}

func TestFindUsers(t *testing.T) {
	mockAppCtx := &app.Context{}
	internalErrorTemplate := exception.Internal("error geting users", errors.New("test error"))
	internalErrorSubmissionsTemplate := exception.Internal("error getting detail submissions", errors.New("test error"))
	commitErrorTemplate := exception.Internal("commit transaction", errors.New("commit error"))
	permissionLimitDeniedTemplate := exception.PermissionDenied("Input of limit must be integer")
	permissionPageDeniedTemplate := exception.PermissionDenied("Input of page must be integer")

	limit := "10"
	page := "1"
	limitInt := 10
	pageInt := 1

	t.Run("FindUsers Invalid Limit", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result, err := mockUsersService.Find(mockAppCtx, "invalid_limit", page)
		assert.NotNil(t, err)
		assert.Equal(t, permissionLimitDeniedTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("FindUsers Invalid Page", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result, err := mockUsersService.Find(mockAppCtx, limit, "invalid_page")
		assert.NotNil(t, err)
		assert.Equal(t, permissionPageDeniedTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("FindUsers Error Getting Users", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		users := &[]domain.UserResponse{}
		pagination := &db.Paginate{}
		mockRepository.On("Find", mockAppCtx, mock.Anything, limitInt, pageInt).Return(users, pagination, errors.New("test error"))
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result, err := mockUsersService.Find(mockAppCtx, limit, page)
		assert.NotNil(t, err)
		assert.Equal(t, internalErrorTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("FindUsers Error Getting Submissions", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		users := &[]domain.UserResponse{{ID: 1}}
		pagination := &db.Paginate{}
		mockRepository.On("Find", mockAppCtx, mock.Anything, limitInt, pageInt).Return(users, pagination, nil)
		mockSubmitRepository.On("DetailByUser", mockAppCtx, mock.Anything, 1).Return(nil, errors.New("test error"))
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		result, err := mockUsersService.Find(mockAppCtx, limit, page)
		assert.NotNil(t, err)
		assert.Equal(t, internalErrorSubmissionsTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("FindUsers Commit Error", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		users := &[]domain.UserResponse{{ID: 1}}
		submissions := &submissionsDomain.Submissions{UserId: 1}
		pagination := &db.Paginate{}
		mockRepository.On("Find", mockAppCtx, mock.Anything, limitInt, pageInt).Return(users, pagination, nil)
		mockSubmitRepository.On("DetailByUser", mockAppCtx, mock.Anything, 1).Return(submissions, nil)
		mockSql.ExpectBegin()
		mockSql.ExpectCommit().WillReturnError(errors.New("commit error"))
		result, err := mockUsersService.Find(mockAppCtx, limit, page)
		assert.NotNil(t, err)
		assert.Equal(t, commitErrorTemplate, err)
		assert.Nil(t, result)
	})

	t.Run("FindUsers Valid Input", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UsersRepository)
		mockSubmitRepository := new(submissionsMock.SubmissionsRepository)
		validate := validator.New()
		mockUsersService := service.NewService(nil, mockRepository, mockSubmitRepository, gormDB, validate)

		users := &[]domain.UserResponse{{ID: 1}}
		submissions := &submissionsDomain.Submissions{UserId: 1}
		pagination := &db.Paginate{Limit: limitInt, Page: pageInt}
		mockRepository.On("Find", mockAppCtx, mock.Anything, limitInt, pageInt).Return(users, pagination, nil)
		mockSubmitRepository.On("DetailByUser", mockAppCtx, mock.Anything, 1).Return(submissions, nil)
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		result, err := mockUsersService.Find(mockAppCtx, limit, page)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, limitInt, result.Pagination.Limit)
		assert.Equal(t, pageInt, result.Pagination.Page)
		assert.Equal(t, 1, len(result.Data))
		assert.Equal(t, 1, result.Data[0].ID)
	})
}
