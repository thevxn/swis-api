package news

import (
	"encoding/xml"
	"net/http"
)

func fetchRSSContents(s Source) (i *[]Item) {
	resp, err := http.Get(s.URL)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var rss = Rss{}

	decoder := xml.NewDecoder(resp.Body)
	if err = decoder.Decode(&rss); err != nil {
		//log.Println(err)
		return nil
	}

	//log.Printf("Channel title: %v\n", rss.Channel.Title)
	//log.Printf("Channel link: %v\n", rss.Channel.Link)

	/*for i, item := range rss.Channel.Items {
		log.Printf("%v. item title: %v\n", i, item.Title)
	}*/
	return &rss.Channel.Items
}
