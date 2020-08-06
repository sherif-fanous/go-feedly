package feedly

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/dghubble/sling"
	"github.com/sfanous/go-feedly/internal/mapstructure"
	"github.com/sfanous/go-feedly/pkg/time"
)

// SearchService provides methods for searching.
type SearchService struct {
	sling *sling.Sling
}

// newSearchService returns a new SearchService.
func newSearchService(sling *sling.Sling) *SearchService {
	return &SearchService{
		sling: sling,
	}
}

type EmbeddedFilter string

const (
	Any   EmbeddedFilter = "any"
	Audio EmbeddedFilter = "audio"
	Doc   EmbeddedFilter = "doc"
	Video EmbeddedFilter = "video"
)

type EngagementFilter string

const (
	High   EngagementFilter = "high"
	Medium EngagementFilter = "medium"
)

type FieldFilter struct {
	All      bool
	Author   bool
	Keywords bool
	Title    bool
}

// EncodeValues implements the query.Encoder interface.
func (f *FieldFilter) EncodeValues(key string, v *url.Values) error {
	if f.All {
		v.Set(key, "all")

		return nil
	}

	fields := make([]string, 0)

	if f.Author {
		fields = append(fields, "author")
	}

	if f.Keywords {
		fields = append(fields, "keywords")
	}

	if f.Title {
		fields = append(fields, "title")
	}

	v.Set(key, strings.Join(fields, ","))

	return nil
}

// SearchFeedsOptionalParams are the optional parameters for SearchService.Topic.
type SearchFeedsOptionalParams struct {
	Count  *int    `url:"count,omitempty"`
	Locale *string `url:"locale,omitempty"`
}

// SearchFeedsResponse represents the response from SearchService.Topic.
type SearchFeedsResponse struct {
	Hint           *string                `json:"hint,omitempty"`
	QueryType      *string                `json:"queryType,omitempty"`
	Related        []string               `json:"related,omitempty"`
	Results        []Feed                 `json:"results,omitempty"`
	Scheme         *string                `json:"scheme,omitempty"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// Feeds returns matching feeds.
func (s *SearchService) Feeds(query string, optionalParams *SearchFeedsOptionalParams) (*SearchFeedsResponse, *http.Response, error) {
	if optionalParams == nil {
		optionalParams = &SearchFeedsOptionalParams{}
	}

	requiredParams := &struct {
		Query string `url:"query,omitempty"`
	}{
		Query: query,
	}

	encodedResponse := make(map[string]interface{})
	decodedResponse := new(SearchFeedsResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("search/feeds").QueryStruct(requiredParams).QueryStruct(optionalParams).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// SearchStreamOptionalParams are the optional parameters for SearchService.Stream.
type SearchStreamOptionalParams struct {
	Continuation *string           `url:"count,omitempty"`
	Count        *int              `url:"count,omitempty"`
	Embedded     *EmbeddedFilter   `url:"embedded,omitempty"`
	Engagement   *EngagementFilter `url:"engagement,omitempty"`
	Fields       *FieldFilter      `url:"fields,omitempty"`
	Locale       *string           `url:"locale,omitempty"`
	NewerThan    *time.Time        `url:"newerThan,omitempty"`
	UnreadOnly   *bool             `url:"unreadOnly,omitempty"`
}

// SearchStreamResponse represents the response from SearchService.Stream.
type SearchStreamResponse struct {
	AdvancedSearch    *bool `json:"advancedSearch,omitempty"`
	SearchElapsedTime *int  `json:"searchElapsedTime,omitempty"`
	SearchTime        *int  `json:"searchTime,omitempty"`
	Stream            `mapstructure:",squash"`
	Terms             []string               `json:"terms,omitempty"`
	UnmappedFields    map[string]interface{} `json:"-" mapstructure:",remain"`
}

// Stream returns matching content in a stream.
func (s *SearchService) Stream(streamID string, query string, optionalParams *SearchStreamOptionalParams) (*SearchStreamResponse, *http.Response, error) {
	if optionalParams == nil {
		optionalParams = &SearchStreamOptionalParams{}
	}

	requiredParams := &struct {
		Query    string `url:"query,omitempty"`
		StreamID string `url:"streamId,omitempty"`
	}{
		Query:    query,
		StreamID: streamID,
	}

	encodedResponse := make(map[string]interface{})
	decodedResponse := new(SearchStreamResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("search/contents").QueryStruct(requiredParams).QueryStruct(optionalParams).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}
