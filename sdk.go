package linkup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

const LinkupServerUrl string = "https://api.linkup.so"

// Helper interface to reduce the scope of the underlying HTTP client for Linkup (mostly for testing purposes)
type LinkupHttpClient interface {
	SearchWithResponse(context.Context, SearchJSONRequestBody, ...RequestEditorFn) (*SearchResponse, error)
	BalanceWithResponse(context.Context, ...RequestEditorFn) (*BalanceResponse, error)
	FetchWithResponse(context.Context, FetchJSONRequestBody, ...RequestEditorFn) (*FetchResponse, error)
}

// Struct type representing a client to perform operations with the Linkup API
type LinkupClient struct {
	apiKey string
	client LinkupHttpClient
}

// Constructor to create a new LinkupClient instance.
// If the API Key is passed as an empty string, it will be loaded
// from the environment
func NewLinkupClient(apiKey string) (*LinkupClient, error) {
	if apiKey == "" {
		key, ok := os.LookupEnv("LINKUP_API_KEY")
		if !ok {
			return nil, errors.New("api key not provided and could not find LINKUP_API_KEY in the environment")
		}
		apiKey = key
	}
	var requestEditor RequestEditorFn = func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+apiKey)
		return nil
	}
	client, err := NewClientWithResponses(LinkupServerUrl, WithRequestEditorFn(requestEditor))
	if err != nil {
		return nil, err
	}
	return &LinkupClient{
		apiKey: apiKey,
		client: client,
	}, nil
}

// Additional search options to be used with search methods for customization
type AdditionalSearchOptions struct {
	// ExcludeDomains The domains you want to exclude of the search. By default, don't restrict the search.
	ExcludeDomains []string `json:"excludeDomains,omitempty"`

	// FromDate The date from which the search results should be considered, in ISO 8601 format (YYYY-MM-DD). It must be before `toDate`, if provided, and later than 1970-01-01.
	FromDate *string `json:"fromDate,omitempty"`

	// IncludeDomains The domains you want to search on. By default, don't restrict the search. You can provide up to 100 domains.
	IncludeDomains []string `json:"includeDomains,omitempty"`

	// IncludeImages Defines whether the API should include images in its results.
	IncludeImages bool `json:"includeImages,omitempty"`

	// IncludeInlineCitations Relevant only when `outputType` is `sourcedAnswer`. Defines whether the answer should include inline citations.
	IncludeInlineCitations bool `json:"includeInlineCitations,omitempty"`

	// IncludeSources Relevant only when `outputType` is `structured`. Defines whether the response should include sources. **Please note that it modifies the schema of the response, see below**
	IncludeSources bool `json:"includeSources,omitempty"`

	// MaxResults The maximum number of results to return.
	MaxResults *float32 `json:"maxResults,omitempty"`

	// ToDate The date until which the search results should be considered, in ISO 8601 format (YYYY-MM-DD). It must be later than `fromDate`, if provided, or than 1970-01-01.
	ToDate *string `json:"toDate,omitempty"`
}

func DefaultAdditionalSearchOptions() AdditionalSearchOptions {
	return AdditionalSearchOptions{
		ExcludeDomains:         nil,
		IncludeDomains:         nil,
		FromDate:               nil,
		ToDate:                 nil,
		IncludeImages:          false,
		MaxResults:             nil,
		IncludeInlineCitations: false,
		IncludeSources:         false,
	}
}

// Additional option to be used with the fetch method for customization
type AdditionalFetchOptions struct {
	// ExtractImages Defines whether the API should extract the images from the webpage in its response.
	ExtractImages bool `json:"extractImages,omitempty"`

	// IncludeRawHtml Defines whether the API should include the raw HTML of the webpage in its response.
	IncludeRawHtml bool `json:"includeRawHtml,omitempty"`

	// RenderJs Defines whether the API should render the JavaScript of the webpage.
	RenderJs bool `json:"renderJs,omitempty"`
}

func DefaultAdditionalFetchOptions() AdditionalFetchOptions {
	return AdditionalFetchOptions{
		ExtractImages:  false,
		IncludeRawHtml: false,
		RenderJs:       false,
	}
}

