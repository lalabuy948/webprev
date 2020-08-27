// Package webprev provides easy extraction of website previews. Generic, Facebook and Twitter cards.
package webprev

import (
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

		if isGivenMeta(attr, "title") {
			webPreview.Generic.Title = extractMetaContentValue(attr)
		}
		if isGivenMeta(attr, "description") {
			webPreview.Generic.Description = extractMetaContentValue(attr)
		}
		if isMetaItemprop(attr) {
			webPreview.Generic.ImgURL = supplementImgURL(url, extractMetaContentValue(attr))
		}
		if isGivenMeta(attr, "og:title") {
			webPreview.OpenGraph.Title = extractMetaContentValue(attr)
		}
		if isGivenMeta(attr, "og:description") {
			webPreview.OpenGraph.Description = extractMetaContentValue(attr)
		}
		if isGivenMeta(attr, "og:image") {
			webPreview.OpenGraph.ImgURL = supplementImgURL(url, extractMetaContentValue(attr))
		}
		if isGivenMeta(attr, "twitter:title") {
			webPreview.Twitter.Title = extractMetaContentValue(attr)
		}
		if isGivenMeta(attr, "twitter:description") {
			webPreview.Twitter.Description = extractMetaContentValue(attr)
		}
		if isGivenMeta(attr, "twitter:image") {
			webPreview.Twitter.ImgURL = supplementImgURL(url, extractMetaContentValue(attr))
		}
	}
}

func isGivenMeta(attrs []html.Attribute, lookup string) bool {
	for i := 0; i < len(attrs); i++ {
		if isNameKey(attrs[i]) || isPropertyKey(attrs[i]) {
			if attrs[i].Val == lookup {
				return true
			}
		}
	}
	return false
}

func isMetaItemprop(attrs []html.Attribute) bool {
	for i := 0; i < len(attrs); i++ {
		if isItempropKey(attrs[i]) && attrs[i].Val == "image" {
			return true
		}
	}
	return false
}

func extractMetaContentValue(attrs []html.Attribute) string {
	for i := 0; i < len(attrs); i++ {
		if isContentKey(attrs[i]) {
			return attrs[i].Val
		}
	}
	return ""
}

func isItempropKey(attr html.Attribute) bool {
	return attr.Key == "itemprop"
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
	if imgURL == "" {
		return imgURL
	}
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
