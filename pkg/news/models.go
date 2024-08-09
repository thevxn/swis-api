package news

import (
	"time"
)

//
// Sources
//

type NewsSources struct {
	UserSources []UserSource `json:"user_sources"`
}

type UserSource struct {
	// News sources linked username.
	User string `json:"user_name" required:"true"`

	// Actual user's news sources.
	Sources []Source `json:"news_sources"`
}

type Source struct {
	// Source unique identificator.
	ID string `json:"source_id" required:"true"`

	// Source's more natural name.
	Name string `json:"source_name"`

	// RSS URL to the source.
	URL string `json:"source_url"`
}

//
// fetched news from RSS channel
//

// typical RSS structure:
// <rss ...>
//   <channel>
//     ...
//     <item>
//       <title> ... </title>
//       ...

// XML exported root
type Rss struct {
	// XML exported Channel.
	Channel Channel `xml:"channel"`
}

type Channel struct {
	// Channel title.
	Title string `xml:"title"`

	// Link to such channel.
	Link string `xml:"link"`

	// Channel description.
	Desc string `xml:"description"`

	// Channel language.
	Lang string `xml:"language"`

	// XML exported Item(s) -- actual news items.
	Items []Item `xml:"item"`
}

type Item struct {
	// Item's title (headline).
	Title string `xml:"title" json:"title"`

	// Item's short description (perex).
	Perex string `xml:"description" json:"perex"`

	// Link to such item -- to the actual article usually.
	Link string `xml:"link" json:"link"`

	// Issuer server name (hostname with subdomain).
	Server string `json:"server"`

	// Date of issue, formatted by issuer.
	PubDate string `xml:"pubDate" json:"pub_date"`

	// Special date of issue format to reparse/order all news items altogether.
	ParseDate time.Time `json:"parse_date_rfc1123z"`
}
