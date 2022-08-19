package link_parser

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"strings"
)

type Link struct {
	Href string
	Text string
}

type Attributes []html.Attribute

// Parse parse all links elements (<a>) from the given html reader
func Parse(htmlReader io.Reader) ([]Link, error) {
	node, err := html.Parse(htmlReader)

	if err != nil {
		return nil, err
	}

	var nodes []html.Node
	var parsedLinks []Link
	FindNodes(*node, atom.A, &nodes)

	for _, el := range nodes {
		parsedLinks = append(parsedLinks, Link{
			Href: FindAttributeAndGetValue(el.Attr, func(attr html.Attribute) bool {
				return attr.Key == "href"
			}),
			Text: ParseTextFromTheNode(el),
		})

	}

	return parsedLinks, nil
}

// FindNodes find all nodes by the following type and save them into nodes slice
func FindNodes(parentNode html.Node, nodeType atom.Atom, nodes *[]html.Node) {
	if parentNode.DataAtom == nodeType {
		*nodes = append(*nodes, parentNode)
	} else {
		if parentNode.FirstChild != nil {
			FindNodes(*parentNode.FirstChild, nodeType, nodes)
		}
	}

	if parentNode.NextSibling != nil {
		FindNodes(*parentNode.NextSibling, nodeType, nodes)
	}
}

func extractText(node html.Node) string {
	var text string

	if node.Type == html.TextNode {
		text = node.Data
	} else if node.FirstChild != nil {
		text += extractText(*node.FirstChild)
	}

	if node.NextSibling != nil {
		text += extractText(*node.NextSibling)
	}

	return text
}

// ParseTextFromTheNode parse all text nodes from any html node. Unescaped html tags and comments are omitted
func ParseTextFromTheNode(node html.Node) string {
	if node.FirstChild == nil {
		return ""
	}

	return strings.Join(strings.Fields(extractText(*node.FirstChild)), " ")
}

// FindAttributeAndGetValue returns attribute value by the following predicate. If noting found then returns empty string
func FindAttributeAndGetValue(a Attributes, searchFunc func(attribute html.Attribute) bool) string {
	for _, attribute := range a {
		if searchFunc(attribute) {
			return attribute.Val
		}
	}

	return ""
}
