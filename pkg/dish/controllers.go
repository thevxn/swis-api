package dish

import (
	//"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	//"go.savla.dev/dish/pkg/socket"
	"go.savla.dev/swis/v5/pkg/core"

	"github.com/gin-gonic/gin"
)

var (
	CacheIncidents *core.Cache
	CacheSockets   *core.Cache
	Dispatcher     *Stream
	pkgName        string = "dish"
)

var Package *core.Package = &core.Package{
	Name: pkgName,
	Cache: []**core.Cache{
		&CacheIncidents,
		&CacheSockets,
	},
	Routes: Routes,
	Subpackages: []string{
		"incidents",
		"sockets",
	},
}

/*

  sockets

*/

// @Summary Get all sockets list
// @Description get socket list, socket array
// @Tags dish
// @Produce  json
// @Success 200 {object} string "ok"
// @Router /dish/sockets [get]
func GetSocketList(ctx *gin.Context) {
	core.PrintAllRootItems(ctx, CacheSockets, pkgName)
	return
}

// Add new socket to the list.
// @Summary Adding new socket to socket array
// @Description add new socket to socket array
// @Tags dish
// @Produce json
// @Param request body dish.Socket true "query params"
// @Success 200 {object} dish.Socket
// @Router /dish/sockets/{key} [post]
func PostNewSocketByKey(ctx *gin.Context) {
	core.AddNewItemByParam[Socket](ctx, CacheSockets, pkgName, Socket{})
	return
}

// edit existing socket by ID
// @Summary Update socket by its ID
// @Description update socket by its ID
// @Tags dish
// @Produce json
// @Param request body dish.Socket.ID true "query params"
// @Success 200 {object} dish.Socket
// @Router /dish/sockets/{key} [put]
func UpdateSocketByKey(ctx *gin.Context) {
	core.UpdateItemByParam[Socket](ctx, CacheSockets, pkgName, Socket{})
	return
}

// remove existing socket by ID
// @Summary Delete socket by its ID
// @Description delete socket by its ID
// @Tags dish
// @Produce json
// @Param  id  path  string  true  "dish ID"
// @Success 200 {object} dish.Socket
// @Router /dish/sockets/{key} [delete]
func DeleteSocketByKey(ctx *gin.Context) {
	core.DeleteItemByParam(ctx, CacheSockets, pkgName)
	return
}

// @Summary Get public socket list
// @Description get public socket list
// @Tags dish
// @Produce json
// @Success 200 {string} string	"ok"
// @Router /dish/sockets/public [get]
func GetSocketListPublic(ctx *gin.Context) {
	var exportedSockets = make(map[string]Socket)
	var counter int = 0

	rawSocketsMap, _ := CacheSockets.GetAll()

	for _, rawSocket := range rawSocketsMap {
		socket, ok := rawSocket.(Socket)
		if !ok {
			continue
		}

		if socket.Public {
			exportedSockets[socket.ID] = socket
			counter++
		}
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"count":   counter,
		"items":   exportedSockets,
		"message": "ok, dumping public sockets",
	})
	return
}

// @Summary Get socket list by host
// @Description get socket list by Host
// @Tags dish
// @Produce json
// @Param host path string true "dish instance name"
// @Success 200 {string} string	"ok"
// @Router /dish/sockets/{host} [get]
// Get sockets by hostname/dish-name.
func GetSocketListByHost(ctx *gin.Context) {
	var host string = ctx.Param("host")
	var exportedSockets = make(map[string]Socket)
	var counter int = 0

	rawSocketsMap, _ := CacheSockets.GetAll()

	for _, rawSocket := range rawSocketsMap {
		socket, ok := rawSocket.(Socket)
		if !ok {
			continue
		}

		if contains(socket.DishTarget, host) && !socket.Muted {
			exportedSockets[socket.ID] = socket
			counter++
		}
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"count":   counter,
		"items":   exportedSockets,
		"message": "ok, dumping sockets by host",
		"host":    host,
	})
	return
}

