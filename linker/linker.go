package linker

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

// Link represents an HTML hyperlink extracted from a document.
// Href contains the value of the href attribute, and Text contains
// the visible text inside the <a> tag.
type Link struct {
	Href string
	Text string
}

// ParseHTML parses HTML from the provided io.Reader and extracts all
// hyperlinks. It walks the entire HTML tree, collects <a> nodes,
// converts them into Link structs, and returns the complete list.
// If parsing fails, ParseHTML returns an empty slice and error.
func ParseHTML(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	var links []Link
	if err != nil {
		fmt.Printf("Error while parsing the html file : %v\n", err)
		return links, err
	}
	linkedNodes := linkNodes(doc)
	buildLinks(&linkedNodes, &links)
	return links, nil
}

// buildLinks converts a list of <a> HTML nodes into Link structs.
// It extracts the href attribute and visible text from each node
// and appends the resulting Link values into the provided slice pointer.
func buildLinks(linkedNodes *[]*html.Node, links *[]Link) {
	for _, node := range *linkedNodes {
		var ret Link
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				ret.Href = attr.Val
				break
			}
		}
		ret.Text = getText(node)
		*links = append(*links, ret)
	}
}

// getText returns the concatenated visible text contained inside an
// HTML node. It recursively walks all descendant nodes and collects
// text from any TextNode it encounters.
func getText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	if node.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		ret += getText(c)
	}
	return ret
}

// linkNodes returns a slice of all <a> element nodes within the HTML
// tree rooted at the provided node. It performs a recursive depth-first
// search and collects nodes whose Data field equals "a".
func linkNodes(node *html.Node) []*html.Node {
	if node.Type == html.ElementNode && node.Data == "a" {
		return []*html.Node{node}
	}
	var ret []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}
