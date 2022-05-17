// dish driver for swapi
package dish

import (
	"net/http"

	"github.com/gin-gonic/gin"

	//"github.com/savla-dev/savla-dish/socket"
)


// Sockets and Socket structs has to refer to savla-dish/zasuvka --- should be imported
// https://github.com/savla-dev/savla-dish/blob/master/zasuvka/zasuvka.go#L11
type Sockets struct {
	Sockets []Socket `json:"sockets"`
}

type Socket struct {
	Name         	  	string 		`json:"socket_name"`
	Host	     	  	string 		`json:"host_name"`
	Port         	  	int    		`json:"port_tcp"`
	ExpectedHttpCodes 	[]int  		`json:"expected_http_code_array"`
	PathHttp	  	string 		`json:"path_http"`
	DishList		[]string 	`json:"dish_source"`
	Status		  	[]bool		`json:"status"`
}


// demo socket list
//var socketArray = []zasuvka.Socket{}
var socketArray = []Socket{
	// TCP port check
	{Name: "frank SSH", Host: "frank.savla.su", Port: 22, DishList: []string{"talion", "frank"}},
	//{Name: "frank public SSH", Host: "frank-public.savla.su", Port: 22, DishList: []string{"frank"}},
	{Name: "fgrebox SSH", Host: "fgrebox.savla.su", Port: 22, DishList: []string{"talion", "frank"}},
	{Name: "squabbit SSH", Host: "squabbit.savla.su", Port: 22, DishList: []string{"talion", "frank"}},
	//{Name: "stafle SSH", Host: "stafle.savla.su", Port: 22, DishList: []string{"talion", "frank"}},
	{Name: "stokrle SSH", Host: "stokrle.savla.su", Port: 22, DishList: []string{"talion", "frank"}},
	{Name: "tackbox SSH", Host: "tackbox.savla.su", Port: 22, DishList: []string{"talion", "frank"}},
	{Name: "talion SSH", Host: "talion.savla.su", Port: 22, DishList: []string{"talion", "frank"}},
	{Name: "viking SSH", Host: "viking.savla.su", Port: 22, DishList: []string{"talion", "frank"}},
	{Name: "viking public SSH", Host: "viking-public.savla.su", Port: 22, DishList: []string{"frank"}},

	// TCP port check, cont'd
	{Name: "frank IP intranet DNS", Host: "10.4.5.130", Port: 53, DishList: []string{"frank", "talion"}},
	{Name: "frank NS intranet DNS", Host: "ns.savla.su", Port: 53, DishList: []string{"frank", "talion"}},
	{Name: "frank OpenTTD", Host: "ottd.savla.dev", Port: 3979, DishList: []string{"frank", "talion"}},

	// intranet -- savla.su
	{Name: "savla-docs HTTP", Host: "http://docs.savla.su", Port: 80, ExpectedHttpCodes: []int{ 200 }, PathHttp: "/howto/docs", DishList: []string{"talion", "frank"}},
	{Name: "elden-bling dev HTTP", Host: "http://elden-bling.savla.su", Port: 80, ExpectedHttpCodes: []int{ 200 }, PathHttp: "/", DishList: []string{"talion", "frank"}},
	{Name: "kanban HTTP", Host: "http://kanban.savla.su", Port: 80, ExpectedHttpCodes: []int{ 200 }, PathHttp: "/login", DishList: []string{"talion", "frank"}},
	//{Name: "passbolt HTTP", Host: "http://passbolt.savla.su", Port: 80, ExpectedHttpCodes: []int{ 200 }, PathHttp: "/", DishList: []string{"talion", "frank"}},
	{Name: "swapi head test HTTP", Host: "http://swapi.savla.su", Port: 80, ExpectedHttpCodes: []int{ 200 }, PathHttp: "/dish/test", DishList: []string{"talion", "frank"}},

	// public endpoints -- savla.dev
	{Name: "savla-dev HTTPS", Host: "https://savla.dev", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/", DishList: []string{"talion", "frank"}},
	{Name: "ks-savla-dev HTTPS", Host: "https://ks.savla.dev", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/", DishList: []string{"talion", "frank"}},
	{Name: "kss-savla-dev HTTPS", Host: "https://kss.savla.dev", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/", DishList: []string{"talion", "frank"}},
	{Name: "modul-savla-dev HTTPS", Host: "https://modul.savla.dev", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/", DishList: []string{"talion", "frank"}},
	{Name: "www-savla-dev HTTPS", Host: "https://www.savla.dev", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/", DishList: []string{"talion", "frank"}},

	// public endpoints -- n0p.cz
	{Name: "red-n0p-cz HTTPS", Host: "https://red.n0p.cz", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/login", DishList: []string{"talion", "frank"}},
	{Name: "text-n0p-cz HTTPS", Host: "https://text.n0p.cz", Port: 443, ExpectedHttpCodes: []int{ 200, 301, 401 }, PathHttp: "/", DishList: []string{"talion", "frank"}},
	{Name: "web-n0p-cz HTTPS", Host: "https://web.n0p.cz", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/", DishList: []string{"talion", "frank"}},
	{Name: "wiki-n0p-cz HTTPS", Host: "https://wiki.n0p.cz", Port: 443, ExpectedHttpCodes: []int{ 200, 301, 401 }, PathHttp: "/start", DishList: []string{"talion", "frank"}},

	// public endpoints -- eldenbling.net
	{Name: "eldenbling-net HTTPS", Host: "https://eldenbling.net", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/", DishList: []string{"talion", "frank"}},
	{Name: "www-eldenbling-net HTTPS", Host: "https://www.eldenbling.net", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/", DishList: []string{"talion", "frank"}},

	// public endpoints -- platispivo.cz
	{Name: "platispivo-cz HTTPS", Host: "https://platispivo.cz", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/", DishList: []string{"talion", "frank"}},
	{Name: "www-platispivo-cz HTTPS", Host: "https://www.platispivo.cz", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/", DishList: []string{"talion", "frank"}},

	// legacy
	//{Name: "kyrspa-wz-cz litter HTTP", Host: "http://kyrspa.wz.cz", Port: 80, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/litter/?page=login", DishList: []string{"talion", "frank"}},
	{Name: "smart-comp-cz HTTP", Host: "https://sc.cz", Port: 443, ExpectedHttpCodes: []int{ 200, 301 }, PathHttp: "/", DishList: []string{"talion", "frank"}},
	
}


func HeadTest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"test": true,
	})
}

// GetSocketList GET method
func GetSocketList(c *gin.Context) {
	var sockets = Sockets{
		Sockets: socketArray,
	}

	c.IndentedJSON(http.StatusOK, sockets)
}

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

// GetSocketListByHost GET
func GetSocketListByHost(c *gin.Context) {
	host := c.Param("host")

	var sockets = Sockets{
		Sockets: []Socket{},
	}

        // loop over socketArray, find 
        for _, s := range socketArray {
		if contains(s.DishList, host) {
			// clear the dish source list for the client (dish)
			//s.DishList = []string{host}
			s.DishList = nil
			sockets.Sockets = append(sockets.Sockets, s)
                }
        }

	if len(sockets.Sockets) > 0 {
                c.IndentedJSON(http.StatusOK, sockets)
		return
	}

        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no sockets for given 'hostname'"})

}

