package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"siteMapBuilder/linker"
	"siteMapBuilder/queue"
	"strings"
)

type LinkNode struct {
	Link  string
	Level int
}

// traversalDFS performs a depth-first search traversal of URLs starting from urlStr.
// It explores links up to the specified depth, marking visited URLs to avoid cycles.
// Results are appended to traversedLinks in DFS order.
func traversalDFS(urlStr string, depth int, visited map[string]bool, traversedLinks *[]string) {
	if depth == 0 {
		return
	}
	visited[urlStr] = true
	*traversedLinks = append(*traversedLinks, urlStr)
	for _, link := range getData(urlStr) {
		if !visited[link] {
			traversalDFS(link, depth-1, visited, traversedLinks)
		}
	}
}

// traversalBFS performs a breadth-first search traversal of URLs starting from urlStr.
// It explores links level by level up to maxDepth, marking visited URLs to avoid cycles.
// Results are appended to traversedLinks in BFS order (level by level).
func traversalBFS(urlStr string, maxDepth int, traversedLinks *[]string) {
	var q queue.Queue[LinkNode]
	visited := make(map[string]bool)

	q.Push(LinkNode{Link: urlStr, Level: 0})
	visited[urlStr] = true

	for !q.IsEmpty() {
		currLink, _ := q.Pop()
		*traversedLinks = append(*traversedLinks, currLink.Link)
		for _, childLink := range getData(currLink.Link) {
			if !visited[childLink] && currLink.Level < maxDepth {
				q.Push(LinkNode{Link: childLink, Level: currLink.Level + 1})
				visited[childLink] = true
			}
		}
	}
}

// getData fetches and parses all links from the given URL.
// It makes an HTTP GET request, extracts all <a> tag hrefs, converts them to absolute URLs,
// and filters to only return links with the same domain as the original URL.
// Returns an empty slice if the HTTP request fails.
//
// Parameters:
//   - urlStr: The URL to fetch and parse
//
// Returns:
//   - A slice of absolute URLs from the same domain
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
	urlFlag := flag.String("url", "https://gophercises.com", "The URL you want to build your siteMap for")
	depthFlag := flag.Int("depth", 10, "Maximum Number of links deep to traverse")
	traversalMethod := flag.String("traversal", "bfs", "Traversal method for parsing links: dfs/bfs")
	flag.Parse()
	var traversedLinks []string
	switch {
	case strings.ToLower(*traversalMethod) == "dfs":
		visited := make(map[string]bool)
		traversalDFS(*urlFlag, *depthFlag, visited, &traversedLinks)
	default:
		traversalBFS(*urlFlag, *depthFlag, &traversedLinks)
	}
	for _, link := range traversedLinks {
		fmt.Println(link)
	}
}
