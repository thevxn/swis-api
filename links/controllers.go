package links

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func findLinkByHash(c *gin.Context) (*int, *Link) {
	for i, l := range links {
		if l.Hash == c.Param("hash") {
			return &i, &l
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "link not found",
	})
	return nil, nil
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
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, listing links",
		"links":   links,
	})
}

// @Summary Get link by :hash
// @Description get link by its :hash param
// @Tags links
// @Produce  json
// @Success 200 {object} links.Link
// @Router /links/{hash} [get]
// GetLinkByHash returns link's properties, given sent hash exists in database.
func GetLinkByHash(c *gin.Context) {
	if _, link := findLinkByHash(c); link != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "ok, dumping given link's details",
			"code":    http.StatusOK,
			"link":    link,
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

	links = append(links, newLink)

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "link added",
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

	links = importLinks.Links

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
	var updatedLink Link

	i, _ := findLinkByHash(c.Copy())
	updatedLink = links[*i]

	// inverse the Active field value
	updatedLink.Active = !updatedLink.Active

	links[*i] = updatedLink
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "link active toggle pressed!",
		"link":    updatedLink,
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

	i, _ := findLinkByHash(c.Copy())

	if err := c.BindJSON(&updatedLink); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}
	links[*i] = updatedLink

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
	i, l := findLinkByHash(c.Copy())

	// delete an element from the array
	// https://www.educative.io/answers/how-to-delete-an-element-from-an-array-in-golang
	newLength := 0
	for index := range links {
		if *i != index {
			links[newLength] = links[index]
			newLength++
		}
	}

	// reslice the array to remove extra index
	links = links[:newLength]

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "link deleted by Hash",
		"link":    *l,
	})
}
