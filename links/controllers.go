package links

import (
	"net/http"

	"go.savla.dev/swis/v5/config"

	"github.com/gin-gonic/gin"
)

var Cache *config.Cache

// @Summary Get all links
// @Description get links complete list
// @Tags links
// @Produce json
// @Success 200 {object} links.Link
// @Router /links [get]
// GetLinks GET method
// GetLinks returns JSON serialized list of links and their properties.
func GetLinks(c *gin.Context) {
	var links = Cache.GetAll()

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
// @Produce json
// @Success 200 {object} links.Link
// @Router /links/{hash} [get]
// GetLinkByHash returns link's properties, given sent hash exists in database.
func GetLinkByHash(c *gin.Context) {
	var hash string = c.Param("hash")
	var link Link

	rawLink, ok := Cache.Get(hash)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "link not found",
		})
		return
	}

	link, ok = rawLink.(Link)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "cannot assert data type, database internal error",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, dumping requested link's details",
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
	var newLink Link

	if err := c.BindJSON(&newLink); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot bind input JSON stream",
		})
		return
	}

	var hash string = newLink.Name

	if _, found := Cache.Get(hash); found {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "link hash name already used!",
			"hash":    hash,
		})
		return
	}

	if saved := Cache.Set(newLink.Name, newLink); !saved {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "user couldn't be saved to database",
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "link added",
		"hash":    hash,
		"link":    newLink,
	})
	return
}

// @Summary Upload links dump backup -- restore all links
// @Description update links JSON dump
// @Tags links
// @Accept json
// @Produce json
// @Router /links/restore [post]
func PostDumpRestore(c *gin.Context) {
	var importLinks = &Links{}
	var counter int = 0

	if err := c.BindJSON(importLinks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err,
		})
		return
	}

	for _, link := range importLinks.Links {
		Cache.Set(link.Name, link)
		counter++
	}

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"count":   counter,
		"message": "links imported successfully, ommiting output",
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
	var hash string = c.Param("hash")
	var link Link

	rawLink, ok := Cache.Get(hash)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "link not found",
		})
		return
	}

	link, ok = rawLink.(Link)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "cannot assert data type, database internal error",
		})
		return
	}

	// inverse the Active field value
	link.Active = !link.Active

	if saved := Cache.Set(hash, link); !saved {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "link couldn't be saved to database",
		})
		return
	}

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
	var hash string = c.Param("hash")
	var updatedLink Link

	if _, ok := Cache.Get(hash); !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "link not found",
			"code":    http.StatusNotFound,
		})
		return
	}

	if err := c.BindJSON(&updatedLink); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot bind input JSON stream",
		})
		return
	}

	if saved := Cache.Set(hash, updatedLink); !saved {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "link couldn't be saved to database",
		})
		return
	}

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

	if _, ok := Cache.Get(hash); !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "link not found",
		})
		return
	}

	if deleted := Cache.Delete(hash); !deleted {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "link couldn't be deleted from database",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "link deleted by Hash",
		"hash":    hash,
	})
	return
}
