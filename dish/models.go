package dish

// Sockets and Socket structs has to refer to savla-dish/zasuvka --- should be imported
// https://github.com/savla-dev/savla-dish/blob/master/zasuvka/zasuvka.go#L11
type Sockets struct {
	Sockets []Socket `json:"sockets"`
}

type Socket struct {
	Name              string   `json:"socket_name"`
	Host              string   `json:"host_name"`
	Port              int      `json:"port_tcp"`
	ExpectedHttpCodes []int    `json:"expected_http_code_array"`
	PathHttp          string   `json:"path_http"`
	DishList          []string `json:"dish_source"`
	Status            []bool   `json:"status"`
}

// flush dish socket list array at start
var socketArray = []Socket{}
