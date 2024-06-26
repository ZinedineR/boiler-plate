package handler_test

import (
	"boiler-plate/internal/base/handler"
	"boiler-plate/internal/submissions/domain"
	handler2 "boiler-plate/internal/submissions/handler"
	"boiler-plate/internal/submissions/mocks"
	"boiler-plate/internal/submissions/service"
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
	assert.Equal(t, mockService, httpHandler.SubmissionsService)
}

func TestSubmissionHandler_Create(t *testing.T) {
	// Setup router and mocks

	t.Run("Valid Submission", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockSubmissionsService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:                mockBaseHTTPHandler,
			SubmissionsService: mockSubmissionsService,
		}

		r.POST("/submissions", mockBaseHTTPHandler.GuestRunAction(httpHandler.Create))
		// Prepare valid submission request
		validRequest := domain.SubmissionRequest{
			UserId: 1,
			Answers: []struct {
				QuestionId int    `json:"question_id"`
				Answer     string `json:"answer"`
			}{
				{QuestionId: 1, Answer: "Answer to question 1"},
				{QuestionId: 2, Answer: "Answer to question 2"},
			},
		}
		requestBodyBytes, _ := json.Marshal(validRequest)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/submissions", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Mock service expectation
		mockSubmissionsService.On("Create", mock.Anything, &validRequest).Return(nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body
		expectedBody := `{"data":{"answers":[{"answer":"Answer to question 1","question_id":1},{"answer":"Answer to question 2","question_id":2}],"user_id":1},"message":"success created","status_code":200}`

		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock service method was called with the expected parameters
		mockSubmissionsService.AssertCalled(t, "Create", mock.Anything, &validRequest)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockSubmissionsService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:                mockBaseHTTPHandler,
			SubmissionsService: mockSubmissionsService,
		}

		r.POST("/submissions", mockBaseHTTPHandler.GuestRunAction(httpHandler.Create))
		// Prepare invalid JSON request
		invalidJSON := `{"user_id": 1, "answers": [{ "question_id": 1, "answer": "Answer 1" }]`
		requestBodyBytes := []byte(invalidJSON)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/submissions", bytes.NewBuffer(requestBodyBytes))
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
		expectedBody := `{"message": "error reading request", "status_code": 400}`
		assert.JSONEq(t, expectedBody, w.Body.String())
	})

	t.Run("Error from Service", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockSubmissionsService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:                mockBaseHTTPHandler,
			SubmissionsService: mockSubmissionsService,
		}

		r.POST("/submissions", mockBaseHTTPHandler.GuestRunAction(httpHandler.Create))
		// Prepare valid submission request
		validRequest := domain.SubmissionRequest{
			UserId: 1,
			Answers: []struct {
				QuestionId int    `json:"question_id"`
				Answer     string `json:"answer"`
			}{
				{QuestionId: 1, Answer: "Answer to question 1"},
				{QuestionId: 2, Answer: "Answer to question 2"},
			},
		}
		requestBodyBytes, _ := json.Marshal(validRequest)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/submissions", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Mock service expectation
		errService := exception.Internal("error inserting submission", errors.New("service error"))
		mockSubmissionsService.On("Create", mock.Anything, &validRequest).Return(errService)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Check response body
		expectedBody := `{"status_code": 500, "message": "error inserting submission", "error": "service error"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock service method was called with the expected parameters
		mockSubmissionsService.AssertCalled(t, "Create", mock.Anything, &validRequest)
	})
}

func TestSubmissionHandler_Detail(t *testing.T) {
	t.Run("Valid Detail Request", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockSubmissionsService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:                mockBaseHTTPHandler,
			SubmissionsService: mockSubmissionsService,
		}

		r.GET("/submissions/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.Detail))

		// Prepare valid detail request
		validID := "1"
		mockDetail := &domain.Submissions{
			ID:             1,
			UserId:         1,
			Answers:        []byte(`{"answers": "data"}`),
			RiskScore:      5,
			RiskCategory:   "low",
			RiskDefinition: "Low risk",
		}

		// Mock service expectation
		mockSubmissionsService.On("Detail", mock.Anything, validID).Return(mockDetail, nil)

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/submissions/"+validID, nil)

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body
		expectedBody := `{"data":{"answers":{"answers":"data"},"created_at":null,"id":1,"risk_category":"low","risk_definition":"Low risk","risk_score":5,"updated_at":null,"user_id":1},"message":"success","status_code":200}`

		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock service method was called with the expected parameters
		mockSubmissionsService.AssertCalled(t, "Detail", mock.Anything, validID)
	})

	t.Run("Error from Service", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockSubmissionsService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:                mockBaseHTTPHandler,
			SubmissionsService: mockSubmissionsService,
		}

		r.GET("/submissions/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.Detail))

		// Prepare valid detail request
		validID := "1"
		errService := exception.NotFound("submission not found")
		mockSubmissionsService.On("Detail", mock.Anything, validID).Return(nil, errService)

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/submissions/"+validID, nil)

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Check response body
		expectedBody := `{"status_code": 404, "message": "submission not found"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock service method was called with the expected parameters
		mockSubmissionsService.AssertCalled(t, "Detail", mock.Anything, validID)
	})
}

