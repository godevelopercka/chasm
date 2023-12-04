package web

import "github.com/gin-gonic/gin"

type handler interface {
	RegisterUserRoutes(server *gin.Engine)
}
