package dish

import (
	"github.com/gin-gonic/gin"
)

// dish CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	// Any generate all possible handler combinations
	// [GIN-debug] GET    /dish/test                --> swis-api/dish.HeadTest (3 handlers)
	// [GIN-debug] POST   /dish/test                --> swis-api/dish.HeadTest (3 handlers)
	// [GIN-debug] PUT    /dish/test                --> swis-api/dish.HeadTest (3 handlers)
	// [GIN-debug] PATCH  /dish/test                --> swis-api/dish.HeadTest (3 handlers)
	// [GIN-debug] HEAD   /dish/test                --> swis-api/dish.HeadTest (3 handlers)
	// [GIN-debug] OPTIONS /dish/test                --> swis-api/dish.HeadTest (3 handlers)
	// [GIN-debug] DELETE /dish/test                --> swis-api/dish.HeadTest (3 handlers)
	// [GIN-debug] CONNECT /dish/test                --> swis-api/dish.HeadTest (3 handlers)
	// [GIN-debug] TRACE  /dish/test                --> swis-api/dish.HeadTest (3 handlers)
	//g.Any("/test", HeadTest)

	g.HEAD("/test", HeadTest)

	g.GET("/sockets", GetSocketList)
	g.GET("/sockets/:host", GetSocketListByHost)
	//g.POST("/sockets/result", PostSocketTestResult)
	g.POST("/sockets/restore", PostDumpRestore)
}
