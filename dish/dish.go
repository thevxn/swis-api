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


// demo socket list
//var socketArray = []zasuvka.Socket{}
var socketArray = []Socket{
	// TCP port check
	{Name: "frank SSH", Host: "frank.savla.su", Port: 22, ExpectedHttpCodes: []int{}, PathHttp: ""},
	{Name: "frank public SSH", Host: "frank-public.savla.su", Port: 22, ExpectedHttpCodes: []int{}, PathHttp: ""},
	{Name: "fgrebox SSH", Host: "fgrebox.savla.su", Port: 22, ExpectedHttpCodes: []int{}, PathHttp: ""},
	{Name: "squiabbit SSH", Host: "squabbit.savla.su", Port: 22, ExpectedHttpCodes: []int{}, PathHttp: ""},
	{Name: "stafle SSH", Host: "stafle.savla.su", Port: 22, ExpectedHttpCodes: []int{}, PathHttp: ""},
	{Name: "stokrle SSH", Host: "storkle.savla.su", Port: 22, ExpectedHttpCodes: []int{}, PathHttp: ""},
	{Name: "tackbox SSH", Host: "tackbox.savla.su", Port: 22, ExpectedHttpCodes: []int{}, PathHttp: ""},
	{Name: "talion SSH", Host: "talion.savla.su", Port: 22, ExpectedHttpCodes: []int{}, PathHttp: ""},
	{Name: "viking SSH", Host: "viking.savla.su", Port: 22, ExpectedHttpCodes: []int{}, PathHttp: ""},
	{Name: "viking public SSH", Host: "viking-public.savla.su", Port: 22, ExpectedHttpCodes: []int{}, PathHttp: ""},

	// TCP port check, cont'd
	{Name: "frank IP intranet DNS", Host: "10.4.5.130", Port: 53, ExpectedHttpCodes: []int{}, PathHttp: ""},
	{Name: "frank NS intranet DNS", Host: "ns.savla.su", Port: 53, ExpectedHttpCodes: []int{}, PathHttp: ""},
	{Name: "frank OpenTTD", Host: "ottd.savla.dev", Port: 3979, ExpectedHttpCodes: []int{}, PathHttp: ""},

	// intranet -- savla.su
	{Name: "savla-docs HTTP", Host: "http://docs.savla.su", Port: 80, ExpectedHttpCodes: []int{ 200 }, PathHttp: "/howto/docs"},
	{Name: "elden-bling dev HTTP", Host: "http://elden-bling.savla.su", Port: 80, ExpectedHttpCodes: []int{ 200 }, PathHttp: "/"},
	{Name: "kanban HTTP", Host: "http://kanban.savla.su", Port: 80, ExpectedHttpCodes: []int{ 200 }, PathHttp: "/login"},
	{Name: "passbolt HTTP", Host: "http://passbolt.savla.su", Port: 80, ExpectedHttpCodes: []int{ 200 }, PathHttp: "/"},
	{Name: "swapi head test HTTP", Host: "http://swapi.savla.su", Port: 80, ExpectedHttpCodes: []int{ 200 }, PathHttp: "/dish/test"},

	// public endpoints -- savla.dev
	{Name: "savla-dev HTTPS", Host: "https://savla.dev", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/"},
	{Name: "ks-savla-dev HTTPS", Host: "https://ks.savla.dev", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/"},
	{Name: "kss-savla-dev HTTPS", Host: "https://kss.savla.dev", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/"},
	{Name: "modul-savla-dev HTTPS", Host: "https://modul.savla.dev", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/"},
	{Name: "www-savla-dev HTTPS", Host: "https://www.savla.dev", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/"},

	// public endpoints -- n0p.cz
	{Name: "red-n0p-cz HTTPS", Host: "https://red.n0p.cz", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/login"},
	{Name: "text-n0p-cz HTTPS", Host: "https://text.n0p.cz", Port: 443, ExpectedHttpCodes: []int{ 301, 401 }, PathHttp: "/"},
	{Name: "web-n0p-cz HTTPS", Host: "https://web.n0p.cz", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/"},
	{Name: "wiki-n0p-cz HTTPS", Host: "https://wiki.n0p.cz", Port: 443, ExpectedHttpCodes: []int{ 200, 301, 401 }, PathHttp: "/start"},

	// public endpoints -- eldenbling.net
	{Name: "eldenbling-net HTTPS", Host: "https://eldenbling.net", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/"},
	{Name: "www-eldenbling-net HTTPS", Host: "https://www.eldenbling.net", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/"},

	// public endpoints -- platispivo.cz
	{Name: "platispivo-cz HTTPS", Host: "https://platispivo.cz", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/"},
	{Name: "www-platispivo-cz HTTPS", Host: "https://www.platispivo.cz", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/"},

	// legacy
	{Name: "kyrspa-wz-cz litter HTTP", Host: "http://kyrspa.wz.cz", Port: 80, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/litter/?page=login"},
	{Name: "smart-comp-cz HTTP", Host: "https://sc.cz", Port: 80, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/"},
	
}


func HeadTest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"test": true,
	})
}

// GetSocketList GET method
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

