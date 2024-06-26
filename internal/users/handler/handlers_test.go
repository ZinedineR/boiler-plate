package handler_test

import (
	"boiler-plate/internal/base/handler"
	"boiler-plate/internal/users/domain"
	handler2 "boiler-plate/internal/users/handler"
	"boiler-plate/internal/users/mocks"
	"boiler-plate/internal/users/service"
	"boiler-plate/pkg/db"
	"boiler-plate/pkg/exception"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHTTPHandler_ValidInitialization(t *testing.T) {
	dbMock, _ := gorm.Open(nil, nil)
	mockService := new(mocks.Service)
	mockBaseHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
	httpHandler := handler2.NewHTTPHandler(mockBaseHandler, mockService)

	assert.NotNil(t, httpHandler)
	assert.Equal(t, mockBaseHandler, httpHandler.App)
	assert.Equal(t, mockService, httpHandler.UsersService)
}

func TestPostHandler(t *testing.T) {
	// Setup router

	errService := exception.Internal("error inserting users", errors.New("service error"))

	t.Run("Positive Case", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and UsersService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockUsersService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:          mockBaseHTTPHandler,
			UsersService: mockUsersService,
		}

		r.POST("/users", mockBaseHTTPHandler.GuestRunAction(httpHandler.Create))
		// Prepare request data
		requestBody := &domain.Users{
			Email:    "test_users",
			Password: "test_password",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockUsersService.On("Create", mock.Anything, requestBody).Return(nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body
		expectedBody := `{"status_code": 200, "message": "success created", "data": {"id" :0, "email": "test_users", "password": "test_password", "created_at": null, "updated_at": null}}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockUsersService.AssertCalled(t, "Create", mock.Anything, requestBody)
	})
	t.Run("Error service", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and UsersService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockUsersService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:          mockBaseHTTPHandler,
			UsersService: mockUsersService,
		}

		r.POST("/users", mockBaseHTTPHandler.GuestRunAction(httpHandler.Create))
		// Prepare request data
		requestBody := &domain.Users{
			Email:    "test_users",
			Password: "test_password",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockUsersService.On("Create", mock.Anything, requestBody).Return(errService)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Check response body
		expectedBody := `{"status_code": 500, "message": "error inserting users", "error": "service error"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockUsersService.AssertCalled(t, "Create", mock.Anything, requestBody)
	})
	t.Run("Error binding json", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and UsersService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockUsersService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:          mockBaseHTTPHandler,
			UsersService: mockUsersService,
		}

		r.POST("/users", mockBaseHTTPHandler.GuestRunAction(httpHandler.Create))
		// Prepare request data
		malformedJSON := `{"users": 1, "password": "test_password"`
		requestBodyBytes, _ := json.Marshal(malformedJSON)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Check response body
		expectedBody := `{"message":"error reading request", "status_code":400}`
		assert.JSONEq(t, expectedBody, w.Body.String())

	})
}