func TestSubmissionHandler_Find(t *testing.T) {
	t.Run("Successful Find Request", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockSubmissionsService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:                mockBaseHTTPHandler,
			SubmissionsService: mockSubmissionsService,
		}

		r.GET("/submissions", mockBaseHTTPHandler.GuestRunAction(httpHandler.Find))

		// Prepare valid find request
		limitParam := "10"
		pageParam := "1"
		mockPagination := db.NewPaginate(10, 1)
		mockFindResponse := &service.FindResponse{
			Pagination: *mockPagination,
			Data: []domain.Submissions{
				{
					ID:             1,
					UserId:         1,
					RiskScore:      5,
					RiskCategory:   "low",
					RiskDefinition: "Low risk",
				},
				{
					ID:             2,
					UserId:         1,
					RiskScore:      5,
					RiskCategory:   "low",
					RiskDefinition: "Low risk",
				},
			},
		}

		// Mock service expectation
		mockSubmissionsService.On("Find", mock.Anything, limitParam, pageParam).Return(mockFindResponse, nil)

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/submissions?pageSize="+limitParam+"&page="+pageParam, nil)

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body (validate pagination and data)
		expectedBody := `{"data":{"data":[{"created_at":null,"id":1,"risk_category":"low","risk_definition":"Low risk","risk_score":5,"updated_at":null,"user_id":1},{"created_at":null,"id":2,"risk_category":"low","risk_definition":"Low risk","risk_score":5,"updated_at":null,"user_id":1}],"pagination":{"limit":10,"page":1}},"message":"success","status_code":200}`

		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock service method was called with the expected parameters
		mockSubmissionsService.AssertCalled(t, "Find", mock.Anything, limitParam, pageParam)
	})

	t.Run("Error from Service", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockSubmissionsService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:                mockBaseHTTPHandler,
			SubmissionsService: mockSubmissionsService,
		}

		r.GET("/submissions", mockBaseHTTPHandler.GuestRunAction(httpHandler.Find))

		// Prepare valid find request
		limitParam := "10"
		pageParam := "1"
		errService := exception.Internal("error finding submissions", errors.New("service error"))
		mockSubmissionsService.On("Find", mock.Anything, limitParam, pageParam).Return(nil, errService)

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/submissions?pageSize="+limitParam+"&page="+pageParam, nil)

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Check response body
		expectedBody := `{"status_code":500,"message":"error finding submissions","error":"service error"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock service method was called with the expected parameters
		mockSubmissionsService.AssertCalled(t, "Find", mock.Anything, limitParam, pageParam)
	})
}

func TestSubmissionHandler_FindByUser(t *testing.T) {
	t.Run("Successful FindByUser Request", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockSubmissionsService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:                mockBaseHTTPHandler,
			SubmissionsService: mockSubmissionsService,
		}

		r.GET("/submissions/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.FindByUser))

		// Prepare valid FindByUser request
		limitParam := "10"
		pageParam := "1"
		userID := "1"
		mockPagination := db.NewPaginate(10, 1)
		mockFindByUserResponse := &service.FindByUserResponse{
			UserId:     1,
			Pagination: *mockPagination,
			Data: []domain.Submissions{
				{
					ID:             1,
					UserId:         1,
					RiskScore:      5,
					RiskCategory:   "low",
					RiskDefinition: "Low risk",
				},
				{
					ID:             2,
					UserId:         1,
					RiskScore:      5,
					RiskCategory:   "low",
					RiskDefinition: "Low risk",
				},
			},
		}

		// Mock service expectation
		mockSubmissionsService.On("FindByUser", mock.Anything, limitParam, pageParam, userID).Return(mockFindByUserResponse, nil)

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/submissions/"+userID+"?pageSize="+limitParam+"&page="+pageParam, nil)

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req
		ginCtx.Params = gin.Params{gin.Param{Key: "id", Value: userID}}

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body (validate pagination and data)
		expectedBody := `{"data":{"data":[{"created_at":null,"id":1,"risk_category":"low","risk_definition":"Low risk","risk_score":5,"updated_at":null,"user_id":1},{"created_at":null,"id":2,"risk_category":"low","risk_definition":"Low risk","risk_score":5,"updated_at":null,"user_id":1}],"pagination":{"limit":10,"page":1},"user_id":1},"message":"success","status_code":200}`

		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock service method was called with the expected parameters
		mockSubmissionsService.AssertCalled(t, "FindByUser", mock.Anything, limitParam, pageParam, userID)
	})

	t.Run("Error from FindByUser Service", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockSubmissionsService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:                mockBaseHTTPHandler,
			SubmissionsService: mockSubmissionsService,
		}

		r.GET("/submissions/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.FindByUser))

		// Prepare valid FindByUser request
		limitParam := "10"
		pageParam := "1"
		userID := "1"
		errService := exception.Internal("error finding submissions by user", errors.New("service error"))
		mockSubmissionsService.On("FindByUser", mock.Anything, limitParam, pageParam, userID).Return(nil, errService)

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/submissions/"+userID+"?pageSize="+limitParam+"&page="+pageParam, nil)

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req
		ginCtx.Params = gin.Params{gin.Param{Key: "id", Value: userID}}

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Check response body
		expectedBody := `{"status_code":500,"message":"error finding submissions by user","error":"service error"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock service method was called with the expected parameters
		mockSubmissionsService.AssertCalled(t, "FindByUser", mock.Anything, limitParam, pageParam, userID)
	})
}

