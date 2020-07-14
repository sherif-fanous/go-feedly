package feedly

import (
	"net/http"

	"github.com/dghubble/sling"
	"github.com/sfanous/go-feedly/internal/mapstructure"
	"github.com/sfanous/go-feedly/pkg/time"
)

// RecommendationService provides methods for searching.
type RecommendationService struct {
	sling *sling.Sling
}

// newRecommendationService returns a new RecommendationService.
func newRecommendationService(sling *sling.Sling) *RecommendationService {
	return &RecommendationService{
		sling: sling,
	}
}

// Topic is a feedly topic.
type Topic struct {
	BestFeedID       *string                `json:"bestFeedId,omitempty"`
	DuplicateTopics  []Topic                `json:"duplicateTopics,omitempty"`
	FeedIDs          []string               `json:"feedIds,omitempty"`
	FeedInfos        []Feed                 `json:"feedInfos,omitempty"`
	FocusScore       *float64               `json:"focusScore,omitempty"`
	Language         *string                `json:"language,omitempty"`
	ParentTopic      *Topic                 `json:"parentTopic,omitempty"`
	RecommendedFeeds []Feed                 `json:"recommendedFeeds,omitempty"`
	RelatedTopics    []Topic                `json:"relatedTopics,omitempty"`
	Size             *int                   `json:"size,omitempty"`
	Topic            *string                `json:"topic,omitempty"`
	TopicID          *string                `json:"topicId,omitempty"`
	Updated          *time.Time             `json:"updated,omitempty"`
	VersionCode      *string                `json:"versionCode,omitempty"`
	Visual           *string                `json:"visual,omitempty"`
	UnmappedFields   map[string]interface{} `json:"-" mapstructure:",remain"`
}

// RecommendationTopicOptionalParams are the optional parameters for RecommendationService.Topic.
type RecommendationTopicOptionalParams struct {
	Count *int `url:"count,omitempty"`
}

// RecommendationTopicResponse represents the response from RecommendationService.Topic.
type RecommendationTopicResponse struct {
	Topics         []Topic                `json:"topics"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// Topic returns recommended feeds.
func (s *RecommendationService) Topic(query string, locale string, optionalParams *RecommendationTopicOptionalParams) (*RecommendationTopicResponse, *http.Response, error) {
	if optionalParams == nil {
		optionalParams = &RecommendationTopicOptionalParams{}
	}

	requiredParams := &struct {
		Locale string `url:"locale,omitempty"`
		Query  string `url:"query,omitempty"`
	}{
		Locale: locale,
		Query:  query,
	}

	encodedResponse := make([]map[string]interface{}, 0)
	decodedResponse := new(RecommendationTopicResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("recommendations/topics").QueryStruct(requiredParams).QueryStruct(optionalParams).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Topics); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}
