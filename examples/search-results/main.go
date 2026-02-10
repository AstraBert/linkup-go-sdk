// In this example we demonstrate how to call the `/search` endpoint
// of the Linkup API returning a list of text or image search results.
package main

import (
	"fmt"
	"log"

	"github.com/AstraBert/linkup-go-sdk"
)

func main() {
	query := "Places to visit on Lake Como"
	// we provide the key as a empty string to load it from the environment
	client, err := linkup.NewLinkupClient("")
	if err != nil {
		log.Fatal(err)
	}
	// only retrieve text-based results
	output, err := client.GetSearchResults(query, linkup.Standard)
	if err != nil {
		log.Fatal(err)
	}
	for _, result := range output.TextResults {
		fmt.Println("Name:", result.Name)
		fmt.Println("Content:", result.Content)
		fmt.Println("Url:", result.Url)
		fmt.Println("Favicon:", result.Favicon)
	}
	// retrieve both text-based and image-based results
	completeOutput, err := client.GetSearchResults(query, linkup.Standard, linkup.AdditionalSearchOptions{IncludeImages: true})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Texts:")
	for _, result := range completeOutput.TextResults {
		fmt.Println("Name:", result.Name)
		fmt.Println("Content:", result.Content)
		fmt.Println("Url:", result.Url)
		fmt.Println("Favicon:", result.Favicon)
	}
	fmt.Println("Images:")
	for _, image := range completeOutput.ImageResults {
		fmt.Println("Name:", image.Name)
		fmt.Println("Url:", image.Url)
	}
}
