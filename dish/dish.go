// dish driver for swapi
package dish

import (
	"net/http"

	"github.com/gin-gonic/gin"

	//"github.com/savla-dev/savla-dish/zasuvka"
)


// Sockets and Socket structs has to refer to savla-dish/zasuvka --- should be imported
// https://github.com/savla-dev/savla-dish/blob/master/zasuvka/zasuvka.go#L11
type Sockets struct {
	Sockets []Socket `json:"sockets"`
}

type Socket struct {
	Name         	  string `json:"socket_name"`
	Host	     	  string `json:"host_name"`
	Port         	  int    `json:"port_tcp"`
	ExpectedHttpCodes []int  `json:"expected_http_code_array"`
	PathHttp	  string `json:"path_http"`
}


// demo socket
//var socketArray = []zasuvka.Socket{}
var socketArray = []Socket{
	{Name: "frank SSH", Host: "frank.savla.su", Port: 22, ExpectedHttpCodes: []int{}, PathHttp: ""},
}


func GetSocketList(c *gin.Context) {
	//var sockets = zasuvka.Sockets{}
	var sockets = Sockets{
		Sockets: socketArray,
	}

	c.IndentedJSON(http.StatusOK, sockets)
}


//func GetSocketListByHost(c *gin.Context) {
	//host := c.Param("host")
	//c.IndentedJSON(http.StatusOK, sockets)
//}

