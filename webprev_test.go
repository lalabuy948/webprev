package webprev

import (
	"testing"
)

func TestPreview(t *testing.T) {
	webPrev, err := Preview("https://basicbb.com")
	if err != nil {
		t.Errorf("Preview() err = %v; want nothing", err)
	}
	if webPrev.Generic.Title == "" {
		t.Errorf("Preview() title = %v; want something", webPrev.Generic.Title)
	}
	//fmt.Println(webPrev.Generic.Title)
	//fmt.Println(webPrev.Generic.Description)
	//fmt.Println(webPrev.Generic.ImgURL)
}

func TestSupplementImgURL(t *testing.T) {
	if got := supplementImgURL("", "bla"); got != "bla" {
		t.Errorf("supplementImgURL('', 'bla') = %v; want true", got)
	}

	if got := supplementImgURL("https://some.domain", "http://other.domain"); got != "http://other.domain" {
		t.Errorf("supplementImgURL('', 'bla') = %v; want http://other.domain", got)
	}

	if got := supplementImgURL("https://some.domain/", "img.jpeg"); got != "https://some.domain/img.jpeg" {
		t.Errorf("supplementImgURL('', 'bla') = %v; want https://some.domain/img.jpeg", got)
	}

	if got := supplementImgURL("https://some.domain", "/img.jpeg"); got != "https://some.domain/img.jpeg" {
		t.Errorf("supplementImgURL('', 'bla') = %v; want https://some.domain/img.jpeg", got)
	}
}
