package linkup

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

// assert that MockClient satisfies the LinkupHttpClient interface
var _ LinkupHttpClient = (*MockClient)(nil)

type MockStructuredStruct struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type MockClient struct {
	fails bool
}

func (m *MockClient) SearchWithResponse(ctx context.Context, body SearchJSONRequestBody, requestEditors ...RequestEditorFn) (*SearchResponse, error) {
	if m.fails {
		return &SearchResponse{
			Body: []byte("an error occurred"),
			HTTPResponse: &http.Response{
				Status:     "429 Too Many Requests",
				StatusCode: 429,
			},
		}, nil
	}
	switch body.OutputType {
	case Structured:
		if body.IncludeSources != nil && *body.IncludeSources {
			responseBody := StructuredWithSourcesDto{Data: map[string]any{"title": "hello", "summary": "lorem ipsum dolor"}, Sources: nil}
			marshaled, err := json.Marshal(responseBody)
			if err != nil {
				return nil, err
			}
			return &SearchResponse{
				Body: marshaled,
				HTTPResponse: &http.Response{
					Status:     "200 OK",
					StatusCode: 200,
				},
			}, nil
		} else {
			responseBody := MockStructuredStruct{Title: "hello", Summary: "lorem ipsum dolor"}
			marshaled, err := json.Marshal(responseBody)
			if err != nil {
				return nil, err
			}
			return &SearchResponse{
				Body: marshaled,
				HTTPResponse: &http.Response{
					Status:     "200 OK",
					StatusCode: 200,
				},
			}, nil
		}
	case SearchResults:
		results := []SearchResultsDto_Results_Item{}
		if body.IncludeImages != nil && *body.IncludeImages {
			imageResult := &SearchResultsDto_Results_Item{}
			_ = imageResult.FromImageSearchResultDto(ImageSearchResultDto{
				Name: "lake",
				Type: "image",
				Url:  "https://image.lake.com",
			})
			results = append(results, *imageResult)
		}
		textResult := &SearchResultsDto_Results_Item{}
		_ = textResult.FromTextSearchResultDto(TextSearchResultDto{
			Content: "This is a lake",
			Favicon: "",
			Name:    "lake",
			Type:    "text",
			Url:     "https://thisisalake.com",
		})
		results = append(results, *textResult)
		responseBody := SearchResultsDto{Results: results}
		marshaled, err := json.Marshal(responseBody)
		if err != nil {
			return nil, err
		}
		return &SearchResponse{
			Body: marshaled,
			HTTPResponse: &http.Response{
				Status:     "200 OK",
				StatusCode: 200,
			},
		}, nil
	default:
		responseBody := SourcedAnswerDto{
			Answer: "This is a lake",
			Sources: []SourceDto{
				{
					Favicon: "",
					Name:    "lake",
					Url:     "https://thisisalake.com",
					Snippet: "A lake in the mountains",
				},
			},
		}
		marshaled, err := json.Marshal(responseBody)
		if err != nil {
			return nil, err
		}
		return &SearchResponse{
			Body: marshaled,
			HTTPResponse: &http.Response{
				Status:     "200 OK",
				StatusCode: 200,
			},
		}, nil
	}
}

func (m *MockClient) BalanceWithResponse(ctx context.Context, requestEditors ...RequestEditorFn) (*BalanceResponse, error) {
	if m.fails {
		return &BalanceResponse{
			Body: []byte("an error occurred: too many requests"),
			HTTPResponse: &http.Response{
				Status:     "429 Too Many Requests",
				StatusCode: 429,
			},
			JSON200: nil,
		}, nil
	}
	return &BalanceResponse{
		Body: []byte("3.14"),
		HTTPResponse: &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
		},
		JSON200: &CreditDto{
			Balance: 3.14,
		},
	}, nil
}

