package news

import (
	//b64 "encoding/base64"
	"encoding/xml"
	"log"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

func findSourcesByUser(c *gin.Context) (s *[]Source) {
	for _, n := range news {
		if n.User == c.Param("user") {
			//c.IndentedJSON(http.StatusOK, a)
			return &n.Sources
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "user's sources not found",
	})
	return nil
}

func fetchRSSContents(s *Source) (i *[]Item) {
	resp, err := http.Get(s.URL)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()

	var rss = Rss{}

	decoder := xml.NewDecoder(resp.Body)
	if err = decoder.Decode(&rss); err != nil {
		log.Println(err)
		return nil
	}

	//log.Printf("Channel title: %v\n", rss.Channel.Title)
	//log.Printf("Channel link: %v\n", rss.Channel.Link)

	/*
		for i, item := range rss.Channel.Items {
			log.Printf("%v. item title: %v\n", i, item.Title)
		}*/
	return &rss.Channel.Items
}

// @Summary Get news by User
// @Description fetch and parse news for :user param
// @Tags news
// @Produce  json
// @Success 200 {object} news.Item
// @Router /news/{name} [get]
// GetNewsByUser returns all possible news from all sources loaded in memory
func GetNewsByUser(c *gin.Context) {
	userSources := findSourcesByUser(c)
	if userSources == nil {
		return
	}

	//var R = []Rss{}
	var items = []Item{}

	for _, s := range *userSources {
		cont := fetchRSSContents(&s)
		if cont == nil {
			continue
		}

		for _, item := range *cont {
			// time layouts (date template constants) --> https://go.dev/src/time/format.go
			item.ParseDate, _ = time.Parse(time.RFC1123Z, item.PubDate)

			// convert news link to server/hostname
			u, _ := url.Parse(item.Link)
			item.Server = string(u.Hostname())

			items = append(items, item)
		}
	}

	// sort items by date DESC
	// https://stackoverflow.com/a/47028486
	sort.Slice(items, func(i, j int) bool {
		return items[i].ParseDate.After(items[j].ParseDate)
	})

	c.IndentedJSON(http.StatusOK, gin.H{
		"news": items,
	})
}

// @Summary Get news source list
// @Description get all news sources
// @Tags news
// @Produce  json
// @Success 200 {object} news.News.Sources
// @Router /news/sources/ [get]
// GetSources
func GetSources(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "ok, dumping news sources",
		"code":    http.StatusOK,
		"sources": sources,
	})
}

// @Summary Get news source list by Username
// @Description get news sources by their :name param
// @Tags news
// @Produce  json
// @Success 200 {object} news.News.Sources
// @Router /news/sources/{name} [get]
// GetSources
func GetSourcesByUser(c *gin.Context) {
	userSources := findSourcesByUser(c)
	if userSources == nil {
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "ok, dumping news sources",
		"code":    http.StatusOK,
		"sources": *userSources,
	})
}

// @Summary Upload news sources dump backup -- restores all sources
// @Description update news sources JSON dump
// @Tags news
// @Accept json
// @Produce json
// @Router /news/sources/restore [post]
// PostDumpRestore
func PostDumpRestore(c *gin.Context) {
	var importSources []News //News.Sources

	// bind received JSON to newUser
	if err := c.BindJSON(&importSources); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	// restore sources
	news = importSources

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "sources imported successfully",
	})
}
