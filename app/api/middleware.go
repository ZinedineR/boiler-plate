package api

import (
	"boiler-plate/pkg/getfilter"
	"context"
	"google.golang.org/protobuf/proto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseHeaderFormat() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func responseHeaderMatcher(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	headers := w.Header()
	if location, ok := headers["Grpc-Metadata-Location"]; ok {
		w.Header().Set("Location", location[0])
		w.WriteHeader(http.StatusFound)
	}

	return nil
}

func AuthMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {

		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Authorization basic token required",
			})
			return
		}
		const (
			expectedUsername = "user"
			expectedPassword = "pass"
		)
		isValid := (username == expectedUsername) && (password == expectedPassword)
		if !isValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "User password invalid",
			})
		}
		c.Next()
	}
}

func FilterMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {

		err := getfilter.Handle(c)
		if err {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
				"status":  http.StatusNotAcceptable,
				"message": "query invalid",
			})
		}

		c.Next()
	}
}
