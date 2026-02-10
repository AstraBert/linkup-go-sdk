// In this example we demonstrate how to call the `/search` endpoint
// of the Linkup API returning an output which follows a specified JSON schema,
// optionally retrieving sources along with it.
package main

import (
	"fmt"
	"log"

	"github.com/AstraBert/linkup-go-sdk"
)

// Define a struct which will be converted to JSON schema
type LakeComoResult struct {
	Castles []string `json:"castles" jsonschema:"title=castles,description=Castles to visit on Lake Como"`
	Cities  []string `json:"cities" jsonschema:"title=cities,description=Cities to visit on Lake Como"`
}

func main() {
	query := "Places to visit on Lake Como"
	// we provide the key as a empty string to load it from the environment
	client, err := linkup.NewLinkupClient("")
	if err != nil {
		log.Fatal(err)
	}
	// we create a JSON schema with the GenerateJSONSchema utility function
	schema, err := linkup.GenerateJSONSchema[LakeComoResult]()
	if err != nil {
		log.Fatal(err)
	}
	// we get the structured result without sources
	output, err := client.GetStructuredResults(query, linkup.Standard, schema)
	if err != nil {
		log.Fatal(err)
	}
	if output.RawJson == nil {
		log.Fatalln("RawJson should be non-null")
	}
	// we convert the raw JSON string to a struct type
	result, err := linkup.GetResultFromRawJSON[LakeComoResult](*output.RawJson)
	if err != nil {
		log.Fatal(err)
	}
	// we assert that the struct type conforms to LakeComoResult
	structuredResult, ok := result.(LakeComoResult)
	if !ok {
		log.Fatalf("Unexpected result: %v\n", result)
	}
	fmt.Println("Castles to visit:")
	for _, castle := range structuredResult.Castles {
		fmt.Println(castle)
	}
	fmt.Println("Cities to visit:")
	for _, city := range structuredResult.Cities {
		fmt.Println(city)
	}
	// we can also get the structured output with sources
	outputWithSources, err := client.GetStructuredResults(query, linkup.Standard, schema, linkup.AdditionalSearchOptions{IncludeSources: true})
	if err != nil {
		log.Fatal(err)
	}
	if outputWithSources.SourcedOutput == nil {
		log.Fatalln("SourcedOutput should be non-null")
	}
	// we convert the SourcedOutput to a struct type
	resultWithSource, err := linkup.GetResultFromSourcedOutput[LakeComoResult](output.SourcedOutput)
	if err != nil {
		log.Fatal(err)
	}
	// we assert that the struct type conforms to LakeComoResult
	structuredResultWithSource, ok := resultWithSource.(LakeComoResult)
	if !ok {
		log.Fatalf("Unexpected result: %v\n", resultWithSource)
	}
	fmt.Println("Castles to visit:")
	for _, castle := range structuredResultWithSource.Castles {
		fmt.Println(castle)
	}
	fmt.Println("Cities to visit:")
	for _, city := range structuredResultWithSource.Cities {
		fmt.Println(city)
	}
	// print all the sources
	if outputWithSources.SourcedOutput.Sources != nil {
		for _, source := range outputWithSources.SourcedOutput.Sources {
			if source.Name != nil {
				fmt.Println("Source Name:", *source.Name)
			}
			if source.Url != nil {
				fmt.Println("Source URL:", *source.Url)
			}
			if source.Content != nil {
				fmt.Println("Source Snippet:", *source.Content)
			}
		}
	}
}
