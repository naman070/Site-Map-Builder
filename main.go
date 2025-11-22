package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"siteMapBuilder/linker"
	"strings"
)

func getData(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Printf("Error in getting data for url: %s, error: %v\n", urlStr, err)
		return []string{}
	}
	defer resp.Body.Close()
	reqUrl := resp.Request.URL // this will adjust the redirected scheme.
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	baseUrlStr := baseUrl.String()
	return filter(parseHtmlTags(resp.Body, baseUrlStr), withPrefix(baseUrlStr))
}

// parseHtmlTags extracts all <a> tag href values, converting relative
// paths into absolute URLs based on baseUrlStr.
func parseHtmlTags(r io.Reader, baseUrlStr string) []string {
	links, _ := linker.ParseHTML(r)
	var ret []string
	for _, link := range links {
		if strings.HasPrefix(link.Href, "/") {
			ret = append(ret, baseUrlStr+link.Href)
		} else if strings.HasPrefix(link.Href, "http") {
			ret = append(ret, link.Href)
		}
	}
	return ret
}

// filter returns a new slice containing only the URLs that satisfy
// the given predicate function.
func filter(urls []string, filterConditionFunc func(string) bool) []string {
	var ret []string
	for _, link := range urls {
		if filterConditionFunc(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

// withPrefix returns a predicate that checks whether a URL
// begins with the provided prefix (used for domain filtering).
func withPrefix(prefix string) func(string) bool {
	return func(urlStr string) bool {
		return strings.HasPrefix(urlStr, prefix)
	}
}

func main() {
	urlFlag := flag.String("url", "https://gophercises.com/cyoa", "The URL you want to build your siteMap for")
	flag.Parse()
	links := getData(*urlFlag)
	for _, link := range links {
		fmt.Println(link)
	}
}
