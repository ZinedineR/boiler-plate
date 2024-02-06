package handler

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	baseModel "ms-batch/pkg/db"
	"ms-batch/pkg/httpclient"
	"net/http"
	"os"
	"strings"
	"time"

	"ms-batch/app/appconf"
	"ms-batch/internal/base/app"
	"ms-batch/pkg/server"

	"github.com/gin-gonic/gin"
)

type HandlerFn func(ctx *app.Context) *server.Response
type HandlerFnInterface func(ctx *app.Context) *server.ResponseInterface

type BaseHTTPHandler struct {
	Handlers   interface{}
	DB         *mongo.Client
	AppConfig  *appconf.Config
	BaseModel  *baseModel.MongoDBClientRepository
	HttpClient httpclient.Client
}

func NewBaseHTTPHandler(
	db *mongo.Client,
	appConfig *appconf.Config,
	baseModel *baseModel.MongoDBClientRepository,
	httpClient httpclient.Client,
) *BaseHTTPHandler {
	return &BaseHTTPHandler{
		DB:         db,
		AppConfig:  appConfig,
		BaseModel:  baseModel,
		HttpClient: httpClient,
	}
}

// Handler Basic Method ======================================================================================================
func (h BaseHTTPHandler) AsErrorMessage(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"responseCode":    "500",
		"responseMessage": message,
	})
}

func (h BaseHTTPHandler) AsErrorDefault(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"message": message,
	})
}

func (h BaseHTTPHandler) AsInvalidClientIdError(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"responseCode":    "400",
		"responseMessage": "invalid clientid",
	})
}

func (h BaseHTTPHandler) AsInvalidClientIdAccessTokenError(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"responseCode":    "4010000",
		"responseMessage": "Invalid Client Key",
	})
}

func (h BaseHTTPHandler) AsInvalidPrivateKeyError(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"responseCode":    "4010000",
		"responseMessage": "Invalid Private Key",
	})
}

func (h BaseHTTPHandler) AsDatabaseError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"responseCode":    "500",
		"responseMessage": "Error in database",
	})
}

func (h BaseHTTPHandler) AsNotVerfied(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"responseCode":    "401",
		"responseMessage": "Account still not verified",
	})
}

func (h BaseHTTPHandler) AsDuplicateEmail(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"responseCode":    "401",
		"responseMessage": "Another account with same email already created",
	})
}

func (h BaseHTTPHandler) AsDataNotFound(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"responseCode":    "404",
		"responseMessage": "Data not Found",
	})
}

func (h BaseHTTPHandler) AsJWTExist(ctx *gin.Context, token string) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"responseCode":    "401",
		"responseMessage": "You already login before",
		"token":           token,
	})
}

func (h BaseHTTPHandler) AsPasswordUnmatched(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"responseCode":    "401",
		"responseMessage": "Password Unmatched",
	})
}

func (h BaseHTTPHandler) AsHashError(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"responseCode":    "500",
		"responseMessage": "Error in hashing",
	})
}

func (h BaseHTTPHandler) AsDataUnauthorized(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"responseCode":    "401",
		"responseMessage": message,
	})
}

func (h BaseHTTPHandler) AsEmailNotFound(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"responseCode":    "401",
		"responseMessage": "Can't send email, contact admin for verification",
	})
}

func (h BaseHTTPHandler) AsInvalidPublicKeyError(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"responseCode":    "4010000",
		"responseMessage": "Invalid Public Key",
	})
}

func (h BaseHTTPHandler) AsInvalidSignatureError(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"responseCode":    "4017300",
		"responseMessage": "Invalid Token (B2B)",
	})
}

func (h BaseHTTPHandler) AsRequiredTimeStampError(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"responseCode":    "4000000",
		"responseMessage": "The timestamp field is required.",
	})
}

func (h BaseHTTPHandler) AsInvalidFieldTimeStampError(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"responseCode":    "4000000",
		"responseMessage": "Invalid Field Format Timestamp",
	})
}