// @Summary Batch update socket's healthy state.
// @Description batch update socket's healthy state.
// @Tags dish
// @Produce json
// @Router /dish/sockets/results [post]
func BatchPostHealthyStatus(ctx *gin.Context) {
	var results = struct {
		Map map[string]bool `json:"dish_results"`
	}{}

	if err := ctx.BindJSON(&results); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
			"message": "cannot bind input JSON stream",
			"package": pkgName,
		})
		return
	}

	var socketsDown []string
	var socketsUp []string
	var count int = 0

	for key, result := range results.Map {
		var socket Socket
		var ok bool

		if rawSocket, found := CacheSockets.Get(key); !found {
			continue
		} else {
			if socket, ok = rawSocket.(Socket); !ok {
				ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"key":     key,
					"message": "cannot assert data type, database internal error",
					"package": pkgName,
				})
				return
			}
		}

		// add socket ID to the exported array (via event dispatcher) if changed its state only
		if socket.Healthy != result {
			if result {
				socketsUp = append(socketsUp, socket.ID)
			} else {
				socketsDown = append(socketsDown, socket.ID)
			}
			count++
			socket.Healthy = result
		}

		socket.TestTimestamp = time.Now().UnixNano()

		if saved := CacheSockets.Set(key, socket); !saved {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"key":     key,
				"message": "cannot update socket's healthy state by key",
			})
			return
		}
	}

	if count > 0 {
		if len(socketsUp) > 0 {
			// generate and send a SSE message
			msg := Message{
				Content:    "socket-up",
				SocketList: socketsUp,
				Timestamp:  time.Now().UnixNano(),
			}

			// emit an server-sent event to subscribers
			if Dispatcher != nil {
				Dispatcher.NewMessage(msg)
			}
		}

		if len(socketsDown) > 0 {
			// generate and send a SSE message
			msg := Message{
				Content:    "socket-down",
				SocketList: socketsDown,
				Timestamp:  time.Now().UnixNano(),
			}

			// emit an server-sent event to subscribers
			if Dispatcher != nil {
				Dispatcher.NewMessage(msg)
			}
		}
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, healthy booleans updated per socket",
		"count":   count,
	})
	return
}

// MaintenanceToggleSocketByKey sets maintenance mode of a socket by its ID
//
// @Summary      Toggle maintenance mode
// @Description  toggle maintenance mode
// @Tags         dish
// @Accept       json
// @Produce      json
// @Param        key  query     string  false  "name socket by key"
// @Success      200  {array}   dish.Socket
// @Failure      404  {object}  dish.Socket
// @Failure      500  {object}  dish.Socket
// @Router       /dish/sockets/{key}/maintenance [put]
func MaintenanceToggleSocketByKey(ctx *gin.Context) {
	var id string = ctx.Param("key")
	var updatedSocket Socket

	rawSocket, ok := CacheSockets.Get(id)
	if !ok {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "socket not found",
			"id":      id,
		})
		return
	}

	updatedSocket, ok = rawSocket.(Socket)
	if !ok {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "cannot assert data type, database internal error",
		})
		return
	}

	// inverse the Maintenance field value
	updatedSocket.Maintenance = !updatedSocket.Maintenance

	if updatedSocket.Maintenance {
		updatedSocket.Muted = true
	}

	if saved := CacheSockets.Set(updatedSocket.ID, updatedSocket); !saved {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "socket couldn't be saved to database",
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "socket mute toggle pressed!",
		"socket":  updatedSocket,
	})
	return
}

