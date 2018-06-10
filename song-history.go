// WEQX Song History Web Scraper
// Parsing is heavily based on https://schier.co/blog/2015/04/26/a-simple-web-scraper-in-go.html
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


func keepLines(s string, n int) string {
	result := strings.Join(strings.Split(s, "\n")[:n], "\n")
	return strings.Replace(result, "\r", "", -1)
}

/* Checks if google.com is reachable. */
func sanityCheck() {
	resp, err := http.Get("https://google.com/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close();
}

const NOT_FOUND = "ERR: Attr not found"
func getAttr(tok html.Token, key string) (hasAttr bool, value string) {
	for _, attr := range tok.Attr {
		if attr.Key == key {
			return true, attr.Val
		}
	}

	return false, NOT_FOUND
}

func collectHistory(historyItems chan string, dateArg string, timeArg string) {
	resp, err := http.PostForm(SONG_HISTORY_URL, url.Values{DATE_KEY: {"06/01/2018"}, TIME_KEY: {"1:00pm"}})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tokType := tokenizer.Next()

		switch {
		case tokType == html.ErrorToken:
			// End of the document, we're done
			return
		case tokType == html.StartTagToken:
			tok := tokenizer.Token()

			hasClass, class := getAttr(tok, "class")
			if hasClass && strings.Contains(class, "songhistoryitem") {
				fmt.Println(tok)
			}


			// Check if the token is an <a> tag
			isAnchor := tok.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one
			// ok, url := getHref(t)
			// if !ok {
			// 	continue
			// }

			// // Make sure the url begines in http**
			// hasProto := strings.Index(url, "http") == 0
			// if hasProto {
			// 	ch <- url
			// }
		}
	}
}

func main() {
	sanityCheck()

	historyItems := make(chan string)

	// Produce song results, eventually for multiple dates concurrently
	go func() {
		collectHistory(historyItems, "06/01/2018", "1:00pm")
		close(historyItems)
	}()

	// Consume song results
	for historyItem := range historyItems {
		fmt.Println(historyItem)
	}
}




// hc := http.Client{}
//     req, err := http.NewRequest("POST", APIURL, nil)

//     form := url.Values{}
//     form.Add("ln", c.ln)
//     form.Add("ip", c.ip)
//     form.Add("ua", c.ua)
//     req.PostForm = form
//     req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

//     glog.Info("form was %v", form)
// 	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