func (m *MockClient) FetchWithResponse(ctx context.Context, body FetchJSONRequestBody, requestEditors ...RequestEditorFn) (*FetchResponse, error) {
	if m.fails {
		return &FetchResponse{
			Body: []byte("an error occurred: too many requests"),
			HTTPResponse: &http.Response{
				Status:     "429 Too Many Requests",
				StatusCode: 429,
			},
			JSON200: nil,
		}, nil
	}
	response := &FetchResponseDto{
		Markdown: "# Hello World!",
	}
	if body.ExtractImages != nil && *body.ExtractImages {
		response.Images = &[]FetchImageDto{
			{Alt: nil, Url: "https://helloworld.image.png"},
		}
	}
	if body.IncludeRawHtml != nil && *body.IncludeRawHtml {
		rawHtml := "<h1>Hello World!</h1>"
		response.RawHtml = &rawHtml
	}
	return &FetchResponse{
		Body: []byte("response"),
		HTTPResponse: &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
		},
		JSON200: response,
	}, nil
}

func TestGetSearchResultsTextOnlySuccess(t *testing.T) {
	client := LinkupClient{
		apiKey: "hello",
		client: &MockClient{fails: false},
	}
	output, err := client.GetSearchResults("lake", Standard)
	if err != nil {
		t.Fatalf("An unexpected error occurred: %s", err.Error())
	}
	if len(output.TextResults) != 1 {
		t.Fatalf("Expected 1 text result, got %d", len(output.TextResults))
	}
	if len(output.ImageResults) != 0 {
		t.Fatalf("Expected 0 image results, got %d", len(output.ImageResults))
	}
	result := output.TextResults[0]
	if result.Content != "This is a lake" || result.Favicon != "" || result.Name != "lake" || result.Type != "text" || result.Url != "https://thisisalake.com" {
		t.Fatalf("Unexpected text result: %v", result)
	}
}

func TestGetSearchResultsWithImagesSuccess(t *testing.T) {
	client := LinkupClient{
		apiKey: "hello",
		client: &MockClient{fails: false},
	}
	output, err := client.GetSearchResults("lake", Standard, AdditionalSearchOptions{IncludeImages: true})
	if err != nil {
		t.Fatalf("An unexpected error occurred: %s", err.Error())
	}
	if len(output.TextResults) != 1 {
		t.Fatalf("Expected 1 text result, got %d", len(output.TextResults))
	}
	if len(output.ImageResults) != 1 {
		t.Fatalf("Expected 1 image result, got %d", len(output.ImageResults))
	}
	textResult := output.TextResults[0]
	if textResult.Content != "This is a lake" || textResult.Favicon != "" || textResult.Name != "lake" || textResult.Type != "text" || textResult.Url != "https://thisisalake.com" {
		t.Fatalf("Unexpected text result: %v", textResult)
	}
	imageResult := output.ImageResults[0]
	if imageResult.Name != "lake" || imageResult.Type != "image" || imageResult.Url != "https://image.lake.com" {
		t.Fatalf("Unexpected image result: %v", imageResult)
	}
}

func TestGetSearchResultsFails(t *testing.T) {
	client := LinkupClient{
		apiKey: "hello",
		client: &MockClient{fails: true},
	}
	_, err := client.GetSearchResults("lake", Standard, AdditionalSearchOptions{IncludeImages: true})
	if err != nil {
		if err.Error() != "response returned a status code of 429: 429 Too Many Requests" {
			t.Fatalf("Unexpected error: %s", err.Error())
		}
	} else {
		t.Fatalf("No error recorded, but one was expected")
	}
}

func TestGetStructuredResultsSuccess(t *testing.T) {
	client := LinkupClient{
		apiKey: "hello",
		client: &MockClient{fails: false},
	}
	schema, err := GenerateJSONSchema[MockStructuredStruct]()
	if err != nil {
		t.Fatalf("An error occurred while generating the JSON schema: %s", err.Error())
	}
	output, err := client.GetStructuredResults("summary", Standard, schema)
	if err != nil {
		t.Fatalf("An unexpected error occurred: %s", err.Error())
	}
	if output.RawJson == nil {
		t.Fatal("output.RawJson should be non-null")
	} else {
		result, err := GetResultFromRawJSON[MockStructuredStruct](*output.RawJson)
		if err != nil {
			t.Fatalf("An unexpected error occurred: %s", err.Error())
		}
		typedResult, ok := result.(MockStructuredStruct)
		if !ok {
			t.Fatal("result should be of type MockStructuredStruct")
		}
		if typedResult.Summary != "lorem ipsum dolor" || typedResult.Title != "hello" {
			t.Fatalf("Unexpected result: %v", typedResult)
		}
	}
	if output.SourcedOutput != nil {
		t.Fatalf("output.SourcedOutput should be non-null, got %v", *output.SourcedOutput)
	}
}

