package six

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func findCalendarByUser(c *gin.Context) (*int, *Calendar) {
	for idx, cal := range sixStruct.Calendars {
		if cal.Owner == c.Param("owner_name") {
			return &idx, &cal
			break
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{
		"message": "user's calendar not found",
		"code":    http.StatusNotFound,
	})
	return nil, nil
}

func findCalendarItemByName(c *gin.Context, cal Calendar) (*int, *Item) {
	if &cal == nil {
		return nil, nil
	}
	items := cal.Items

	for itemIdx, item := range items {
		// catching nil pointer dereference exception...
		// this panics the server when Name is empty in memory!
		if item.Name == c.Param("item_name") {
			return &itemIdx, &item
			break
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{
		"message": "item not found in user's calendar",
		"code":    http.StatusNotFound,
	})
	return nil, nil
}

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
	_, userCalendar := findCalendarByUser(c)
	if userCalendar == nil {
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message":  "ok, showing user's calendar items",
		"code":     http.StatusOK,
		"calendar": userCalendar.Items,
	})
}

// @Summary Add new item to user's calendar
// @Description add new item to user's calendar
// @Tags six
// @Produce json
// @Param request body six.Item true "six.Item"
// @Success 200 {object} six.Item
// @Router /six/calendar/{owner_name} [post]
func PostCalendarItemByUser(c *gin.Context) {
	var newItem Item

	// bind received JSON to newItem
	if err := c.BindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	// find the right calendar
	calIdx, cal := findCalendarByUser(c)
	if calIdx == nil || cal == nil {
		return
	}

	// check for already existing item
	/*_, item := findCalendarItemByName(c, *cal)
	if item != nil {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"message": "this item name already exists!",
			"code":    http.StatusConflict,
			"item":    newItem,
		})
		return
	}*/

	// add item to calendar
	cal.Items = append(cal.Items, newItem)

	// update calendar
	sixStruct.Calendars[*calIdx] = *cal

	// HTTP 200 OK
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "item added",
		"item":    newItem,
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

// (DELETE /six/calendar/{owner_name}/item/{item_name})
// @Summary Delete calendar item by its name
// @Description delete calendar item by its name
// @Tags six
// @Produce json
// @Param  id  path  string  true  "item_name"
// @Success 200 {object} six.Item
// @Router /six/calendar/{owner_name}/item/{item_name} [delete]
func DeleteCalendarItemNameByUser(c *gin.Context) {

	// find the right calendar
	calIdx, cal := findCalendarByUser(c)
	if calIdx == nil || cal == nil {
		return
	}

	items := cal.Items

	// find the right calendar item -- TODO: this function repeats calendar searching (the very previous paragraph of this method)!
	itemIdx, item := findCalendarItemByName(c, *cal)
	if itemIdx == nil || item == nil {
		return
	}

	// delete an element from the array
	// https://www.educative.io/answers/how-to-delete-an-element-from-an-array-in-golang
	newLength := 0
	for index := range items {
		if *itemIdx != index {
			items[newLength] = items[index]
			newLength++
		}
	}

	// reslice the array to remove extra index
	items = items[:newLength]

	// add item to calendar
	cal.Items = items

	// update calendar
	sixStruct.Calendars[*calIdx] = *cal

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "socket deleted by ID",
		"item":    *item,
	})
}
