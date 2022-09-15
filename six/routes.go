package six

import "github.com/gin-gonic/gin"

func Routes(g *gin.RouterGroup) {
	g.GET("/", GetSixStruct)
}
