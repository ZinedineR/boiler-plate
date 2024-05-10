package api

import (
	"boiler-plate/pkg/getfilter"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseHeaderFormat() gin.HandlerFunc {
	return func(c *gin.Context) {
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
