package main

import (
	"net/http"
	"net/url"
)

import "io/ioutil"
import "fmt"
import "strings"

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

const SONG_HISTORY_URL = "http://www.weqx.com/song-history"
const DATE_KEY = "playlistdate"
const TIME_KEY = "playlisttime"

func main() {
	sanityCheck()

	resp, err := http.PostForm(SONG_HISTORY_URL, url.Values{DATE_KEY: {"06/01/2018"}, TIME_KEY: {"1:00pm"}})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
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