// MuteToggleSocketByKey sets muted state of a socket by its ID
//
// @Summary      Toggle muted state
// @Description  toggle muted state
// @Tags         dish
// @Accept       json
// @Produce      json
// @Param        key  query     string  false  "name socket by key"
// @Success      200  {array}   dish.Socket
// @Failure      404  {object}  dish.Socket
// @Failure      500  {object}  dish.Socket
// @Router       /dish/sockets/{key}/mute [put]
func MuteToggleSocketByKey(ctx *gin.Context) {
	var id string = ctx.Param("key")
	var updatedSocket Socket

	rawSocket, ok := CacheSockets.Get(id)
	if !ok {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "socket not found",
			"id":      id,
		})
		return
	}

	updatedSocket, ok = rawSocket.(Socket)
	if !ok {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "cannot assert data type, database internal error",
		})
		return
	}

	// inverse the Muted field value
	updatedSocket.Muted = !updatedSocket.Muted

	if updatedSocket.Muted {
		updatedSocket.MutedFrom = time.Now().Unix()
	}

	if saved := CacheSockets.Set(updatedSocket.ID, updatedSocket); !saved {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "socket couldn't be saved to database",
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "socket mute toggle pressed!",
		"socket":  updatedSocket,
	})
	return
}

// PublicToggleSocketByKey sets public state of a socket by its ID
//
// @Summary      Toggle public state
// @Description  toggle public state
// @Tags         dish
// @Accept       json
// @Produce      json
// @Param        key  query     string  false  "name socket by key"
// @Success      200  {array}   dish.Socket
// @Failure      404  {object}  dish.Socket
// @Failure      500  {object}  dish.Socket
// @Router       /dish/sockets/{key}/public [put]
func PublicToggleSocketByKey(ctx *gin.Context) {
	var id string = ctx.Param("key")
	var updatedSocket Socket

	rawSocket, ok := CacheSockets.Get(id)
	if !ok {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "socket not found",
			"id":      id,
		})
		return
	}

	updatedSocket, ok = rawSocket.(Socket)
	if !ok {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "cannot assert data type, database internal error",
		})
		return
	}

	// inverse the Muted field value
	updatedSocket.Public = !updatedSocket.Public

	if saved := CacheSockets.Set(updatedSocket.ID, updatedSocket); !saved {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "socket couldn't be saved to database",
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "socket public toggle pressed!",
		"socket":  updatedSocket,
	})
	return
}

// GetSSEvents
//
// @Summary      Subscribe to dish SSE dispatcher
// @Description  subscribe to dish SSE dispatcher
// @Tags         dish
// @Accept       json
// @Produce      json
// @Success      200  {array}   dish.Message
// @Router       /dish/sockets/status [get]
func GetSSEvents(ctx *gin.Context) {
	// initialize client channel
	clientChan := make(chan string)

	// send new connection to event server
	Dispatcher.NewClients <- clientChan

	defer func() {
		// send closed connection to event server
		Dispatcher.ClosedClients <- clientChan
	}()

	// set the stream headers
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")

	if closed := ctx.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-clientChan:
			if ok {
				ctx.SSEvent("message", msg)
				return true
			}

		// connection is closed then defer will be executed
		case <-ctx.Done():
			return false
		}

		return true
	}); closed {
		//log.Println("closing eventChannel")
		//close(eventChannel)
	}
	return
}

/*

  incidents

*/

// GetIncidentList lists all available incidents
//
// @Summary      Get all incidents
// @Description  get all incidents
// @Tags         dish
// @Accept       json
// @Produce      json
// @Success      200  {array}   dish.Incident
// @Router       /dish/incidents [get]
func GetIncidentList(ctx *gin.Context) {
	core.PrintAllRootItems(ctx, CacheIncidents, pkgName)
	return
}

