package worldlogo

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

type WorldLogo struct {
	Name string
	Key  string
	Src  string
}

func TakeWorldLogo(ctx context.Context, srcUrl string) (res []WorldLogo, err error) {
	// load html page
	node, err := loadHtmlPage(ctx, srcUrl)
	if err != nil {
		return nil, fmt.Errorf("loadHtmlPage: %w", err)
	}

	// parse page
	return extractWorldLogoFromPage(node), nil
}

func loadHtmlPage(ctx context.Context, url string) (*html.Node, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response code err: %d", resp.StatusCode)
	}

	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("html.Parse: %w", err)
	}
	return node, nil
}

func extractWorldLogoFromPage(node *html.Node) (res []WorldLogo) {
	if node == nil {
		return nil
	}

	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "class" && attr.Val == "logo" {
				item := extractWorldLogoFromDiv(node)
				if item != nil {
					res = append(res, *item)
				}
				break
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		items := extractWorldLogoFromPage(child)
		if len(items) > 0 {
			res = append(res, items...)
		}
	}

	return res
}

func extractWorldLogoFromDiv(node *html.Node) *WorldLogo {
	var res WorldLogo

	// get key from root node
	if node == nil {
		return nil
	}
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			res.Key = attr.Val
			break
		}
	}

	// get link
	src := node.FirstChild.FirstChild.FirstChild
	for _, attr := range src.Attr {
		if attr.Key == "src" {
			res.Src = attr.Val
			break
		}
	}

	// get name
	name := node.FirstChild.FirstChild.NextSibling.FirstChild
	res.Name = name.Data

	return &res
}
