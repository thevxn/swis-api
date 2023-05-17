package links

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var l sync.Map

// @Summary Get all links
// @Description get links complete list
// @Tags links
// @Produce  json
// @Success 200 {object} links.Link
// @Router /links [get]
// GetLinks GET method
// GetLinks returns JSON serialized list of links and their properties.
func GetLinks(c *gin.Context) {
	var links = make(map[string]Link)

	l.Range(func(rawKey, rawVal interface{}) bool {
		// very insecure assert
		k, ok := rawKey.(string)
		v, ok := rawVal.(Link)

		if !ok {
			return false
		}

		links[k] = v
		return true
	})

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, listing links",
		"links":   links,
	})
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
	var hash string = c.Param("hash")
	var link Link

	rawLink, ok := l.Load(hash)
	link, ok = rawLink.(Link)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "link not found",
			"code":    http.StatusNotFound,
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "ok, dumping given link's details",
		"code":    http.StatusOK,
		"link":    link,
	})
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
	var newLink = &Link{}

	// bind received JSON to newLink
	if err := c.BindJSON(newLink); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	var hash string = newLink.Name

	if _, found := l.Load(hash); found {
		// Link already exists, such hash/name is already used...
		c.IndentedJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "link hash name already used!",
			"hash":    hash,
		})
		return
	}

	l.Store(newLink.Name, newLink)

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "link added",
		"hash":    hash,
		"link":    newLink,
	})
	return
}

// @Summary Upload links dump backup -- restores all links
// @Description update links JSON dump
// @Tags links
// @Accept json
// @Produce json
// @Router /links/restore [post]
// PostDumpRestore
func PostDumpRestore(c *gin.Context) {
	var importLinks = &Links{}

	if err := c.BindJSON(importLinks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err,
		})
		return
	}

	for _, link := range importLinks.Links {
		l.Store(link.Name, link)
	}

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "links imported successfully",
		"links":   l,
	})
	return
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
	var link Link
	var hash string = c.Param("hash")

	rawLink, ok := l.Load(hash)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "links not found",
			"code":    http.StatusNotFound,
		})
		return
	}

	link, typeOk := rawLink.(Link)
	if !typeOk {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"message": "stored value is not type Link",
			"code":    http.StatusConflict,
		})
		return
	}

	// inverse the Active field value
	link.Active = !link.Active

	l.Store(hash, link)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "link active toggle pressed!",
		"link":    link,
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
	var hash string = c.Param("hash")

	if _, ok := l.Load(hash); !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "links not found",
			"code":    http.StatusNotFound,
		})
		return
	}

	if err := c.BindJSON(&updatedLink); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	l.Store(hash, updatedLink)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "socket updated",
		"link":    updatedLink,
	})
	return
}

// @Summary Delete link by its Hash
// @Description delete link by its Hash
// @Tags links
// @Produce json
// @Param  id  path  string  true  "link Hash"
// @Success 200 {object} links.Link
// @Router /links/{hash} [delete]
func DeleteLinkByHash(c *gin.Context) {
	var hash string = c.Param("hash")

	l.Delete(hash)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "link deleted by Hash",
		"hash":    hash,
	})
	return
}
