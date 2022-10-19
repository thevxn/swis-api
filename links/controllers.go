package links

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var links = Links{
	l: make(map[string]Link),
}

func findLinkByHash(c *gin.Context) *Link {
	// TODO: nil map test!

	hash := c.Param("hash")

	links.RLock()
	defer links.RUnlock()
	if link, ok := links.l[hash]; ok {
		return &link
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "link not found",
	})
	return nil
}

// @Summary Get all links
// @Description get links complete list
// @Tags links
// @Produce  json
// @Success 200 {object} links.Link
// @Router /links [get]
// GetLinks GET method
// GetLinks returns JSON serialized list of links and their properties.
func GetLinks(c *gin.Context) {

	links.RLock()
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, listing links",
		"links":   links.l,
	})

	links.RUnlock()
	return
}

// @Summary Get link by :hash
// @Description get link by its :hash param
// @Tags links
// @Produce  json
// @Success 200 {object} links.Link
// @Router /links/{hash} [get]
// GetLinkByHash returns link's properties, given sent hash exists in database.
func GetLinkByHash(c *gin.Context) {
	if link := findLinkByHash(c); link != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "ok, dumping given link's details",
			"code":    http.StatusOK,
			"link":    *link,
		})
	}
	return
}

// @Summary Add new link to links
// @Description add new link to links array
// @Tags links
// @Produce json
// @Param request body links.Link true "query params"
// @Success 200 {object} links.Link
// @Router /links [post]
// PostNewLink enables one to add new link to links model.
func PostNewLink(c *gin.Context) {
	var newLink Link

	// bind received JSON to newLink
	if err := c.BindJSON(&newLink); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	// hotfix for new hash name
	hash := newLink.Name

	links.RLock()
	if _, found := links.l[hash]; found {
		// Link already exists, such hash/name is already used...
		c.IndentedJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "link hash name already used!",
			"hash":    hash,
		})
		links.RUnlock()
		return
	}
	links.RUnlock()

	// add newLink to the hash map
	links.Lock()
	links.l[hash] = newLink
	links.Unlock()

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "link added",
		"hash":    hash,
		"link":    newLink,
	})
}

// @Summary Upload links dump backup -- restores all links
// @Description update links JSON dump
// @Tags links
// @Accept json
// @Produce json
// @Router /links/restore [post]
// PostDumpRestore
func PostDumpRestore(c *gin.Context) {
	var importLinks Links

	if err := c.BindJSON(&importLinks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	links.Lock()
	links.l = importLinks.l
	links.Unlock()

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "links imported successfully",
	})
}

// (PUT /links/{hash}/active)
// @Summary Toggle active boolean for {hash}
// @Description toggle active boolean for {hash}
// @Tags links
// @Produce json
// @Param  id  path  string  true  "hash"
// @Success 200 {object} links.Link
// @Router /links/{hash}/active [put]
func ActiveToggleByHash(c *gin.Context) {
	//var updatedLink Link

	hash := c.Param("hash")
	link := findLinkByHash(c)

	// inverse the Active field value
	link.Active = !link.Active

	links.Lock()
	links.l[hash] = *link
	links.Unlock()

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "link active toggle pressed!",
		"link":    *link,
	})
	return
}

// (PUT /links/{hash})
// @Summary Update link by its Hash
// @Description update link by its Hash
// @Tags links
// @Produce json
// @Param request body links.Link.Hash true "query params"
// @Success 200 {object} links.Link
// @Router /links/{hash} [put]
func UpdateLinkByHash(c *gin.Context) {
	var updatedLink Link

	hash := c.Param("hash")
	_ = findLinkByHash(c)

	if err := c.BindJSON(&updatedLink); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	links.Lock()
	links.l[hash] = updatedLink
	links.Unlock()

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "socket updated",
		"link":    updatedLink,
	})
	return
}

// (DELETE /links/{})
// @Summary Delete link by its Hash
// @Description delete link by its Hash
// @Tags links
// @Produce json
// @Param  id  path  string  true  "link Hash"
// @Success 200 {object} links.Link
// @Router /links/{hash} [delete]
func DeleteLinkByHash(c *gin.Context) {

	hash := c.Param("hash")
	link := findLinkByHash(c)

	links.Lock()
	delete(links.l, hash)
	links.Unlock()

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "link deleted by Hash",
		"link":    *link,
	})
	return
}
