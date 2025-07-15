package ports

import "github.com/gin-gonic/gin"

type Handler interface {
	RegisterRoute(router *gin.Engine)
}
