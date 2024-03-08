package core

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func MountPackage(pkg *Package) *gin.Engine {
	if pkg == nil {
		fmt.Errorf("failed to mount a package: Package input cannot be nil")
		return nil
	}

	if pkg.Name == "" || *pkg.Name == nil {
		fmt.Errorf("failed to mount a package: Name cannot be blank")
		return nil
	}

	if err := initCache(pkg.Cache); err != nil {
		fmt.Errorf("failed to mount '%s' package: %s", pkg.Name, err.Error())
		return nil
	}

	if err := mountRouterGroup(pkg.Name, pkg.Routes); err != nil {
		fmt.Errorf("failed to mount '%s' package: %s", pkg.Name, err.Error())
		return nil
	}
}

func initCache(cache **Cache) error {
	if cache == nil {
		c := &Cache{}
		cache = &c

		return nil
	}

	if *cache == nil {
		*cache = &Cache{}
	}

	return nil
}

func mountRouterGroup(router *gin.Server, groupName string, subRoutes func(r *gin.RouterGroup)) error {
	if router == nil {
		return errors.New("nil router pointer")
	}

	if groupName == "" {
		return errors.New("blank groupName parameter")
	}

	if *subRoutes == nil {
		return errors.New("nil pointer to routes function")
	}

	// register pkg's route group name
	pkgRouter := router.Group(groupName)

	// register pkg's routes
	subRoutes(pkgRouter)

	return nil
}
