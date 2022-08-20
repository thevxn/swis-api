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

	// testing route for disg
	g.HEAD("/test", HeadTest)

	// get all sockets loaded
	g.GET("/sockets", GetSocketList)

	// get sockets by hostname/dish-name
	g.GET("/sockets/:host", GetSocketListByHost)

	// add new socket to the list
	g.POST("/sockets", PostNewSocket)

	// edit existing socket by ID
	g.PUT("/sockets/:id", UpdateSocketByID)
	g.PATCH("/sockets/:id", UpdateSocketByID)

	// remove existing socket by ID
	g.DELETE("/sockets/:id", DeleteSocketByID)

	// restore all sockets from JSON dump (JSON-bind)
	g.POST("/sockets/restore", PostDumpRestore)
}
