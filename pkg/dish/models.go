package dish

type Root struct {
	Incidents map[string]Incident `json:"incidents"`
	Sockets   map[string]Socket   `json:"sockets"`
}

type Socket struct {
	// Socket ID, snake_cased for socket editing and deleting.
	ID string `json:"socket_id" binding:"required" validation:"required"`

	// GEneric name of the socket, to be used in dish results as failed one endpoint for example.
	Name string `json:"socket_name" binding:"required"`

	// More verbose name/description of the socket.
	Description string `json:"socket_description"`

	// Hostname (server.random.com) or HTTP/S URI (https://endpoint.space).
	Host string `json:"host_name" binding:"required" validation:"required"`

	// Socket TCP port part
	// Even default port 80 should be added here.
	Port int `json:"port_tcp" binding:"required" validation:"required"`

	// If the Host is HTTP/S endpoint, one can specify which HTTP Result/Response codes are okay and not to alert upon.
	ExpectedHTTPCodes []int `json:"expected_http_code_array"`

	// PathHTTP is any URL the site is about to be tested on, e.g. /dish/sockets
	PathHTTP string `json:"path_http"`

	// DishTarget is a string array, usually containing dish's host short name (e.g. frank).
	// To be referred as /dish/sockets/frank for example.
	DishTarget []string `json:"dish_target"`

	// Muted bool indicates that the socket is not propagated to any dish if true.
	Muted bool `json:"muted" default:true`

	// MutedFrom UNIX timestamp.
	MutedFrom int64 `json:"muted_from"`

	// FailCount indicates how many times socket has to be in failed state before alerting.
	FailCount int `json:"fail_count" default:0`

	// ResponseTime is the time for the request to be processed.
	ResponseTime float64 `json:"response_time"`

	// TestTimestamp tells the time of the last socket testing being executed upon.
	TestTimestamp int64 `json:"test_timestamp"`

	// Healthy boolean indicates wheter is socket okay, or the way around.
	Healthy bool `json:"healthy" default:false`

	// Public boolean tells the frontendee to show itself.
	Public bool `json:"public" default:false`

	// Maintenance boolean states for the M. mode being applied to such socket/endpoint.
	Maintenance bool `json:"maintenance" default>false`
}

type Incident struct {
	// Incident ID, stringified timestamp usually.
	ID string `json:"id"`

	// Incident name.
	Name string `json:"name" binding:"required"`

	// Further details about the incident like place, state of operation etc.
	Description string `json:"desc"`

	// Type of incident, e.g. planned, maintenance, outage etc
	Type string `json:"type"`

	// ID of the referencing socket(s).
	//SocketID []string `json:"socket_id"`
	SocketID string `json:"socket_id"`

	// The very start datetime of such incident.
	StartTimestamp int64 `json:"start_date"`

	// Estimated end of incident handling/resolving.
	EndTimestamp int64 `json:"end_date"`

	// Reason of the incident that happened.
	Reason string `json:"reason"`

	// Public indicates the state of visibility for all.
	Public bool `json:"public" default:false`

	// Other commentary to the incident.
	Comment string `json:"comment"`
}

type ClientChan chan Message

// Stream is a SSE data structure
type Stream struct {
	// Message is a string of volatile length that carries the very event saying and its metadata.
	Message chan string `json:"message"`

	// Client connections monitoring
	NewClients    chan chan string
	ClosedClients chan chan string
	TotalClients  map[chan string]bool
}

type Message struct {
	Content    string   `json:"content"`
	SocketList []string `json:"socket_list"`
	Timestamp  int64    `json:"timestamp"`
}