func TestUpdateHandler(t *testing.T) {
	// Setup router

	errService := exception.Internal("error updating users", errors.New("service error"))

	t.Run("Positive Case", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and UsersService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockUsersService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:          mockBaseHTTPHandler,
			UsersService: mockUsersService,
		}

		r.PUT("/users/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.Update))
		// Prepare request data
		requestBody := &domain.Users{
			Email:    "updated_users",
			Password: "updated_password",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP PUT request
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockUsersService.On("Update", mock.Anything, "1", requestBody).Return(nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body
		expectedBody := `{"status_code": 200, "message": "success update", "data": {"id" :0, "email": "updated_users", "password": "updated_password", "created_at": null, "updated_at": null}}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockUsersService.AssertCalled(t, "Update", mock.Anything, "1", requestBody)
	})
	t.Run("Error service", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and UsersService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockUsersService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:          mockBaseHTTPHandler,
			UsersService: mockUsersService,
		}

		r.PUT("/users/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.Update))
		// Prepare request data
		requestBody := &domain.Users{
			Email:    "updated_users",
			Password: "updated_password",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP PUT request
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockUsersService.On("Update", mock.Anything, "1", requestBody).Return(errService)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Check response body
		expectedBody := `{"status_code": 500, "message": "error updating users", "error": "service error"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockUsersService.AssertCalled(t, "Update", mock.Anything, "1", requestBody)
	})
	t.Run("Error binding json", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and UsersService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockUsersService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:          mockBaseHTTPHandler,
			UsersService: mockUsersService,
		}

		r.PUT("/users/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.Update))
		// Prepare request data
		malformedJSON := `{"users": 1, "password": "updated_password"`
		requestBodyBytes, _ := json.Marshal(malformedJSON)

		// Create HTTP PUT request
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Check response body
		expectedBody := `{"message":"error reading request", "status_code":400}`
		assert.JSONEq(t, expectedBody, w.Body.String())
	})
}

func TestDeleteHandler(t *testing.T) {
	// Setup router
	errService := exception.Internal("error deleting users", errors.New("service error"))

	t.Run("Positive Case", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and UsersService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockUsersService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:          mockBaseHTTPHandler,
			UsersService: mockUsersService,
		}

		r.DELETE("/users/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.Delete))

		// Create HTTP DELETE request
		req, _ := http.NewRequest("DELETE", "/users/1", nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockUsersService.On("Delete", mock.Anything, "1").Return(nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body
		expectedBody := `{"status_code": 200, "message": "success delete id: 1"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockUsersService.AssertCalled(t, "Delete", mock.Anything, "1")
	})
	t.Run("Error service", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and UsersService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockUsersService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:          mockBaseHTTPHandler,
			UsersService: mockUsersService,
		}

		r.DELETE("/users/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.Delete))

		// Create HTTP DELETE request
		req, _ := http.NewRequest("DELETE", "/users/1", nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockUsersService.On("Delete", mock.Anything, "1").Return(errService)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Check response body
		expectedBody := `{"status_code": 500, "message": "error deleting users", "error": "service error"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockUsersService.AssertCalled(t, "Delete", mock.Anything, "1")
	})
}

func TestDetailHandler(t *testing.T) {
	// Setup router

	errService := exception.Internal("error fetching users detail", errors.New("service error"))

	t.Run("Positive Case", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and UsersService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockUsersService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:          mockBaseHTTPHandler,
			UsersService: mockUsersService,
		}

		r.GET("/users/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.Detail))

		// Prepare mock response data
		users := &domain.UserResponse{
			ID:             1,
			Email:          "test_users",
			Password:       "test_password",
			RiskScore:      28,
			RiskCategory:   "Moderate Risk",
			RiskDefinition: "This user is classified as moderate risk based on their investment profile.",
		}

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/users/1", nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockUsersService.On("Detail", mock.Anything, "1").Return(users, nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body
		expectedBody := `{"status_code": 200, "message": "success", "data": {"email":"test_users", "id":1, "password":"test_password", "risk_category":"Moderate Risk", "risk_definition":"This user is classified as moderate risk based on their investment profile.", "risk_score":28, "created_at": null,  "updated_at": null}}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockUsersService.AssertCalled(t, "Detail", mock.Anything, "1")
	})
	t.Run("Error service", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and UsersService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockUsersService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:          mockBaseHTTPHandler,
			UsersService: mockUsersService,
		}

		r.GET("/users/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.Detail))

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/users/1", nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockUsersService.On("Detail", mock.Anything, "1").Return(nil, errService)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Check response body
		expectedBody := `{"status_code": 500, "message": "error fetching users detail", "error": "service error"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockUsersService.AssertCalled(t, "Detail", mock.Anything, "1")
	})
}

func TestFindHandler(t *testing.T) {
	// Setup router

	errService := exception.Internal("error fetching users", errors.New("service error"))
	limit := "1"
	page := "10"
	t.Run("Positive Case", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and UsersService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockUsersService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:          mockBaseHTTPHandler,
			UsersService: mockUsersService,
		}

		r.GET("/users", mockBaseHTTPHandler.GuestRunAction(httpHandler.Find))

		// Prepare mock response data
		users := &[]domain.UserResponse{
			{
				ID:             1,
				Email:          "test_users_1",
				Password:       "test_password_1",
				RiskScore:      28,
				RiskCategory:   "Moderate Risk",
				RiskDefinition: "This user is classified as moderate risk based on their investment profile.",
			},
			{
				ID:             2,
				Email:          "test_users_2",
				Password:       "test_password_2",
				CreatedAt:      nil,
				RiskScore:      28,
				RiskCategory:   "Moderate Risk",
				RiskDefinition: "This user is classified as moderate risk based on their investment profile.",
			},
		}
		pagination := db.NewPaginate(1, 10)
		expectedResponse := &service.FindResponse{
			Pagination: *pagination,
			Data:       *users,
		}
		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/users?pageSize="+limit+"&page="+page, nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockUsersService.On("Find", mock.Anything, limit, page).Return(expectedResponse, nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body
		//expectedBody := `{"status_code": 200, "message": "success", "data": [{"id": 1, "users": "test_users_1", "password": "test_password_1", "created_at": null}, {"id": 2, "users": "test_users_2", "password": "test_password_2", "created_at": null}]}`
		expectedBody := `{"data": {"data": [{"created_at": null, "email": "test_users_1", "id": 1, "password": "test_password_1", "risk_category": "Moderate Risk", "risk_definition": "This user is classified as moderate risk based on their investment profile.", "risk_score": 28, "updated_at": null}, {"created_at": null, "email": "test_users_2", "id": 2, "password": "test_password_2", "risk_category": "Moderate Risk", "risk_definition": "This user is classified as moderate risk based on their investment profile.", "risk_score": 28, "updated_at": null}], "pagination": {"limit": 1, "page": 10}}, "message": "success", "status_code": 200}`

		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockUsersService.AssertCalled(t, "Find", mock.Anything, limit, page)
	})
	t.Run("Error service", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and UsersService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockUsersService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:          mockBaseHTTPHandler,
			UsersService: mockUsersService,
		}

		r.GET("/users", mockBaseHTTPHandler.GuestRunAction(httpHandler.Find))

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/users?pageSize="+limit+"&page="+page, nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockUsersService.On("Find", mock.Anything, limit, page).Return(nil, errService)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Check response body
		expectedBody := `{"status_code": 500, "message": "error fetching users", "error": "service error"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockUsersService.AssertCalled(t, "Find", mock.Anything, limit, page)
	})
}
