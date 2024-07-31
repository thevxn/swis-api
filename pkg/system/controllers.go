package system

import (
	"fmt"
	"net/http"
	//"strconv"
	//"time"

	"go.savla.dev/swis/v5/pkg/core"

	"github.com/gin-gonic/gin"
)

var (
	Cache   *core.Cache
	pkgName string = "system"
)

var Package *core.Package = &core.Package{
	Name: pkgName,
	Cache: []**core.Cache{
		&Cache,
	},
	Routes: Routes,
}

func GetMountedPackages(ctx *gin.Context) {
	items, _ := Cache.Get("mounted")

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"items":   items,
		"message": fmt.Sprintf("ok, listing mounted package list"),
		"package": pkgName,
	})
	return
}
