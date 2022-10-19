package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gsabadini/go-bank-transfer/adapter/api/action"
)

func (g ginEngine) healthcheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		action.HealthCheck(c.Writer, c.Request)
	}
}