// String enum representing the depth that a search should have (alias type)
type SearchDepth = QuerySearchDtoDepth

// Struct type representing results from the `/v1/search` endpoint when `searchResults` is used as output type
type SearchResultsOutput struct {
	ImageResults []ImageSearchResultDto
	TextResults  []TextSearchResultDto
}

// Struct type representing results from the `/v1/search` endpoint when `sourcedAnswer` is used as output type
type SourcedAnswerOutput = SourcedAnswerDto

// Struct type representing results from the `/v1/search` endpoint when `structured` is used as output type
type StructuredOutput struct {
	// Raw JSON output following the provided schema.
	// This field is non-null only if `includeSources` is set to `false`
	// and can be parsed back into the original struct type
	// using the `GetResultFromRawJSON` generic function
	RawJson *string

	// Structured output with sources.
	// This field is non-null only if `includeSources` is set to `true`
	SourcedOutput *StructuredWithSourcesDto
}

// Struct type representing the results from the `/v1/fetch` endpoint
type FetchOutput = FetchResponseDto

// Method to query the /v1/search API endpoint with `searchResults` as output type.
func (l *LinkupClient) GetSearchResults(
	query string,
	depth SearchDepth,
	searchOptions ...AdditionalSearchOptions,
) (*SearchResultsOutput, error) {
	var options AdditionalSearchOptions
	switch len(searchOptions) {
	case 0:
		options = DefaultAdditionalSearchOptions()
	default:
		options = searchOptions[0]
	}
	searchQuery := SearchJSONRequestBody{
		Depth:                  depth,
		Q:                      query,
		ExcludeDomains:         &options.ExcludeDomains,
		IncludeDomains:         &options.IncludeDomains,
		FromDate:               options.FromDate,
		ToDate:                 options.ToDate,
		IncludeImages:          &options.IncludeImages,
		MaxResults:             options.MaxResults,
		IncludeInlineCitations: &options.IncludeInlineCitations,
		IncludeSources:         &options.IncludeSources,
		StructuredOutputSchema: nil,
		OutputType:             SearchResults,
	}
	response, err := l.client.SearchWithResponse(context.Background(), searchQuery)
	if err != nil {
		return nil, err
	}
	if 200 <= response.StatusCode() && response.StatusCode() <= 299 {
		var results SearchResultsDto
		err := json.Unmarshal(response.Body, &results)
		if err != nil {
			return nil, err
		}
		output := SearchResultsOutput{
			TextResults:  make([]TextSearchResultDto, 0, len(results.Results)),
			ImageResults: make([]ImageSearchResultDto, 0, len(results.Results)),
		}
		for i, r := range results.Results {
			result, errTxt := r.AsTextSearchResultDto()
			if errTxt != nil {
				log.Printf("Skipping %d-th result, as it cannot be represented as a text result nor as an image result\n", i)
				continue
			} else {
				if result.Type == "image" {
					output.ImageResults = append(output.ImageResults, ImageSearchResultDto{
						Type: ImageSearchResultDtoType(result.Type),
						Url:  result.Url,
						Name: result.Name,
					})
				} else {
					output.TextResults = append(output.TextResults, result)
				}
			}
		}
		if len(output.TextResults)+len(output.ImageResults) == 0 {
			return nil, errors.New("no valid results were found")
		}
		return &output, nil
	}
	return nil, fmt.Errorf("response returned a status code of %d: %s", response.StatusCode(), response.Status())
}