func TestGetStructuredResultsWithSourcesSuccess(t *testing.T) {
	client := LinkupClient{
		apiKey: "hello",
		client: &MockClient{fails: false},
	}
	schema, err := GenerateJSONSchema[MockStructuredStruct]()
	if err != nil {
		t.Fatalf("An error occurred while generating the JSON schema: %s", err.Error())
	}
	output, err := client.GetStructuredResults("summary", Standard, schema, AdditionalSearchOptions{IncludeSources: true})
	if err != nil {
		t.Fatalf("An unexpected error occurred: %s", err.Error())
	}
	if output.RawJson != nil {
		t.Fatalf("output.RawJson should be null, got %s", *output.RawJson)
	}
	if output.SourcedOutput == nil {
		t.Fatal("output.SourcedOutput should be non-null")
	} else {
		data := output.SourcedOutput.Data
		summary, ok := data["summary"]
		if !ok {
			t.Fatal("summary not in data")
		}
		title, ok := data["title"]
		if !ok {
			t.Fatal("title not in data")
		}
		if summary != "lorem ipsum dolor" || title != "hello" {
			t.Fatalf("Unexpected data: {'title': '%v', 'summary': '%v'}", title, summary)
		}
		if output.SourcedOutput.Sources != nil {
			t.Fatal("sources should be null")
		}
	}
}

func TestGetStructuredResultsFails(t *testing.T) {
	client := LinkupClient{
		apiKey: "hello",
		client: &MockClient{fails: true},
	}
	schema, err := GenerateJSONSchema[MockStructuredStruct]()
	if err != nil {
		t.Fatalf("An error occurred while generating the JSON schema: %s", err.Error())
	}
	_, err = client.GetStructuredResults("lake", Standard, schema)
	if err != nil {
		if err.Error() != "response returned a status code of 429: 429 Too Many Requests" {
			t.Fatalf("Unexpected error: %s", err.Error())
		}
	} else {
		t.Fatalf("No error recorded, but one was expected")
	}
}

func TestGetSourcedAnswerSuccess(t *testing.T) {
	client := LinkupClient{
		apiKey: "hello",
		client: &MockClient{fails: false},
	}
	output, err := client.GetSourcedAnswer("lake", Standard)
	if err != nil {
		t.Fatalf("An unexpected error occurred: %s", err.Error())
	}
	if output.Answer != "This is a lake" {
		t.Fatalf("Unexpected answer: %s", output.Answer)
	}
	if len(output.Sources) != 1 {
		t.Fatalf("Expected 1 source, got %d", len(output.Sources))
	}
	source := output.Sources[0]
	if source.Favicon != "" || source.Name != "lake" || source.Url != "https://thisisalake.com" || source.Snippet != "A lake in the mountains" {
		t.Fatalf("Unexpected source value: %v", source)
	}
}

func TestGetSourcedAnswerFails(t *testing.T) {
	client := LinkupClient{
		apiKey: "hello",
		client: &MockClient{fails: true},
	}
	_, err := client.GetSourcedAnswer("lake", Standard)
	if err != nil {
		if err.Error() != "response returned a status code of 429: 429 Too Many Requests" {
			t.Fatalf("Unexpected error: %s", err.Error())
		}
	} else {
		t.Fatalf("No error recorded, but one was expected")
	}
}

func TestGetBalanceSuccess(t *testing.T) {
	client := LinkupClient{
		apiKey: "hello",
		client: &MockClient{fails: false},
	}
	balance, err := client.GetBalance()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if balance != 3.14 {
		t.Fatalf("Expecting a balance of %f, got %f", 3.14, balance)
	}
}

func TestGetBalanceFails(t *testing.T) {
	client := LinkupClient{
		apiKey: "hello",
		client: &MockClient{fails: true},
	}
	_, err := client.GetBalance()
	if err != nil {
		if err.Error() != "response returned a status code of 429: 429 Too Many Requests" {
			t.Fatalf("Unexpected error: %s", err.Error())
		}
	} else {
		t.Fatalf("No error recorded, but one was expected")
	}
}