func TestSubmissionHandler_Delete(t *testing.T) {
	t.Run("Successful Delete Request", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockSubmissionsService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:                mockBaseHTTPHandler,
			SubmissionsService: mockSubmissionsService,
		}

		r.DELETE("/submissions/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.Delete))

		// Mock service expectation
		id := "1"
		mockSubmissionsService.On("Delete", mock.Anything, id).Return(nil)

		// Create HTTP DELETE request
		req, _ := http.NewRequest("DELETE", "/submissions/"+id, nil)

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req
		ginCtx.Params = gin.Params{gin.Param{Key: "id", Value: id}}

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body
		expectedBody := `{"status_code":200,"message":"success delete id: 1"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock service method was called with the expected parameters
		mockSubmissionsService.AssertCalled(t, "Delete", mock.Anything, id)
	})

	t.Run("Error from Delete Service", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		dbMock, _ := gorm.Open(nil, nil)
		mockSubmissionsService := new(mocks.Service)
		mockBaseHTTPHandler := handler.NewBaseHTTPHandler(dbMock, nil, nil, nil)
		httpHandler := handler2.HTTPHandler{
			App:                mockBaseHTTPHandler,
			SubmissionsService: mockSubmissionsService,
		}

		r.DELETE("/submissions/:id", mockBaseHTTPHandler.GuestRunAction(httpHandler.Delete))

		// Mock service expectation
		id := "1"
		errService := exception.Internal("error deleting submission", errors.New("service error"))
		mockSubmissionsService.On("Delete", mock.Anything, id).Return(errService)

		// Create HTTP DELETE request
		req, _ := http.NewRequest("DELETE", "/submissions/"+id, nil)

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req
		ginCtx.Params = gin.Params{gin.Param{Key: "id", Value: id}}

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Check response body
		expectedBody := `{"status_code":500,"message":"error deleting submission","error":"service error"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		// Assert that the mock service method was called with the expected parameters
		mockSubmissionsService.AssertCalled(t, "Delete", mock.Anything, id)
	})
}
