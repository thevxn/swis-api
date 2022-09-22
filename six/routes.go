package six

import "github.com/gin-gonic/gin"

func Routes(g *gin.RouterGroup) {
	g.GET("/", GetSixStruct)
	g.GET("/calendar/:owner_name", GetCalendarByUser)
	g.POST("/calendar/:owner_name", PostCalendarItemByUser)
	g.POST("/restore", PostDumpRestore)
}
