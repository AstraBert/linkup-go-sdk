// In this example we demonstrate how to call the `/search` endpoint
// of the Linkup API returning an LLM-generated answer grounded with sources.
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
	output, err := client.GetSourcedAnswer(query, linkup.Standard)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Answer:", output.Answer)
	for _, source := range output.Sources {
		fmt.Println("Source Name:", source.Name)
		fmt.Println("Source URL:", source.Url)
		fmt.Println("Source Snippet:", source.Snippet)
	}
}
