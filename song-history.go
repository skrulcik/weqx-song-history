// WEQX Song History Web Scraper
// Parsing is heavily based on https://schier.co/blog/2015/04/26/a-simple-web-scraper-in-go.html
// Inferred request/response structure:
//
// Request
//     playlistdate: mm/dd/yyyy date
//     playlisttime: hh:MM[am,pm] time
// Response HTML
//     div.songhistoryresult has a title attribute with the song title and
//         artist
//
// Note that the response HTML contains more information
package main

import (
	"net/http"
	"net/url"
)

import "golang.org/x/net/html"
import "fmt"
import "strings"


const SONG_HISTORY_URL = "http://www.weqx.com/song-history"
const DATE_KEY = "playlistdate"
const TIME_KEY = "playlisttime"
const TOP_5 = 5 // # of results to look for after 5PM
const NOT_FOUND = "ERR: Attr not found"

// Retrieves the value of the given attribute for the given token.
// If the attribute exists, getAttr returns (true, <attribute value>)
// If it does not exist, getAttr returns (false, <undefined>)
func getAttr(tok html.Token, attribute string) (hasAttr bool, value string) {
	for _, attr := range tok.Attr {
		if attr.Key == attribute {
			return true, attr.Val
		}
	}

	return false, NOT_FOUND
}


func tryWithAttribute(tok html.Token, attribute string, f func(val string)) {
	hasAttr, val := getAttr(tok, attribute)
	if hasAttr {
		f(val)
	}
}


func collectHistory(dateArg string, timeArg string) {

	resp, err := http.PostForm(SONG_HISTORY_URL, url.Values{DATE_KEY: {"06/01/2018"}, TIME_KEY: {"1:00pm"}})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	var found int
	for found = 0; found < TOP_5; {
		tokType := tokenizer.Next()

		if tokType == html.ErrorToken {
			fmt.Errorf("%s", tokenizer.Token())
			continue
		}

		switch {
		case tokType == html.ErrorToken:
			// End of the document, we're done
			return
		case tokType == html.StartTagToken:
			tok := tokenizer.Token()

			tryWithAttribute(tok, "class", func(classes string) {
				if strings.Contains(classes, "songhistoryitem") {
					found++
					tryWithAttribute(tok, "title", func(songTitle string) {
						fmt.Println(songTitle)
					})
				}
			})
		}
	}
	if found != TOP_5 {
		fmt.Println("found was only %d, expected %d", found, TOP_5)
	}
}

func main() {
	collectHistory("06/01/2018", "1:00pm")
}