func TestFetchSuccess(t *testing.T) {
	client := LinkupClient{
		apiKey: "hello",
		client: &MockClient{fails: false},
	}
	output, err := client.Fetch("https://fetch.com", AdditionalFetchOptions{IncludeRawHtml: true, ExtractImages: true})
	if err != nil {
		t.Fatalf("An unexpected error occurred: %s", err.Error())
	}
	if output.Markdown != "# Hello World!" {
		t.Fatalf("Unexpected markdown: %s", output.Markdown)
	}
	if output.RawHtml == nil || (output.RawHtml != nil && *output.RawHtml != "<h1>Hello World!</h1>") {
		if output.RawHtml == nil {
			t.Fatal("Unexpected html: null")
		} else {
			t.Fatalf("Unexpected html: %s", *output.RawHtml)
		}
	}
	if output.Images == nil {
		t.Fatal("images should be non-null")
	}
	if len(*output.Images) != 1 {
		t.Fatalf("Expected 1 image, got %d", len(*output.Images))
	}
	image := (*output.Images)[0]
	if image.Alt != nil || image.Url != "https://helloworld.image.png" {
		t.Fatalf("Unexpected image: %v", image)
	}
}

func TestFetchMarkdownOnlySuccess(t *testing.T) {
	client := LinkupClient{
		apiKey: "hello",
		client: &MockClient{fails: false},
	}
	output, err := client.Fetch("https://fetch.com")
	if err != nil {
		t.Fatalf("An unexpected error occurred: %s", err.Error())
	}
	if output.Markdown != "# Hello World!" {
		t.Fatalf("Unexpected markdown: %s", output.Markdown)
	}
	if output.RawHtml != nil {
		t.Fatalf("Unexpected html: %s", *output.RawHtml)
	}
	if output.Images != nil {
		t.Fatal("images should be null")
	}
}

func TestFetchMarkdownImagesSuccess(t *testing.T) {
	client := LinkupClient{
		apiKey: "hello",
		client: &MockClient{fails: false},
	}
	output, err := client.Fetch("https://fetch.com", AdditionalFetchOptions{ExtractImages: true})
	if err != nil {
		t.Fatalf("An unexpected error occurred: %s", err.Error())
	}
	if output.Markdown != "# Hello World!" {
		t.Fatalf("Unexpected markdown: %s", output.Markdown)
	}
	if output.RawHtml != nil {
		t.Fatalf("Unexpected html: %s", *output.RawHtml)
	}
	if output.Images == nil {
		t.Fatal("images should be non-null")
	}
	if len(*output.Images) != 1 {
		t.Fatalf("Expected 1 image, got %d", len(*output.Images))
	}
	image := (*output.Images)[0]
	if image.Alt != nil || image.Url != "https://helloworld.image.png" {
		t.Fatalf("Unexpected image: %v", image)
	}
}

func TestFetchMarkdownHtmlSuccess(t *testing.T) {
	client := LinkupClient{
		apiKey: "hello",
		client: &MockClient{fails: false},
	}
	output, err := client.Fetch("https://fetch.com", AdditionalFetchOptions{IncludeRawHtml: true})
	if err != nil {
		t.Fatalf("An unexpected error occurred: %s", err.Error())
	}
	if output.Markdown != "# Hello World!" {
		t.Fatalf("Unexpected markdown: %s", output.Markdown)
	}
	if output.RawHtml == nil || (output.RawHtml != nil && *output.RawHtml != "<h1>Hello World!</h1>") {
		if output.RawHtml == nil {
			t.Fatal("Unexpected html: null")
		} else {
			t.Fatalf("Unexpected html: %s", *output.RawHtml)
		}
	}
	if output.Images != nil {
		t.Fatal("images should be null")
	}
}

func TestFetchFails(t *testing.T) {
	client := LinkupClient{
		apiKey: "hello",
		client: &MockClient{fails: true},
	}
	_, err := client.Fetch("https://fetch.com", AdditionalFetchOptions{IncludeRawHtml: true, ExtractImages: true})
	if err != nil {
		if err.Error() != "response returned a status code of 429: 429 Too Many Requests" {
			t.Fatalf("Unexpected error: %s", err.Error())
		}
	} else {
		t.Fatalf("No error recorded, but one was expected")
	}
}
