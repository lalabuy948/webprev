package main

import (
	"fmt"
	"log"
	"webprev"
)

func main() {
	webPrev, err := webprev.Preview("https://telegram.org")
	if err != nil {
		log.Fatal(err)
	}

	// default / generic ones
	fmt.Println("title: ", webPrev.Generic.Title)
	fmt.Println("description: ", webPrev.Generic.Description)
	fmt.Println("img_url: ", webPrev.Generic.ImgURL)

	// OpenGraph aka facebook
	fmt.Println("op:title: ", webPrev.OpenGraph.Title)
	fmt.Println("op:description: ", webPrev.OpenGraph.Description)
	fmt.Println("op:img_url: ", webPrev.OpenGraph.ImgURL)

	// twitter card
	fmt.Println("twitter:title: ", webPrev.Twitter.Title)
	fmt.Println("twitter:description: ", webPrev.Twitter.Description)
	fmt.Println("twitter:img_url: ", webPrev.Twitter.ImgURL)
}
