package core

import (
	"github.com/gin-gonic/gin"
)

func SetupTestEnv(pkg *Package) *gin.Engine {
	if pkg == nil {
		return nil
	}

	// initialize pkg's cache
	setupTestCache(pkg.Cache)

	// register pkg's routes
	router := setupTestRouter(pkg.Name, pkg.Routes)

	return router
}

// https://circleci.com/blog/gin-gonic-testing/
func setupTestRouter(group string, subRoutes func(pkgRouter *gin.RouterGroup)) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	if group == "" {
		return nil
	}

	// register pkg's route group name
	pkgRouter := router.Group(group)

	if subRoutes == nil {
		return nil
	}

	// register pkg's routes
	subRoutes(pkgRouter)

	return router
}

func setupTestCache(cache **Cache) {
	if cache == nil {
		c := &Cache{}
		cache = &c

		return
	}

	if *cache == nil {
		*cache = &Cache{}
	}

	return
}