// Method to query the /v1/search API endpoint with `sourcedAnswer` as output type.
func (l *LinkupClient) GetSourcedAnswer(
	query string,
	depth SearchDepth,
	searchOptions ...AdditionalSearchOptions,
) (*SourcedAnswerOutput, error) {
	var options AdditionalSearchOptions
	switch len(searchOptions) {
	case 0:
		options = DefaultAdditionalSearchOptions()
	default:
		options = searchOptions[0]
	}
	searchQuery := SearchJSONRequestBody{
		Depth:                  depth,
		Q:                      query,
		ExcludeDomains:         &options.ExcludeDomains,
		IncludeDomains:         &options.IncludeDomains,
		FromDate:               options.FromDate,
		ToDate:                 options.ToDate,
		IncludeImages:          &options.IncludeImages,
		MaxResults:             options.MaxResults,
		IncludeInlineCitations: &options.IncludeInlineCitations,
		IncludeSources:         &options.IncludeSources,
		StructuredOutputSchema: nil,
		OutputType:             SourcedAnswer,
	}
	response, err := l.client.SearchWithResponse(context.Background(), searchQuery)
	if err != nil {
		return nil, err
	}
	if 200 <= response.StatusCode() && response.StatusCode() <= 299 {
		var results SourcedAnswerDto
		err := json.Unmarshal(response.Body, &results)
		if err != nil {
			return nil, err
		}
		return &results, nil
	}
	return nil, fmt.Errorf("response returned a status code of %d: %s", response.StatusCode(), response.Status())
}

// Method to query the /v1/search API endpoint with `structured` as output type.
// Use `GenerateJSONSchema` to create the JSON schema needed as an argument to this method.
func (l *LinkupClient) GetStructuredResults(
	query string,
	depth SearchDepth,
	jsonSchema json.RawMessage,
	searchOptions ...AdditionalSearchOptions,
) (*StructuredOutput, error) {
	var options AdditionalSearchOptions
	switch len(searchOptions) {
	case 0:
		options = DefaultAdditionalSearchOptions()
	default:
		options = searchOptions[0]
	}
	searchQuery := SearchJSONRequestBody{
		Depth:                  depth,
		Q:                      query,
		ExcludeDomains:         &options.ExcludeDomains,
		IncludeDomains:         &options.IncludeDomains,
		FromDate:               options.FromDate,
		ToDate:                 options.ToDate,
		IncludeImages:          &options.IncludeImages,
		MaxResults:             options.MaxResults,
		IncludeInlineCitations: &options.IncludeInlineCitations,
		IncludeSources:         &options.IncludeSources,
		StructuredOutputSchema: jsonSchema,
		OutputType:             Structured,
	}
	response, err := l.client.SearchWithResponse(context.Background(), searchQuery)
	if err != nil {
		return nil, err
	}
	if 200 <= response.StatusCode() && response.StatusCode() <= 299 {
		output := &StructuredOutput{}
		if options.IncludeSources {
			var sourcedOuput StructuredWithSourcesDto
			err := json.Unmarshal(response.Body, &sourcedOuput)
			if err != nil {
				return nil, err
			}
			output.SourcedOutput = &sourcedOuput
		} else {
			bodyStr := string(response.Body)
			output.RawJson = &bodyStr
		}
		return output, nil
	}
	return nil, fmt.Errorf("response returned a status code of %d: %s", response.StatusCode(), response.Status())
}

// Get the credit balance for the account associated with the API key the client are using
func (l *LinkupClient) GetBalance() (float32, error) {
	response, err := l.client.BalanceWithResponse(context.Background())
	if err != nil {
		return 0, err
	}
	if response.JSON200 != nil {
		return response.JSON200.Balance, nil
	}
	return 0, fmt.Errorf("response returned a status code of %d: %s", response.StatusCode(), response.Status())
}

func (l *LinkupClient) Fetch(
	url string,
	fetchOptions ...AdditionalFetchOptions,
) (*FetchOutput, error) {
	var options AdditionalFetchOptions
	switch len(fetchOptions) {
	case 0:
		options = DefaultAdditionalFetchOptions()
	default:
		options = fetchOptions[0]
	}
	fetchQuery := FetchJSONRequestBody{
		Url:            url,
		RenderJs:       &options.RenderJs,
		IncludeRawHtml: &options.IncludeRawHtml,
		ExtractImages:  &options.ExtractImages,
	}
	response, err := l.client.FetchWithResponse(context.Background(), fetchQuery)
	if err != nil {
		return nil, err
	}
	if response.JSON200 != nil {
		return response.JSON200, nil
	}
	return nil, fmt.Errorf("response returned a status code of %d: %s", response.StatusCode(), response.Status())
}
