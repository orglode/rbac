package http

import (
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	r.GET("test", test)
}
