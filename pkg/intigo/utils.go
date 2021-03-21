package intitools

import (
	"fmt"

	"golang.org/x/net/html"
)

func (c *Client) getElementByName(name string, n *html.Node) (element *html.Node, ok bool) {
	for _, a := range n.Attr {
		if a.Key == "name" && a.Val == name {
			return n, true
		}
	}
	for m := n.FirstChild; m != nil; m = m.NextSibling {
		if element, ok = c.getElementByName(name, m); ok {
			return
		}
	}
	return
}

func (c *Client) getElementValue(name string, n *html.Node) (string, error) {
	element, ok := c.getElementByName(name, n)
	if !ok {
		return "", fmt.Errorf("Cannot find element %s", name)
	}
	for _, a := range element.Attr {
		if a.Key == "value" {
			return a.Val, nil
		}
	}

	return "", fmt.Errorf("Cannot find value of element %s", name)

}
