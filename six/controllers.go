package six

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// (GET /six)
// @Summary Get the six struct
// @Description get the six struct
// @Tags six
// @Produce  json
// @Success 200 {object} six.SixStruct
// @Router /six [get]
func GetSixStruct(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":       http.StatusOK,
		"message":    "ok, dumping the six struct",
		"calendars":  sixStruct.Calendars,
		"todo_lists": sixStruct.TodoLists,
	})
}

// @Summary
// @Description
// @Tags six
// @Produce json
// @Success 200 {object} six.Calendar
// @Router /six/calendar/{owner_name} [get]
// GetCalendarByUser
func GetCalendarByUser(c *gin.Context) {
	var userCalendar Calendar

	for _, cal := range sixStruct.Calendars {
		if cal.Owner == c.Param("owner_name") {
			userCalendar = cal
			break
		}

		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "user's calendar not found",
			"code":    http.StatusNotFound,
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message":  "ok, showing user's calendar items",
		"code":     http.StatusOK,
		"calendar": userCalendar.Items,
	})
}

// (POST /six/restore)
// @Summary Upload six dump backup -- restores all loaded calendars and todo lists
// @Description upload six JSON dump
// @Tags six
// @Accept json
// @Produce json
// @Router /six/restore [post]
// restore all six structs from JSON dump (JSON-bind)
func PostDumpRestore(c *gin.Context) {
	var importSix SixStruct

	if err := c.BindJSON(&importSix); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	sixStruct = importSix

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "six struct imported, omitting output",
	})
}
