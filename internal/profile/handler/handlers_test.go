package handler_test

import (
	"boiler-plate/internal/base/handler"
	"boiler-plate/internal/profile/domain"
	handler2 "boiler-plate/internal/profile/handler"
	"boiler-plate/internal/profile/mocks"
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

func TestPostHandler(t *testing.T) {
	// Setup router

	errService := exception.Internal("error inserting profile", errors.New("service error"))

	t.Run("Positive Case", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and ProfileService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockProfileService := new(mocks.ProfileService)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:            mockBaseHTTPHandler,
			ProfileService: mockProfileService,
		}

		r.POST("/profile", mockBaseHTTPHandler.GuestRunAction(httpHandler.Create))
		// Prepare request data
		requestBody := &domain.Profile{
			Profile:  "test_profile",
			Password: "test_password",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/profile", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockProfileService.On("Create", mock.Anything, requestBody).Return(nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body
		expectedBody := `{"status_code": 200, "message": "success created", "data": {"id" :0, "profile": "test_profile", "password": "test_password", "created_at": null}}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockProfileService.AssertCalled(t, "Create", mock.Anything, requestBody)
	})
	t.Run("Error service", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and ProfileService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockProfileService := new(mocks.ProfileService)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:            mockBaseHTTPHandler,
			ProfileService: mockProfileService,
		}

		r.POST("/profile", mockBaseHTTPHandler.GuestRunAction(httpHandler.Create))
		// Prepare request data
		requestBody := &domain.Profile{
			Profile:  "test_profile",
			Password: "test_password",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/profile", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockProfileService.On("Create", mock.Anything, requestBody).Return(errService)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Check response body
		expectedBody := `{"status_code": 500, "message": "error inserting profile", "error": "service error"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockProfileService.AssertCalled(t, "Create", mock.Anything, requestBody)
	})
	t.Run("Error binding json", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and ProfileService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockProfileService := new(mocks.ProfileService)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:            mockBaseHTTPHandler,
			ProfileService: mockProfileService,
		}

		r.POST("/profile", mockBaseHTTPHandler.GuestRunAction(httpHandler.Create))
		// Prepare request data
		malformedJSON := `{"profile": 1, "password": "test_password"`
		requestBodyBytes, _ := json.Marshal(malformedJSON)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/profile", bytes.NewBuffer(requestBodyBytes))
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

	errService := exception.Internal("error updating profile", errors.New("service error"))

	t.Run("Positive Case", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and ProfileService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockProfileService := new(mocks.ProfileService)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:            mockBaseHTTPHandler,
			ProfileService: mockProfileService,
		}

		r.PUT("/profile/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.Update))
		// Prepare request data
		requestBody := &domain.Profile{
			Profile:  "updated_profile",
			Password: "updated_password",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP PUT request
		req, _ := http.NewRequest("PUT", "/profile/1", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockProfileService.On("Update", mock.Anything, "1", requestBody).Return(nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body
		expectedBody := `{"status_code": 200, "message": "success update", "data": {"id" :0, "profile": "updated_profile", "password": "updated_password", "created_at": null}}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockProfileService.AssertCalled(t, "Update", mock.Anything, "1", requestBody)
	})
	t.Run("Error service", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and ProfileService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockProfileService := new(mocks.ProfileService)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:            mockBaseHTTPHandler,
			ProfileService: mockProfileService,
		}

		r.PUT("/profile/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.Update))
		// Prepare request data
		requestBody := &domain.Profile{
			Profile:  "updated_profile",
			Password: "updated_password",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP PUT request
		req, _ := http.NewRequest("PUT", "/profile/1", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockProfileService.On("Update", mock.Anything, "1", requestBody).Return(errService)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Check response body
		expectedBody := `{"status_code": 500, "message": "error updating profile", "error": "service error"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockProfileService.AssertCalled(t, "Update", mock.Anything, "1", requestBody)
	})
	t.Run("Error binding json", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and ProfileService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockProfileService := new(mocks.ProfileService)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:            mockBaseHTTPHandler,
			ProfileService: mockProfileService,
		}

		r.PUT("/profile/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.Update))
		// Prepare request data
		malformedJSON := `{"profile": 1, "password": "updated_password"`
		requestBodyBytes, _ := json.Marshal(malformedJSON)

		// Create HTTP PUT request
		req, _ := http.NewRequest("PUT", "/profile/1", bytes.NewBuffer(requestBodyBytes))
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

	errService := exception.Internal("error deleting profile", errors.New("service error"))

	t.Run("Positive Case", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and ProfileService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockProfileService := new(mocks.ProfileService)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:            mockBaseHTTPHandler,
			ProfileService: mockProfileService,
		}

		r.DELETE("/profile/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.Delete))

		// Create HTTP DELETE request
		req, _ := http.NewRequest("DELETE", "/profile/1", nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockProfileService.On("Delete", mock.Anything, "1").Return(nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body
		expectedBody := `{"status_code": 200, "message": "success delete id: 1"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockProfileService.AssertCalled(t, "Delete", mock.Anything, "1")
	})
	t.Run("Error service", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and ProfileService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockProfileService := new(mocks.ProfileService)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:            mockBaseHTTPHandler,
			ProfileService: mockProfileService,
		}

		r.DELETE("/profile/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.Delete))

		// Create HTTP DELETE request
		req, _ := http.NewRequest("DELETE", "/profile/1", nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockProfileService.On("Delete", mock.Anything, "1").Return(errService)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Check response body
		expectedBody := `{"status_code": 500, "message": "error deleting profile", "error": "service error"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockProfileService.AssertCalled(t, "Delete", mock.Anything, "1")
	})
}

func TestDetailHandler(t *testing.T) {
	// Setup router

	errService := exception.Internal("error fetching profile detail", errors.New("service error"))

	t.Run("Positive Case", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and ProfileService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockProfileService := new(mocks.ProfileService)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:            mockBaseHTTPHandler,
			ProfileService: mockProfileService,
		}

		r.GET("/profile/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.Detail))

		// Prepare mock response data
		profile := &domain.Profile{
			ID:        1,
			Profile:   "test_profile",
			Password:  "test_password",
			CreatedAt: nil,
		}

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/profile/1", nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockProfileService.On("Detail", mock.Anything, "1").Return(profile, nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body
		expectedBody := `{"status_code": 200, "message": "success", "data": {"id": 1, "profile": "test_profile", "password": "test_password", "created_at": null}}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockProfileService.AssertCalled(t, "Detail", mock.Anything, "1")
	})
	t.Run("Error service", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and ProfileService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockProfileService := new(mocks.ProfileService)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:            mockBaseHTTPHandler,
			ProfileService: mockProfileService,
		}

		r.GET("/profile/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.Detail))

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/profile/1", nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockProfileService.On("Detail", mock.Anything, "1").Return(nil, errService)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Check response body
		expectedBody := `{"status_code": 500, "message": "error fetching profile detail", "error": "service error"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockProfileService.AssertCalled(t, "Detail", mock.Anything, "1")
	})
}

func TestFindHandler(t *testing.T) {
	// Setup router

	errService := exception.Internal("error fetching profiles", errors.New("service error"))

	t.Run("Positive Case", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and ProfileService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockProfileService := new(mocks.ProfileService)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:            mockBaseHTTPHandler,
			ProfileService: mockProfileService,
		}

		r.GET("/profiles", mockBaseHTTPHandler.GuestRunAction(httpHandler.Find))

		// Prepare mock response data
		profiles := &[]domain.Profile{
			{
				ID:        1,
				Profile:   "test_profile_1",
				Password:  "test_password_1",
				CreatedAt: nil,
			},
			{
				ID:        2,
				Profile:   "test_profile_2",
				Password:  "test_password_2",
				CreatedAt: nil,
			},
		}

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/profiles", nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockProfileService.On("Find", mock.Anything).Return(profiles, nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body
		expectedBody := `{"status_code": 200, "message": "success", "data": [{"id": 1, "profile": "test_profile_1", "password": "test_password_1", "created_at": null}, {"id": 2, "profile": "test_profile_2", "password": "test_password_2", "created_at": null}]}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockProfileService.AssertCalled(t, "Find", mock.Anything)
	})
	t.Run("Error service", func(t *testing.T) {
		t.Parallel()
		// Mock GORM DB and ProfileService
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockProfileService := new(mocks.ProfileService)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:            mockBaseHTTPHandler,
			ProfileService: mockProfileService,
		}

		r.GET("/profiles", mockBaseHTTPHandler.GuestRunAction(httpHandler.Find))

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/profiles", nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockProfileService.On("Find", mock.Anything).Return(nil, errService)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Check response body
		expectedBody := `{"status_code": 500, "message": "error fetching profiles", "error": "service error"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock was called with the expected parameters
		mockProfileService.AssertCalled(t, "Find", mock.Anything)
	})
}
