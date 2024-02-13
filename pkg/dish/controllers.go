package dish

import (
	//"encoding/json"
	"io"
	"net/http"
	"time"

	//"go.savla.dev/dish/pkg/socket"
	"go.savla.dev/swis/v5/pkg/core"

	"github.com/gin-gonic/gin"
)

var (
	CacheIncidents *core.Cache
	CacheSockets   *core.Cache
	eventChannel   chan string
	pkgName        string = "dish"
)

/*

  sockets

*/

// Get all sockets loaded.
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
	core.AddNewItemByParam(ctx, CacheSockets, pkgName, Socket{})
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
	core.UpdateItemByParam(ctx, CacheSockets, pkgName, Socket{})
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

		// nasty tweak incoming
		if (contains(socket.DishTarget, host) && !socket.Muted) || (host == "public" && contains(socket.DishTarget, host) && socket.Maintenance) {
			exportedSockets[socket.ID] = socket
			counter++
		}
	}

	if len(exportedSockets) > 0 {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"count":   counter,
			"items":   exportedSockets,
			"message": "ok, dumping socket by host",
			"host":    host,
		})
		return
	}

	ctx.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "no sockets for given 'host'",
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

	var sockets []string
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
			sockets = append(sockets, socket.ID)
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
		// generate and send a SSE message
		msg := Message{
			Content:    "sockets updated",
			SocketList: sockets,
			Timestamp:  time.Now().UnixNano(),
		}

		//log.Println("sockets updated message sent")
		//Dispatcher.NewEvent(msg)
		if EventChannel != nil {
			EventChannel <- msg
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
//	@Summary      Togle maintenance mode
//	@Description  toggle maintenance mode
//	@Tags         dish
//	@Accept       json
//	@Produce      json
//	@Param        key  query     string  false  "name socket by key"
//	@Success      200  {array}   dish.Socket
//	@Failure      404  {object}  dish.Socket
//	@Failure      500  {object}  dish.Socket
//	@Router       /dish/sockets/:key/maintenance [put]
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
//	@Summary      Togle muted state
//	@Description  toggle muted state
//	@Tags         dish
//	@Accept       json
//	@Produce      json
//	@Param        key  query     string  false  "name socket by key"
//	@Success      200  {array}   dish.Socket
//	@Failure      404  {object}  dish.Socket
//	@Failure      500  {object}  dish.Socket
//	@Router       /dish/sockets/:key/mute [put]
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
	//log.Println("opening eventChannel")
	//eventChannel = make(chan string)

	// set the stream headers
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")

	if closed := ctx.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-EventChannel:
			if ok {
				ctx.SSEvent("swis-event", msg)
				return true
			}

		// connection is closed then defer will be executed
		case <-ctx.Done():
			return false
		}

		return false
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
//	@Summary      Get all incidents
//	@Description  get all incidents
//	@Tags         dish
//	@Accept       json
//	@Produce      json
//	@Success      200  {array}   dish.Incident
//	@Router       /dish/incidents [get]
func GetIncidentList(ctx *gin.Context) {
	core.PrintAllRootItems(ctx, CacheIncidents, pkgName)
	return
}

// PostNewIncidentByKey produces and broadcasts a new incident to be happening
//
//	@Summary      Add new incident
//	@Description  add new incident
//	@Tags         dish
//	@Accept       json
//	@Produce      json
//	@Param        key  query     dish.Incident true "socket body"
//	@Success      201  {array}   dish.Incident
//	@Failure      400  {object}  dish.Incident
//	@Failure      409  {object}  dish.Incident
//	@Failure      500  {object}  dish.Incident
//	@Router       /dish/incidents/{key} [post]
func PostNewIncidentByKey(ctx *gin.Context) {
	core.AddNewItemByParam(ctx, CacheIncidents, pkgName, Incident{})
	return
}

// UpdateIncidentByKey update requested incident
//
//	@Summary      Update incident by its key
//	@Description  update incident by its key
//	@Tags         dish
//	@Accept       json
//	@Produce      json
//	@Param        key  query     dish.Incident.ID true "socket body"
//	@Success      200  {array}   dish.Incident
//	@Failure      400  {object}  dish.Incident
//	@Failure      404  {object}  dish.Incident
//	@Failure      500  {object}  dish.Incident
//	@Router       /dish/incidents/{key} [put]
func UpdateIncidentByKey(ctx *gin.Context) {
	core.UpdateItemByParam(ctx, CacheIncidents, pkgName, Incident{})
	return
}

// DeleteIncidentByKey deletes given incident
//
//	@Summary      Delete incident by its key
//	@Description  delete incident by its key
//	@Tags         dish
//	@Accept       json
//	@Produce      json
//	@Param        key  query     dish.Incident.ID true "socket name"
//	@Success      200  {array}   dish.Incident
//	@Failure      404  {object}  dish.Incident
//	@Failure      500  {object}  dish.Incident
//	@Router       /dish/incidents/{key} [delete]
func DeleteIncidentByKey(ctx *gin.Context) {
	core.DeleteItemByParam(ctx, CacheIncidents, pkgName)
	return
}

// GetIncidentListBySocketID returns list of incidents for such Socket.ID
//
//	@Summary      Get incident list by socket ID
//	@Description  get incident list by socket ID
//	@Tags         dish
//	@Accept       json
//	@Produce      json
//	@Param        key  query  dish.Socket.ID true "socket name"
//	@Success      200  {array}   dish.Incident
//	@Failure      404  {object}  dish.Incident
//	@Router       /dish/incidents/{key} [get]
func GetIncidentListBySocketID(ctx *gin.Context) {
	var key string = ctx.Param("key")
	//var exportedIncidents = make(map[string]Incident)
	var exportedIncidents []Incident = []Incident{}
	var counter int = 0

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

  restoration

*/

// GetDishRoot returns all possible items of dish package
//
//	@Summary      Get all root items
//	@Description  get all root items
//	@Tags         dish
//	@Accept       json
//	@Produce      json
//	@Param        key  query     dish.Incident.ID true "socket name"
//	@Success      200  {object}   dish.Root
//	@Router       /dish [get]
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
//	@Summary      Restore dish package items
//	@Description  restore dish package items
//	@Tags         dish
//	@Accept       json
//	@Produce      json
//	@Success      201  {array}   dish.Incident
//	@Failure      404  {array}   dish.Incident
//	@Failure      500  {array}   dish.Incident
//	@Router       /dish/restore [post]
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
		CacheIncidents.Set(key, item)
		counter[0]++
	}

	for key, item := range importDish.Sockets {
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
