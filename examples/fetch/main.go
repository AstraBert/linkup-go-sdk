// In this example we demonstrate how to call the `/fetch` endpoint of the Linkup API
// in order to scrape the content of a URL as markdown, with the possibility of adding
// images and raw HTML
package main

import (
	"fmt"
	"log"

	"github.com/AstraBert/linkup-go-sdk"
)

func main() {
	url := "https://clelia.dev/2026-02-09-the-anatomy-of-a-document-processing-agent"
	// we provide the key as a empty string to load it from the environment
	client, err := linkup.NewLinkupClient("")
	if err != nil {
		log.Fatal(err)
	}
	// fetch only as markdown
	result, err := client.Fetch(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Markdown)
	// fetch with markdown, raw HTML and images
	completeResult, err := client.Fetch(url, linkup.AdditionalFetchOptions{IncludeRawHtml: true, ExtractImages: true})
	if err != nil {
		log.Fatal(err)
	}
	if completeResult.RawHtml == nil {
		log.Fatalln("RawHtml should be non-null")
	}
	fmt.Println(*completeResult.RawHtml)
	if completeResult.Images == nil {
		log.Fatalln("Images should be non-null")
	}
	for _, image := range *completeResult.Images {
		fmt.Println("Image URL:", image.Url)
		if image.Alt != nil {
			fmt.Println("Alt text:", *image.Alt)
		}
	}
}
