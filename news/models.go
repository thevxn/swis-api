package news

import "time"

// Sources

type News struct {
	User    string   `json:"news_user"`
	Sources []Source `json:"news_sources"`
}

type Source struct {
	ID   string `json:"source_id"`
	Name string `json:"source_name"`
	URL  string `json:"source_url"`
}

var sources = []Source{
	//{Name: "Aktuálně.cz", URL: "https://www.aktualne.cz/rss/"},
	{Name: "ČT24 Hlavní zprávy", URL: "http://www.ceskatelevize.cz/ct24/rss/hlavni-zpravy"},
	{Name: "iRozhlas.cz", URL: "https://www.irozhlas.cz/rss/irozhlas"},
	{Name: "Seznam Zprávy", URL: "https://api.seznamzpravy.cz/v1/documenttimelines/5ac49a0272c43201ee1d957f?rss=1"},
	//{Name: "Root.cz Zprávičky", URL: "https://www.root.cz/rss/zpravicky/"},
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
	Title     string    `xml:"title" json:"title"`
	Perex     string    `xml:"description" json:"perex"`
	Link      string    `xml:"link" json:"link"`
	Server    string    `json:"server"`
	PubDate   string    `xml:"pubDate" json:"pub_date"`
	ParseDate time.Time `json:"parse_date_rfc1123z"`
}

// XML exported Channel
type Channel struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Desc  string `xml:"description"`
	Lang  string `xml:"language"`
	Items []Item `xml:"item"`
}

// XML exported root
type Rss struct {
	Channel Channel `xml:"channel"`
}
