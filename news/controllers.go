package news

import (
	//b64 "encoding/base64"
	"encoding/xml"
	"log"
	"net/http"
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

// GetSources
func GetSources(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "ok, dumping news sources",
		"code":    http.StatusOK,
		"news":    news,
	})
}
