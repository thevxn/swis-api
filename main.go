// go:build ignore

// Package swis-api is RESTful API core backend for sakalWeb Information System v5.
package main

import (
	"net/http"

	"swis-api/dish"
	"swis-api/groups"
	"swis-api/users"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// reqs from this IPs are treated as proxies, ergo log the real client IP address
	/*swapiProxies := []string{
		"10.4.5.130/25",
	}*/

	//router.SetTrustedProxies(swapiProxies)
	
	router.GET("/", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"message": "welcome to sakalweb API (swapi) root",
		})
	})

	// dish CRUD
	router.HEAD("/dish/test", dish.HeadTest)
	router.GET("/dish/sockets", dish.GetSocketList)
	router.GET("/dish/sockets/:host", dish.GetSocketListByHost)

	// groups CRUD
	router.GET("/groups", groups.GetGroups)
	router.GET("/groups/:id", groups.GetGroupByID)
	router.POST("/groups", groups.PostGroup)
	//router.PUT("/groups/:id", groups.PutGroupByID)
	//router.DELETE("/groups/:id", groups.DeleteGroupByID)

	// users CRUD
	router.GET("/users", users.GetUsers)
	router.GET("/users/:id", users.GetUserByID)
	router.POST("/users", users.PostUser)
	//router.PUT("/users/:id", users.PutUserByID)
	//router.DELETE("/users/:id", users.DeleteUserByID)

	// attach router to http.Server and start it
	router.Run(":8080")
}

