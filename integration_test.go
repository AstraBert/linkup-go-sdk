package linkup

import (
	"log"
	"os"
	"testing"
)

func hasLinkupApiKey() bool {
	_, ok := os.LookupEnv("LINKUP_API_KEY")
	return ok
}

func TestGetSearchResultsIntegration(t *testing.T) {
	if hasLinkupApiKey() {
		client, err := NewLinkupClient("")
		if err != nil {
			t.Fatalf("An error occurred: %s", err.Error())
		}
		output, err := client.GetSearchResults("Lake Como", Standard, AdditionalSearchOptions{IncludeImages: true})
		if err != nil {
			t.Fatalf("Unexpected error performing the search: %s", err.Error())
		}
		if (len(output.ImageResults) + len(output.TextResults)) == 0 {
			t.Fatalf("No results produced")
		}
	} else {
		log.Println("Skipping TestGetSearchResultsIntegration...")
		t.Skip("LINKUP_API_KEY not available")
	}
}

type LakeComoResult struct {
	PlacesToVisit *[]string `json:"places_to_visit,omitempty" jsonschema:"title=places to visit,description=places to visit on Lake Como"`
	Book          string    `json:"book" jsonschema:"title=book,description=most famous book written about Lake Como"`
}

func TestGetStructuredResultsIntegration(t *testing.T) {
	if hasLinkupApiKey() {
		client, err := NewLinkupClient("")
		if err != nil {
			t.Fatalf("An error occurred: %s", err.Error())
		}
		schema, err := GenerateJSONSchema[LakeComoResult]()
		if err != nil {
			t.Fatalf("An error occurred: %s", err.Error())
		}
		output, err := client.GetStructuredResults("Lake Como", Standard, schema)
		if err != nil {
			t.Fatalf("Unexpected error performing the search: %s", err.Error())
		}
		if output.RawJson == nil {
			t.Fatal("Expected RawJson to be non-null")
		} else {
			val, err := GetResultFromRawJSON[LakeComoResult](*output.RawJson)
			if err != nil {
				t.Fatalf("Unexpected error converting to original struct: %s", err.Error())
			}
			_, ok := val.(LakeComoResult)
			if !ok {
				t.Fatalf("Expected value to be of type LakeComoResult, got %v", val)
			}
		}
		if output.SourcedOutput != nil {
			t.Fatalf("SourcedOutput should be null")
		}
	} else {
		log.Println("Skipping TestGetStructuredResultsIntegration...")
		t.Skip("LINKUP_API_KEY not available")
	}
}

func TestGetStructuredResultsWithSourcesIntegration(t *testing.T) {
	if hasLinkupApiKey() {
		client, err := NewLinkupClient("")
		if err != nil {
			t.Fatalf("An error occurred: %s", err.Error())
		}
		schema, err := GenerateJSONSchema[LakeComoResult]()
		if err != nil {
			t.Fatalf("An error occurred: %s", err.Error())
		}
		output, err := client.GetStructuredResults("Lake Como", Standard, schema, AdditionalSearchOptions{IncludeSources: true})
		if err != nil {
			t.Fatalf("Unexpected error performing the search: %s", err.Error())
		}
		if output.RawJson != nil {
			t.Fatal("Expected RawJson to be null")
		}
		if output.SourcedOutput == nil {
			t.Fatalf("SourcedOutput should be not-null")
		} else {
			val, err := GetResultFromSourcedOutput[LakeComoResult](output.SourcedOutput)
			if err != nil {
				t.Fatalf("Unexpected error converting to original struct: %s", err.Error())
			}
			_, ok := val.(LakeComoResult)
			if !ok {
				t.Fatalf("Expected value to be of type LakeComoResult, got %v", val)
			}
			if len(output.SourcedOutput.Sources) == 0 {
				t.Fatalf("Expected at least one source")
			}
		}
	} else {
		log.Println("Skipping TestGetStructuredResultsWithSourcesIntegration...")
		t.Skip("LINKUP_API_KEY not available")
	}
}

func TestGetSourcedAnswerIntegration(t *testing.T) {
	if hasLinkupApiKey() {
		client, err := NewLinkupClient("")
		if err != nil {
			t.Fatalf("An error occurred: %s", err.Error())
		}
		output, err := client.GetSourcedAnswer("Lake Como", Standard)
		if err != nil {
			t.Fatalf("Unexpected error performing the search: %s", err.Error())
		}
		if output.Answer == "" {
			t.Fatal("Expected answer to be non-empty")
		}
		if len(output.Sources) == 0 {
			t.Fatal("Expected at least one source")
		}
	} else {
		log.Println("Skipping TestGetSourcedAnswerIntegration...")
		t.Skip("LINKUP_API_KEY not available")
	}
}

func TestGetBalanceIntegration(t *testing.T) {
	if hasLinkupApiKey() {
		client, err := NewLinkupClient("")
		if err != nil {
			t.Fatalf("An error occurred: %s", err.Error())
		}
		output, err := client.GetBalance()
		if err != nil {
			t.Fatalf("Unexpected error performing the GetBalance operation: %s", err.Error())
		}
		if output < 0 {
			t.Fatalf("Output cannot be less than 0")
		}
	} else {
		log.Println("Skipping TestGetBalanceIntegration...")
		t.Skip("LINKUP_API_KEY not available")
	}
}

func TestFetchIntegration(t *testing.T) {
	if hasLinkupApiKey() {
		client, err := NewLinkupClient("")
		if err != nil {
			t.Fatalf("An error occurred: %s", err.Error())
		}
		output, err := client.Fetch("https://clelia.dev/2026-01-31-why-dont-i-vibe-code-more", AdditionalFetchOptions{IncludeRawHtml: true})
		if err != nil {
			t.Fatalf("Unexpected error performing the GetBalance operation: %s", err.Error())
		}
		if output.Markdown == "" {
			t.Fatal("Expected markdown to be non-empty")
		}
		if output.RawHtml == nil || (output.RawHtml != nil && *output.RawHtml == "") {
			t.Fatal("Expected raw HTML to be non-null and non-empty")
		}
	} else {
		log.Println("Skipping TestFetchIntegration...")
		t.Skip("LINKUP_API_KEY not available")
	}
}
