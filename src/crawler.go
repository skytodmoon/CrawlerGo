package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/jackdanger/collectlinks"
)

var visited = make(map[string]bool)

func main() {
	//url := "http://www.npc.gov.cn/npc/c12488/list.shtml"
	url := "http://www.baidu.com"

	queue := make(chan string)
	go func() {
		queue <- url
	}()
	for uri := range queue {
		download(uri, queue)

	}
}

func download(url string, queue chan string) {
	visited[url] = true
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	//header
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error", err)
		return
	}
	//close after finish
	defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println("read error", err)
	// 	return
	// }
	// fmt.Println(string(body))
	links := collectlinks.All(resp.Body)
	for _, link := range links {
		absolute := urlJoin(link, url)
		if url != " " && !visited[absolute] && !checkWithBlackList(absolute) {
			fmt.Println("parse url", absolute)
			go func() {
				queue <- absolute
			}()
		}

	}
}

func urlJoin(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return " "
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return " "
	}
	return baseURL.ResolveReference(uri).String()
}

func checkWithBlackList(url string) bool {
	check := strings.Contains(url, "javascript") || strings.Contains(url, "download") || strings.Contains(url, "mail") || len(url) < 50
	//fmt.Println("contain invalid:", check)
	return check
}