func (h BaseHTTPHandler) AsInvalidLengthTimeStampError(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"responseCode":    "4000000",
		"responseMessage": "The field timestamp must be a string or array type with a maximum length of '25'.",
	})
}

func (h BaseHTTPHandler) AsInvalidClientSecretError(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"responseCode":    "4010000",
		"responseMessage": "Invalid Client Secret",
	})
}

func (h BaseHTTPHandler) AsInvalidHttpMethodError(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"responseCode":    "4010000",
		"responseMessage": "http methods is invalid",
	})
}

func (h BaseHTTPHandler) AsInvalidJsonFormat(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"responseCode":    "400",
		"responseMessage": msg,
	})
}

func (h BaseHTTPHandler) AsRequiredClientSecretError(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"responseCode":    "4000000",
		"responseMessage": "The clientSecret field is required.",
	})
}

func (h BaseHTTPHandler) AsRequiredClientIdError(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"responseCode":    "4000000",
		"responseMessage": "The param ID is required.",
	})
}

func (h BaseHTTPHandler) AsRequiredGrantTypeError(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"responseCode":    "4007302",
		"responseMessage": "Bad Request. The grantType field is required.",
	})
}

func (h BaseHTTPHandler) AsRequiredGrantTypeClientCredentialsError(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"responseCode":    "4007300",
		"responseMessage": "grant_type must be set to client_credentials",
	})
}

func (h BaseHTTPHandler) AsRequiredSignatureError(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"responseCode":    "4000000",
		"responseMessage": "The signature field is required.",
	})
}

func (h BaseHTTPHandler) AsRequiredPrivateKeyError(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"responseCode":    "4000000",
		"responseMessage": "The privateKey field is required.",
	})
}

func (h BaseHTTPHandler) AsRequiredContentTypeError(ctx *gin.Context) {
	ctx.JSON(http.StatusUnsupportedMediaType, gin.H{
		"responseCode":    "4000000",
		"responseMessage": "Content Type application/json is required.",
	})
}

func (h BaseHTTPHandler) AsInvalidTokenError(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"responseCode":    "4010001",
		"responseMessage": "Access Token Invalid",
	})
}

func (h BaseHTTPHandler) AsRequiredBearer(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"responseCode":    "4000002",
		"responseMessage": "Bearer authorization is required",
	})
}

func (h BaseHTTPHandler) AsRequiredHttpMethodError(ctx *gin.Context) {
	ctx.JSON(http.StatusUnsupportedMediaType, gin.H{
		"responseCode":    "4000000",
		"responseMessage": "The HttpMethod field is required.",
	})
}

func (h BaseHTTPHandler) AsRequiredEndpoinUrlError(ctx *gin.Context) {
	ctx.JSON(http.StatusUnsupportedMediaType, gin.H{
		"responseCode":    "4000000",
		"responseMessage": "The EndpointUrl field is required.",
	})
}

func (h BaseHTTPHandler) AsRequiredAccessTokenError(ctx *gin.Context) {
	ctx.JSON(http.StatusUnsupportedMediaType, gin.H{
		"responseCode":    "4000000",
		"responseMessage": "The AccessToken field is required.",
	})
}
func (h BaseHTTPHandler) AsRequiredBodyError(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"responseCode":    "4000000",
		"responseMessage": "A non-empty request body is required.",
	})
}

// Data Not Found return AsJsonInterface 404 when data doesn't exist
func (h BaseHTTPHandler) DataNotFound(ctx *app.Context) *server.ResponseInterface {
	type Response struct {
		StatusCode int         `json:"responseCode"`
		Message    interface{} `json:"responseMessage"`
	}
	resp := Response{
		StatusCode: http.StatusNotFound,
		Message:    "Data not found in database.",
	}
	return h.AsJsonInterface(ctx, http.StatusNotFound, resp)

}

// DataReadError return AsJsonInterface error if persist a problem in declared condition
func (h BaseHTTPHandler) DataReadError(ctx *app.Context, code int, description string) *server.ResponseInterface {
	type Response struct {
		StatusCode int         `json:"responseCode"`
		Message    interface{} `json:"responseMessage"`
	}
	resp := Response{
		StatusCode: code,
		Message:    description,
	}
	return h.AsJsonInterface(ctx, code, resp)
}

