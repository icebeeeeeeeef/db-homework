package web

import "github.com/gin-gonic/gin"

type handler interface {
	RegisterRoutes(r *gin.Engine)
}
