package auth

import "github.com/gin-gonic/gin"

// contains checks if a string is present in a slice
// https://freshman.tech/snippets/go/check-if-slice-contains-element/
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// helper function
func checkMethodUsagePermission(pp AuthParams, ctx *gin.Context) bool {
	roles := pp.Roles
	method := ctx.Request.Method

	switch method {
	case "PUT", "POST", "PATCH", "UPDATE":
		return contains(roles, "power") || contains(roles, "admin")

	case "DELETE":
		return contains(roles, "admin")
	}

	return true
}

func respondWithError(ctx *gin.Context, code int, message interface{}) {
	ctx.AbortWithStatusJSON(code, gin.H{
		"code":    code,
		"message": message,
	})
	return
}