// AsJson to response custom message: 200, 201 with message (Mobile use 500 error)
func (b BaseHTTPHandler) AsJson(ctx *app.Context, status int, message string, data interface{}) *server.Response {

	return &server.Response{
		Status:       status,
		Message:      message,
		Data:         data,
		ResponseType: server.DefaultResponseType,
	}
}

func (b BaseHTTPHandler) AsJsonInterface(ctx *app.Context, status int, data interface{}) *server.ResponseInterface {

	return &server.ResponseInterface{
		Status: status,
		Data:   data,
	}
}

// ThrowExceptionJson for some exception not handle in Yii2 framework
func (b BaseHTTPHandler) ThrowExceptionJson(ctx *app.Context, status, code int, name, message string) *server.Response {
	return &server.Response{
		Status:  status,
		Message: "",
		Log:     nil,
	}
}

func (b BaseHTTPHandler) UserAuthentication(c *gin.Context) (*app.Context, error) {
	return app.NewContext(c, b.AppConfig), nil
}

func (b BaseHTTPHandler) UserRunAction(handler HandlerFnInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx, err := b.UserAuthentication(c)
		if err != nil {
			logrus.Errorln(fmt.Sprintf("REQUEST ID: %s , message: Unauthorized", ctx.APIReqID))
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized",
				"data":    err.Error(),
			})
			return
		}

		defer func() {
			if err0 := recover(); err0 != nil {
				logrus.Errorln(err0)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "Request is halted unexpectedly, please contact the administrator.",
					"data":    nil,
				})
			}
		}()

		// Extract bearer token from request header
		tokenString := strings.Split(ctx.GetHeader("Authorization"), " ")[1:]
		if len(tokenString) == 0 {
			logrus.Errorln(fmt.Sprintf("REQUEST ID: %s , message: Unauthorized", ctx.APIReqID))
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "request does not contain an access token",
			})
			return
		}
		authToken := tokenString[0]

		// Validate token
		token, err := jwt.ParseWithClaims(authToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
		})
		if err != nil || !token.Valid {
			logrus.Errorln(fmt.Sprintf("REQUEST ID: %s , message: Unauthorized", ctx.APIReqID))
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			return
		}

		// Get user data from token
		claims := token.Claims.(jwt.MapClaims)
		ctx.Set("user_id", claims["id"])
		ctx.Set("user_roles", claims["roles"])

		// Execute handler
		resp := handler(ctx)
		httpStatus := resp.Status

		if resp.Data == nil {
			c.Status(httpStatus)
			return
		}
		end := time.Now().Sub(start)
		logrus.Infoln(fmt.Sprintf("REQUEST ID: %s , LATENCY: %vms", ctx.APIReqID, end.Milliseconds()))
		c.JSON(httpStatus, resp.Data)

	}
}

func (b BaseHTTPHandler) GuestAuthentication(c *gin.Context) (*app.Context, error) {
	return app.NewContext(c, b.AppConfig), nil
}

func (b BaseHTTPHandler) GuestRunAction(handler HandlerFnInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx, err := b.GuestAuthentication(c)
		if err != nil {
			logrus.Errorln(fmt.Sprintf("REQUEST ID: %s , message: Unauthorized", ctx.APIReqID))
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized",
				"data":    err.Error(),
			})
			return
		}

		defer func() {
			if err0 := recover(); err0 != nil {
				logrus.Errorln(err0)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "Request is halted unexpectedly, please contact the administrator.",
					"data":    nil,
				})
			}
		}()

		resp := handler(ctx)
		httpStatus := resp.Status

		if resp.Data == nil {
			c.Status(httpStatus)
			return
		}
		end := time.Now().Sub(start)
		logrus.Infoln(fmt.Sprintf("REQUEST ID: %s , LATENCY: %vms", ctx.APIReqID, end.Milliseconds()))
		c.JSON(httpStatus, resp.Data)

	}
}