// PostNewIncident produces and broadcasts a new incident to be happening
//
// @Summary      Add new incident
// @Description  add new incident
// @Tags         dish
// @Accept       json
// @Produce      json
// @Param        key  query     dish.Incident true "socket body"
// @Success      201  {object}   dish.Incident
// @Failure      400  {object}  dish.Incident
// @Failure      409  {object}  dish.Incident
// @Failure      500  {object}  dish.Incident
// @Router       /dish/incidents [post]
func PostNewIncident(ctx *gin.Context) {
	//core.AddNewItemByParam(ctx, CacheIncidents, pkgName, Incident{})

	var id string

	// loop until new incident ID is generated and usable
	for {
		id = strconv.FormatInt(time.Now().Unix(), 10)

		_, found := CacheIncidents.Get(id)
		if !found {
			break
		}
	}

	var newIncident Incident

	if err := ctx.BindJSON(&newIncident); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
			"message": "cannot parse input JSON stream",
		})
		return
	}

	newIncident.ID = id

	if saved := CacheIncidents.Set(id, newIncident); !saved {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "incident couldn't be saved to database",
		})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, gin.H{
		"code":     http.StatusCreated,
		"message":  "new incident created",
		"incident": newIncident,
		"id":       id,
	})
	return
}

// UpdateIncidentByKey update requested incident
//
// @Summary      Update incident by its key
// @Description  update incident by its key
// @Tags         dish
// @Accept       json
// @Produce      json
// @Param        key  query     dish.Incident.ID true "socket body"
// @Success      200  {array}   dish.Incident
// @Failure      400  {object}  dish.Incident
// @Failure      404  {object}  dish.Incident
// @Failure      500  {object}  dish.Incident
// @Router       /dish/incidents/{key} [put]
func UpdateIncidentByKey(ctx *gin.Context) {
	core.UpdateItemByParam[Incident](ctx, CacheIncidents, pkgName, Incident{})
	return
}

// DeleteIncidentByKey deletes given incident
//
// @Summary      Delete incident by its key
// @Description  delete incident by its key
// @Tags         dish
// @Accept       json
// @Produce      json
// @Param        key  query     dish.Incident.ID true "socket name"
// @Success      200  {array}   dish.Incident
// @Failure      404  {object}  dish.Incident
// @Failure      500  {object}  dish.Incident
// @Router       /dish/incidents/{key} [delete]
func DeleteIncidentByKey(ctx *gin.Context) {
	core.DeleteItemByParam(ctx, CacheIncidents, pkgName)
	return
}

// GetGlobalIncidentList returns list of global incidents (no socketID assigned)
//
// @Summary      Get global incident list
// @Description  get global incident list
// @Tags         dish
// @Accept       json
// @Produce      json
// @Success      200  {array}   dish.Incident
// @Router       /dish/incidents/global [get]
func GetGlobalIncidentList(ctx *gin.Context) {
	var exportedIncidents []Incident = []Incident{}
	var counter int = 0

	rawIncidentsMap, _ := CacheIncidents.GetAll()

	for _, rawIncident := range rawIncidentsMap {
		incident, ok := rawIncident.(Incident)
		if !ok {
			continue
		}

		if incident.SocketID == "" {
			exportedIncidents = append(exportedIncidents, incident)
			counter++
		}
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"count":   counter,
		"items":   exportedIncidents,
		"message": "ok, dumping global incidents list",
	})
	return
}

// GetPublicIncidentList returns list of public incidents only
//
// @Summary      Get public incident list
// @Description  get public incident list
// @Tags         dish
// @Accept       json
// @Produce      json
// @Success      200  {array}   dish.Incident
// @Router       /dish/incidents/public [get]
func GetPublicIncidentList(ctx *gin.Context) {
	var exportedIncidents []Incident = []Incident{}
	var counter int = 0

	rawIncidentsMap, _ := CacheIncidents.GetAll()

	for _, rawIncident := range rawIncidentsMap {
		incident, ok := rawIncident.(Incident)
		if !ok {
			continue
		}

		if incident.Public {
			exportedIncidents = append(exportedIncidents, incident)
			counter++
		}
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"count":   counter,
		"items":   exportedIncidents,
		"message": "ok, dumping public incidents list",
	})
	return
}

