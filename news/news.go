package news

import (
	//b64 "encoding/base64"
	"encoding/xml"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)


type News struct {
	User	string		`json:"news_user"`
	Sources []Source 	`json:"news_sources"`
}

type Source struct {
	ID	string	`json:"source_id"` 
	Name	string	`json:"source_name"`
	URL	string	`json:"source_url"`
}

var sources = []Source{
	{Name: "Aktuálně.cz", URL: "https://www.aktualne.cz/rss/"},
	{Name: "ČT24 Hlavní zprávy", URL: "http://www.ceskatelevize.cz/ct24/rss/hlavni-zpravy"},
	{Name: "iRozhlas.cz", URL: "https://www.irozhlas.cz/rss/irozhlas"},
	{Name: "Seznam Zprávy", URL: "https://api.seznamzpravy.cz/v1/documenttimelines/5ac49a0272c43201ee1d957f?rss=1"},
	{Name: "Root.cz Zprávičky", URL: "https://www.root.cz/rss/zpravicky/"},
}

var news = []News{
	{User: "krusty", Sources: sources},
}


// typical RSS structure:
// <rss ...>
//   <channel>
//     ...
//     <item>
//       <title> ...
//       ...

// XML exported Item
type Item struct {
	Title		string	`xml:"title" json:"title"`
	Perex		string	`xml:"description" json:"perex"`
	//Server	string	
	Link		string	`xml:"link" json:"link"`
	PubDate		string	`xml:"pubDate" json:"timestamp"`
}

// XML exported Channel
type Channel struct {
	Title 	string	`xml:"title"`
	Link 	string	`xml:"link"`
	Desc	string	`xml:"description"`
	Lang	string	`xml:"language"`
	Items	[]Item	`xml:"item"`
}

// XML exported root
type Rss struct {
	Channel Channel `xml:"channel"`
}


func findSourcesByUser(c *gin.Context) (s *[]Source) {
	for _, n := range news {
		if n.User == c.Param("user") {
			//c.IndentedJSON(http.StatusOK, a)
			return &n.Sources
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code": http.StatusNotFound,
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
	err = decoder.Decode(&rss)
	if err != nil {
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
		//R := *fetchRSSContents(&s)
		for _, item := range *fetchRSSContents(&s) {
			items = append(items, item)
		}
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"news": items,
	})
}

// GetSources
func GetSources(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"news": news,
	})
}

/*
// GetUserByID returns user's properties, given sent ID exists in database.
func GetUserByID(c *gin.Context) {
	//id := c.Param("id")

	if _, user := findUserByID(c); user != nil {
		// user found
		c.IndentedJSON(http.StatusOK, user)
	}

	//c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

// PostUser enables one to add new user to users model.
func PostUser(c *gin.Context) {
	var newUser User

	// bind received JSON to newUser
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	// add new user
	users = append(users, newUser)

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated,
		"message": "user added",
		"user": newUser,
	})
}
*/
