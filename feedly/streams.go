package feedly

import (
	"net/http"
	"net/url"

	"github.com/dghubble/sling"
	"github.com/sfanous/go-feedly/internal/mapstructure"
	"github.com/sfanous/go-feedly/pkg/time"
)

// StreamService provides methods for managing streams.
type StreamService struct {
	sling *sling.Sling
}

// newStreamService returns a new StreamService.
func newStreamService(sling *sling.Sling) *StreamService {
	return &StreamService{
		sling: sling,
	}
}

type ContentRank string

const (
	Engagement ContentRank = "EngagementFilter"
	Newest     ContentRank = "newest"
	Oldest     ContentRank = "oldest"
)

// Stream is a Feedly stream.
type Stream struct {
	Alternate []struct {
		HRef           *string                `json:"href,omitempty"`
		Type           *string                `json:"type,omitempty"`
		UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
	} `json:"alternate,omitempty"`
	Continuation   *string                `json:"continuation,omitempty"`
	Direction      *string                `json:"direction,omitempty"`
	ID             *string                `json:"id,omitempty"`
	Items          []Entry                `json:"items,omitempty"`
	Title          *string                `json:"title,omitempty"`
	Updated        *time.Time             `json:"updated,omitempty"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// StreamContentOptionalParams are the optional parameters for StreamService.Content.
type StreamContentOptionalParams struct {
	Continuation      *string      `url:"count,omitempty"`
	Count             *int         `url:"count,omitempty"`
	FindURLDuplicates *bool        `url:"findUrlDuplicates,omitempty"`
	ImportantOnly     *bool        `url:"importantOnly,omitempty"`
	NewerThan         *time.Time   `url:"newerThan,omitempty"`
	Ranked            *ContentRank `url:"ranked,omitempty"`
	ShowMuted         *bool        `url:"showMuted,omitempty"`
	Similar           *bool        `url:"similar,omitempty"`
	UnreadOnly        *bool        `url:"unreadOnly,omitempty"`
}

// StreamContentResponse represents the response from StreamService.Content.
type StreamContentResponse struct {
	Stream *Stream `json:"stream"`
}

// Content returns the content of a stream.
func (s *StreamService) Content(streamID string, optionalParams *StreamContentOptionalParams) (*StreamContentResponse, *http.Response, error) {
	if optionalParams == nil {
		optionalParams = &StreamContentOptionalParams{}
	}

	encodedResponse := make(map[string]interface{})
	decodedResponse := new(StreamContentResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("streams/"+url.PathEscape(streamID)+"/contents").QueryStruct(optionalParams).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Stream); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// StreamEntryIDsOptionalParams are the optional parameters for StreamService.EntryIDs.
type StreamEntryIDsOptionalParams struct {
	Continuation *string      `url:"count,omitempty"`
	Count        *int         `url:"count,omitempty"`
	NewerThan    *time.Time   `url:"newerThan,omitempty"`
	Ranked       *ContentRank `url:"ranked,omitempty"`
	UnreadOnly   *bool        `url:"unreadOnly,omitempty"`
}

// StreamEntryIDsResponse represents the response from StreamService.EntryIDs.
type StreamEntryIDsResponse struct {
	Continuation   *string                `json:"continuation,omitempty"`
	IDs            []string               `json:"ids,omitempty" mapstructure:"ids"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// EntryIDs returns the IDs of entries in a stream.
func (s *StreamService) EntryIDs(streamID string, optionalParams *StreamEntryIDsOptionalParams) (*StreamEntryIDsResponse, *http.Response, error) {
	if optionalParams == nil {
		optionalParams = &StreamEntryIDsOptionalParams{}
	}

	encodedResponse := make(map[string]interface{})
	decodedResponse := new(StreamEntryIDsResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("streams/"+url.PathEscape(streamID)+"/ids").QueryStruct(optionalParams).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, decodedResponse); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}
