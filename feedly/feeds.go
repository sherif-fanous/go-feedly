package feedly

import (
	"net/http"
	"net/url"

	"github.com/dghubble/sling"
	"github.com/sfanous/go-feedly/internal/mapstructure"
	"github.com/sfanous/go-feedly/pkg/time"
)

// FeedService provides methods for managing feeds.
type FeedService struct {
	sling *sling.Sling
}

// newFeedService returns a new FeedService.
func newFeedService(sling *sling.Sling) *FeedService {
	return &FeedService{
		sling: sling,
	}
}

// Feed is a Feedly feed.
type Feed struct {
	AdPlatform                  *string                `json:"adPlatform,omitempty"`
	AdPosition                  *string                `json:"adPosition,omitempty"`
	AdSlotID                    *string                `json:"adSlotId,omitempty"`
	AccentColor                 *string                `json:"accentColor,omitempty"`
	Added                       *time.Time             `json:"added,omitempty"`
	AnalyticsEngine             *string                `json:"analyticsEngine,omitempty"`
	AnalyticsID                 *string                `json:"analyticsId,omitempty"`
	AverageReadTime             *float64               `json:"averageReadTime,omitempty"`
	ContentType                 *string                `json:"contentType,omitempty"`
	CoverColor                  *string                `json:"coverColor,omitempty"`
	CoverURL                    *string                `json:"coverUrl,omitempty"`
	Coverage                    *float64               `json:"coverage,omitempty"`
	CoverageScore               *float64               `json:"coverageScore,omitempty"`
	Curated                     *bool                  `json:"curated,omitempty"`
	DeliciousTags               []string               `json:"deliciousTags,omitempty"`
	Description                 *string                `json:"description,omitempty"`
	EstimatedEngagement         *int                   `json:"estimatedEngagement,omitempty"`
	Featured                    *bool                  `json:"featured,omitempty"`
	FeedID                      *string                `json:"feedId,omitempty"`
	IconURL                     *string                `json:"iconUrl,omitempty"`
	ID                          *string                `json:"id,omitempty"`
	Language                    *string                `json:"language,omitempty"`
	LastUpdated                 *time.Time             `json:"lastUpdated,omitempty"`
	LeoScore                    *float64               `json:"leoScore,omitempty"`
	Logo                        *string                `json:"logo,omitempty"`
	MustRead                    *bool                  `json:"mustRead,omitempty"`
	NumLongReadEntriesPastMonth *int                   `json:"numLongReadEntriesPastMonth,omitempty"`
	NumReadEntriesPastMonth     *int                   `json:"numReadEntriesPastMonth,omitempty"`
	NumTaggedEntriesPastMonth   *int                   `json:"numTaggedEntriesPastMonth,omitempty"`
	Partial                     *bool                  `json:"Partial,omitempty"`
	Promotion                   []string               `json:"promotion,omitempty"`
	RelatedLayout               *string                `json:"relatedLayout,omitempty"`
	RelatedTarget               *string                `json:"relatedTarget,omitempty"`
	RelevanceScore              *float64               `json:"relevanceScore,omitempty"`
	Score                       *int                   `json:"score,omitempty"`
	Sponsored                   *bool                  `json:"sponsored,omitempty"`
	State                       *string                `json:"state,omitempty"`
	Subscribers                 *int                   `json:"subscribers,omitempty"`
	TagCounts                   map[string]int         `json:"tagCounts,omitempty"`
	Title                       *string                `json:"title,omitempty"`
	Topics                      []string               `json:"topics,omitempty"`
	TotalReadingTimePastMonth   *int                   `json:"totalReadingTimePastMonth,omitempty"`
	TotalTagCount               *int                   `json:"totalTagCount,omitempty"`
	TwitterFollowers            *int                   `json:"twitterFollowers,omitempty"`
	TwitterScreenName           *string                `json:"twitterScreenName,omitempty"`
	Updated                     *time.Time             `json:"updated,omitempty"`
	Velocity                    *float64               `json:"velocity,omitempty"`
	VisualURL                   *string                `json:"visualUrl,omitempty"`
	Website                     *string                `json:"website,omitempty"`
	WebsiteTitle                *string                `json:"websiteTitle,omitempty"`
	Wordmark                    *string                `json:"wordmark,omitempty"`
	UnmappedFields              map[string]interface{} `json:"-" mapstructure:",remain"`
}

// FeedMetadataResponse represents the response from FeedService.Metadata.
type FeedMetadataResponse struct {
	Feed           *Feed                  `json:"feed"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// Metadata returns the metadata for a single feed.
func (s *FeedService) Metadata(feedID string) (*FeedMetadataResponse, *http.Response, error) {
	encodedResponse := make(map[string]interface{})
	decodedResponse := new(FeedMetadataResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("feeds/"+url.PathEscape(feedID)).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Feed); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// FeedMultipleMetadataResponse represents the response from FeedService.MultipleMetadata.
type FeedMultipleMetadataResponse struct {
	Feeds          []Feed                 `json:"feeds"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// MultipleMetadata returns the metadata for a list of feeds.
func (s *FeedService) MultipleMetadata(feedIDs []string) (*FeedMultipleMetadataResponse, *http.Response, error) {
	encodedResponse := make([]map[string]interface{}, 0)
	decodedResponse := new(FeedMultipleMetadataResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Post("feeds/.mget").BodyJSON(feedIDs).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Feeds); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}
