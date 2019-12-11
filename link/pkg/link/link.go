package link

import (
	"golang.org/x/net/html"
	"io"
)

type Link struct {
	Href string
	Text string
}

func Parse(reader io.Reader) ([]Link, error) {
	node, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	nodes := linkNodes(node)
	var links []Link
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}

	return links, nil
}

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

func buildLink(node *html.Node) Link {
	var ret Link
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}

	ret.Text = text(node)

	return ret
}

func text(node *html.Node) string {

	if node.Type == html.TextNode {
		return node.Data
	}

	if node.Type != html.ElementNode {
		return ""
	}

	var ret string
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		ret += text(c)
	}

	return ret
}
