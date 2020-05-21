package feedly

import (
	"net/http"
	"net/url"

	"github.com/dghubble/sling"
	"github.com/sfanous/go-feedly/internal/mapstructure"
	"github.com/sfanous/go-feedly/pkg/time"
)

// MixService provides methods for managing mixes.
type MixService struct {
	sling *sling.Sling
}

// newMixService returns a new MixService.
func newMixService(sling *sling.Sling) *MixService {
	return &MixService{
		sling: sling,
	}
}

// MixMostEngagingOptionalParams are the optional parameters for MixService.MostEngaging.
type MixMostEngagingOptionalParams struct {
	Backfill   *bool      `url:"backfill,omitempty"`
	Count      *int       `url:"count,omitempty"`
	Hours      *int       `url:"hours,omitempty"`
	Locale     *string    `url:"locale,omitempty"`
	NewerThan  *time.Time `url:"newerThan,omitempty"`
	UnreadOnly *bool      `url:"unreadOnly,omitempty"`
}

// MixMostEngagingResponse represents the response from MixService.MostEngaging.
type MixMostEngagingResponse struct {
	Stream *Stream `json:"stream"`
}

// MostEngaging returns a mix of the most engaging content available in a stream.
func (s *MixService) MostEngaging(streamID string, optionalParams *MixMostEngagingOptionalParams) (*MixMostEngagingResponse, *http.Response, error) {
	if optionalParams == nil {
		optionalParams = &MixMostEngagingOptionalParams{}
	}

	encodedResponse := make(map[string]interface{})
	decodedResponse := new(MixMostEngagingResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("mixes/"+url.PathEscape(streamID)+"/contents").QueryStruct(optionalParams).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Stream); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}
