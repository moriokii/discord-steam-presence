package parser

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func Parse(accountid int) (name string, rich string, logo string) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://steamcommunity.com/miniprofile/%d?appid=undefined", accountid), nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		panic(err)
	}

	return doc.Find(".miniprofile_game_name").Text(), doc.Find(".rich_presence").Text(), doc.Find(".game_logo").AttrOr("src", "none")
}
