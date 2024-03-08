package core

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

func RegisterAndMount(parentRouter *gin.Engine, pkgs ...*Package) {
	if parentRouter == nil {
		return
	}

	for _, pkg := range pkgs {
		if pkg == nil {
			continue
		}

		// TODO: catch the error
		MountPackage(parentRouter, pkg)
	}
}

func MountPackage(router *gin.Engine, pkg *Package) bool {
	if pkg == nil {
		fmt.Errorf("failed to mount a package: Package input cannot be nil")
		return false
	}

	if pkg.Name == "" || &pkg.Name == nil {
		fmt.Errorf("failed to mount a package: Name cannot be blank")
		return false
	}

	if err := initCache(pkg.Cache); err != nil {
		fmt.Errorf("failed to mount '%s' package: %s", pkg.Name, err.Error())
		return false
	}

	if err := mountRouterGroup(router, pkg.Name, pkg.Routes); err != nil {
		fmt.Errorf("failed to mount '%s' package: %s", pkg.Name, err.Error())
		return false
	}

	return true
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

func mountRouterGroup(router *gin.Engine, groupName string, subRoutes func(r *gin.RouterGroup)) error {
	if router == nil {
		return errors.New("nil router pointer")
	}

	if groupName == "" {
		return errors.New("blank groupName parameter")
	}

	if &subRoutes == nil {
		return errors.New("nil pointer to routes function")
	}

	// register pkg's route group name
	pkgRouter := router.Group(groupName)

	// register pkg's routes
	subRoutes(pkgRouter)

	return nil
}