// GetIncidentListBySocketID returns list of incidents for such Socket.ID
//
// @Summary      Get incident list by socket ID
// @Description  get incident list by socket ID
// @Tags         dish
// @Accept       json
// @Produce      json
// @Param        key  query  dish.Socket.ID true "socket name"
// @Success      200  {array}   dish.Incident
// @Failure      404  {object}  dish.Incident
// @Router       /dish/incidents/{key} [get]
func GetIncidentListBySocketID(ctx *gin.Context) {
	var key string = ctx.Param("key")
	//var exportedIncidents = make(map[string]Incident)
	var exportedIncidents []Incident = []Incident{}
	var counter int = 0

	if CacheSockets == nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"count":   counter,
			"message": "cannot access socket cache",
			"key":     key,
		})
		return
	}

	if _, ok := CacheSockets.Get(key); !ok {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"count":   0,
			"items":   exportedIncidents,
			"message": "no such socket",
			"key":     key,
		})
		return
	}

	rawIncidentsMap, _ := CacheIncidents.GetAll()

	for _, rawIncident := range rawIncidentsMap {
		incident, ok := rawIncident.(Incident)
		if !ok {
			continue
		}

		if incident.SocketID == key {
			//exportedIncidents[incident.SocketID] = incident
			exportedIncidents = append(exportedIncidents, incident)
			counter++
		}
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"count":   counter,
		"items":   exportedIncidents,
		"message": "ok, dumping incidents list by socketID",
		"key":     key,
	})
	return
}

/*

  restoration and types

*/

// @Summary List package model's field types
// @Description list package model's field types
// @Tags dish
// @Accept json
// @Produce json
// @Router /dish/incidents/types [get]
func ListTypesIncidents(ctx *gin.Context) {
	core.ParsePackageType(ctx, pkgName, Incident{})
	return
}

// @Summary List package model's field types
// @Description list package model's field types
// @Tags dish
// @Accept json
// @Produce json
// @Router /dish/sockets/types [get]
func ListTypesSockets(ctx *gin.Context) {
	core.ParsePackageType(ctx, pkgName, Socket{})
	return
}

// GetDishRoot returns all possible items of dish package
//
// @Summary      Get all root items
// @Description  get all root items
// @Tags         dish
// @Accept       json
// @Produce      json
// @Param        key  query     dish.Incident.ID true "socket name"
// @Success      200  {object}   dish.Root
// @Router       /dish [get]
func GetDishRoot(ctx *gin.Context) {
	incidents, _ := CacheIncidents.GetAll()
	sockets, _ := CacheSockets.GetAll()

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":      http.StatusOK,
		"message":   "ok, dumping dish root",
		"incidents": incidents,
		"sockets":   sockets,
	})
}

// PostDumpRestore is a helper function to load all root items
//
// @Summary      Restore dish package items
// @Description  restore dish package items
// @Tags         dish
// @Accept       json
// @Produce      json
// @Success      201  {array}   dish.Incident
// @Failure      404  {array}   dish.Incident
// @Failure      500  {array}   dish.Incident
// @Router       /dish/restore [post]
func PostDumpRestore(ctx *gin.Context) {
	var counter []int = []int{0, 0}

	var importDish = struct {
		Incidents map[string]Incident `json:"incidents"`
		Sockets   map[string]Socket   `json:"sockets"`
	}{}

	if err := ctx.BindJSON(&importDish); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
			"message": "cannot parse input JSON stream",
		})
		return
	}

	for key, item := range importDish.Incidents {
		if key == "" {
			continue
		}

		CacheIncidents.Set(key, item)
		counter[0]++
	}

	for key, item := range importDish.Sockets {
		if key == "" {
			continue
		}

		CacheSockets.Set(key, item)
		counter[1]++
	}

	// HTTP 201 Created
	ctx.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"counter": counter,
		"message": "dish dump imported successfully",
	})
}
