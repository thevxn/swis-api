package six

import "github.com/gin-gonic/gin"

func Routes(g *gin.RouterGroup) {
	g.GET("/", GetSixStruct)
	g.GET("/calendar/:owner_name", GetCalendarByUser)
	g.POST("/calendar/:owner_name", PostCalendarItemByUser)
	//g.PUT("/calendar/:owner_name/item/:item_name", UpdateCalendarItemNameByUser)
	g.DELETE("/calendar/:owner_name/item/:item_name", DeleteCalendarItemNameByUser)
	g.POST("/restore", PostDumpRestore)
}
