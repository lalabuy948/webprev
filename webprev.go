package webprev

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type (
	//WebCard contains title, description and image url.
	WebCard struct {
		Title       string
		Description string
		ImgURL      string
	}
	//WebPreview contains 3 web preview cards - Generic/default, OpenGraph/Facebook and Twitter.
	WebPreview struct {
		Generic   WebCard
		OpenGraph WebCard
		Twitter   WebCard
	}
)

//Preview makes request by given url, parse html amd return WebPreview struct.
func Preview(url string) (webPreview WebPreview, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	webPreview = parseHtml(url, resp.Body)
	err = resp.Body.Close()
	return
}

func parseHtml(url string, body io.ReadCloser) (webPreview WebPreview) {
	tokenizer := html.NewTokenizer(body)
	for {
		tokenType := tokenizer.Next()
		switch {
		case tokenType == html.ErrorToken:
			if webPreview.Generic.ImgURL == "" {
				webPreview.Generic.ImgURL = webPreview.OpenGraph.ImgURL
				if webPreview.Generic.ImgURL == "" {
					webPreview.Generic.ImgURL = webPreview.Twitter.ImgURL
				}
			}
			return
		case tokenType == html.StartTagToken:
			token := tokenizer.Token()
			parseMetaTags(url, token, &webPreview)
			if token.Data == "title" && webPreview.Generic.Title == "" {
				tt := tokenizer.Next()
				if tt == html.TextToken {
					token := tokenizer.Token()
					webPreview.Generic.Title = token.Data
				}
			}
		}
	}
}

func parseMetaTags(url string, token html.Token, webPreview *WebPreview) {
	if token.Data == "meta" {
		attr := token.Attr
		for i := 0; i < len(attr); i++ {
			if i > 0 {
				if isNameKey(attr[i-1]) && attr[i-1].Val == "title" && isContentKey(attr[i]) {
					if webPreview.Generic.Title == "" {
						webPreview.Generic.Title = attr[i].Val
					}
				}
				if isNameKey(attr[i]) && attr[i].Val == "title" && isContentKey(attr[i-1]) {
					if webPreview.Generic.Title == "" {
						webPreview.Generic.Title = attr[i-1].Val
					}
				}

				if isNameKey(attr[i-1]) && attr[i-1].Val == "description" && isContentKey(attr[i]) {
					if webPreview.Generic.Description == "" {
						webPreview.Generic.Description = attr[i].Val
					}
				}
				if isNameKey(attr[i]) && attr[i].Val == "description" && isContentKey(attr[i-1]) {
					if webPreview.Generic.Description == "" {
						webPreview.Generic.Description = attr[i-1].Val
					}
				}

				if attr[i].Key == "itemprop" && attr[i].Val == "image" && isContentKey(attr[i-1]) {
					if webPreview.Generic.ImgURL == "" {
						webPreview.Generic.ImgURL = supplementImgURL(url, attr[i-1].Val)
					}
				}
				if attr[i-1].Key == "itemprop" && attr[i-1].Val == "image" && isContentKey(attr[i]) {
					if webPreview.Generic.ImgURL == "" {
						webPreview.Generic.ImgURL = supplementImgURL(url, attr[i].Val)
					}
				}

				if isPropertyKey(attr[i-1]) && attr[i-1].Val == "og:title" && isContentKey(attr[i]) {
					if webPreview.OpenGraph.Title == "" {
						webPreview.OpenGraph.Title = attr[i].Val
					}
				}
				if isPropertyKey(attr[i]) && attr[i].Val == "og:title" && isContentKey(attr[i-1]) {
					if webPreview.OpenGraph.Title == "" {
						webPreview.OpenGraph.Title = attr[i-1].Val
					}
				}

				if isPropertyKey(attr[i-1]) && attr[i-1].Val == "og:description" && isContentKey(attr[i]) {
					if webPreview.OpenGraph.Description == "" {
						webPreview.OpenGraph.Description = attr[i].Val
					}
				}
				if isPropertyKey(attr[i]) && attr[i].Val == "og:description" && isContentKey(attr[i-1]) {
					if webPreview.OpenGraph.Description == "" {
						webPreview.OpenGraph.Description = attr[i-1].Val
					}
				}

				if isPropertyKey(attr[i-1]) && attr[i-1].Val == "og:image" && isContentKey(attr[i]) {
					if webPreview.OpenGraph.ImgURL == "" {
						webPreview.OpenGraph.ImgURL = supplementImgURL(url, attr[i].Val)
					}
				}
				if isPropertyKey(attr[i]) && attr[i].Val == "og:image" && isContentKey(attr[i-1]) {
					if webPreview.OpenGraph.ImgURL == "" {
						webPreview.OpenGraph.ImgURL = supplementImgURL(url, attr[i-1].Val)
					}
				}

				if isPropertyKey(attr[i-1]) && attr[i-1].Val == "twitter:title" && isContentKey(attr[i]) {
					if webPreview.Twitter.Title == "" {
						webPreview.Twitter.Title = attr[i].Val
					}
				}
				if isPropertyKey(attr[i]) && attr[i].Val == "twitter:title" && isContentKey(attr[i-1]) {
					if webPreview.Twitter.Title == "" {
						webPreview.Twitter.Title = attr[i-1].Val
					}
				}

				if isPropertyKey(attr[i-1]) && attr[i-1].Val == "twitter:description" && isContentKey(attr[i]) {
					if webPreview.Twitter.Description == "" {
						webPreview.Twitter.Description = attr[i].Val
					}
				}
				if isPropertyKey(attr[i]) && attr[i].Val == "twitter:description" && isContentKey(attr[i-1]) {
					if webPreview.Twitter.Description == "" {
						webPreview.Twitter.Description = attr[i-1].Val
					}
				}

				if isPropertyKey(attr[i-1]) && attr[i-1].Val == "twitter:image" && isContentKey(attr[i]) {
					fmt.Println("got img")
					if webPreview.Twitter.ImgURL == "" {
						webPreview.Twitter.ImgURL = supplementImgURL(url, attr[i].Val)
					}
				}
				if isPropertyKey(attr[i]) && attr[i].Val == "twitter:image" && isContentKey(attr[i-1]) {
					if webPreview.Twitter.ImgURL == "" {
						webPreview.Twitter.ImgURL = supplementImgURL(url, attr[i-1].Val)
					}
				}
			}
		}
	}
}

func isPropertyKey(attr html.Attribute) bool {
	return attr.Key == "property"
}

func isNameKey(attr html.Attribute) bool {
	return attr.Key == "name"
}

func isContentKey(attr html.Attribute) bool {
	return attr.Key == "content"
}

func supplementImgURL(url string, imgURL string) string {
	if len(imgURL) > 8 {
		if imgURL[0:5] == "https" || imgURL[0:4] == "http" {
			return imgURL
		}
	}
	if !strings.HasPrefix(imgURL, url) {
		if strings.HasSuffix(url, "/") && !strings.HasPrefix(imgURL, "/") {
			return url + imgURL
		}
		if !strings.HasSuffix(url, "/") && strings.HasPrefix(imgURL, "/") {
			return url + imgURL
		}
		if strings.HasSuffix(url, "/") && strings.HasPrefix(imgURL, "/") {
			return url + imgURL[0:]
		}
		if !strings.HasSuffix(url, "/") && !strings.HasPrefix(imgURL, "/") {
			return url + "/" + imgURL
		}
	}
	return imgURL
}
