package dish

type Sockets struct {
	//Sockets []Socket `json:"sockets"`
	Sockets map[string]Socket `json:"sockets"`
}

type Socket struct {
	// Socket ID, snake_cased for socket editing and deleting
	ID string `json:"socket_id" validate:"required"`

	// GEneric name of the socket, to be used in dish results as failed one endpoint for example
	Name string `json:"socket_name"`

	// More verbose name/description of the socket
	Description string `json:"socket_description"`

	// Hostname (server.random.com) or HTTP/S URI (http://endpoint.space)
	Host string `json:"host_name" validate:"required"`

	// Socket TCP port part
	// Even default port 80 should be added here
	Port int `json:"port_tcp" validate:"required"`

	// If the Host is HTTP/S endpoint, one can specify which HTTP Result/Response codes are okay and not to alert upon
	ExpectedHTTPCodes []int `json:"expected_http_code_array"`

	// PathHTTP is any URL the site is about to be tested on, e.g. /dish/sockets
	PathHTTP string `json:"path_http"`

	// DishTarget is a string array, usually containing dish's host short name (e.g. frank)
	// to be referred e.g. /dish/sockets/frank
	DishTarget []string `json:"dish_target"`

	// Muted bool indicates that the socket is not propagated to any dish
	Muted bool `json:"muted" default:false`

	// MutedFrom UNIX timestamp.
	MutedFrom int64 `json:"muted_from"`

	// FailCount indicates how many times socket has to be in failed state before alerting.
	FailCount int `json:"fail_count" default:0`

	// Status object for dish results to be returned/updated (by dish itself)
	// Note: discontinued as dish now reports to pushgateway of prometheus'
	//Status []bool `json:"status"`
}
