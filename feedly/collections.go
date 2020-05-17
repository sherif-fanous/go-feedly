package feedly

import (
	"io"
	"net/http"
	"net/url"

	"github.com/dghubble/sling"
	"github.com/sfanous/go-feedly/internal/mapstructure"
	"github.com/sfanous/go-feedly/internal/mime"
	"github.com/sfanous/go-feedly/pkg/time"
)

// CollectionService provides methods for managing collections of feed subscriptions, aka categories.
type CollectionService struct {
	sling *sling.Sling
}

// newCollectionService returns a new CollectionService.
func newCollectionService(sling *sling.Sling) *CollectionService {
	return &CollectionService{
		sling: sling,
	}
}

// Collection is a Feedly collection.
type Collection struct {
	ACL []struct {
		Scope  *string `json:"scope,omitempty"`
		Target *string `json:"target,omitempty"`
	}
	Cover          *string                `json:"cover,omitempty"`
	Created        *time.Time             `json:"created,omitempty"`
	Customizable   *bool                  `json:"customizable,omitempty"`
	Description    *string                `json:"description,omitempty"`
	Engagement     *int                   `json:"engagement,omitempty"`
	Enterprise     *bool                  `json:"enterprise,omitempty"`
	Feeds          []Feed                 `json:"feeds,omitempty"`
	ID             *string                `json:"id,omitempty"`
	Label          *string                `json:"label,omitempty"`
	NumFeeds       *int                   `json:"numFeeds,omitempty"`
	Topics         []string               `json:"topics,omitempty"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// CollectionAddFeedResponse represents the response from CollectionService.AddFeed.
type CollectionAddFeedResponse struct {
	Feeds []Feed `json:"feeds"`
}

// AddFeed adds a feed to an existing collection.
func (s *CollectionService) AddFeed(collectionID string, feed *Feed) (*CollectionAddFeedResponse, *http.Response, error) {
	encodedResponse := make([]map[string]interface{}, 0)
	decodedResponse := new(CollectionAddFeedResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Put("collections/"+url.PathEscape(collectionID)+"/feeds").BodyJSON(feed).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Feeds); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// CollectionAddMultipleFeedsResponse represents the response from CollectionService.AddMultipleFeeds.
type CollectionAddMultipleFeedsResponse struct {
	Feeds []Feed `json:"feeds"`
}

// AddMultipleFeeds adds a one or more feeds to an existing collection.
func (s *CollectionService) AddMultipleFeeds(collectionID string, feeds []Feed) (*CollectionAddMultipleFeedsResponse, *http.Response, error) {
	encodedResponse := make([]map[string]interface{}, 0)
	decodedResponse := new(CollectionAddMultipleFeedsResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Post("collections/"+url.PathEscape(collectionID)+"/feeds/.mput").BodyJSON(feeds).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Feeds); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// CollectionCreateOptionalParams are the optional parameters for CollectionService.Create.
type CollectionCreateOptionalParams struct {
	Description *string `json:"description,omitempty"`
	Feeds       []Feed  `json:"feeds,omitempty"`
	ID          *string `json:"id,omitempty"`
}

// CollectionCreateResponse represents the response from CollectionService.Create.
type CollectionCreateResponse struct {
	Collections []Collection `json:"collections"`
}

// Create creates a new collection.
func (s *CollectionService) Create(label string, optionalParams *CollectionCreateOptionalParams) (*CollectionCreateResponse, *http.Response, error) {
	if optionalParams == nil {
		optionalParams = &CollectionCreateOptionalParams{}
	}

	bodyJSON := &struct {
		*CollectionCreateOptionalParams
		Label string `json:"label"`
	}{
		CollectionCreateOptionalParams: optionalParams,
		Label:                          label,
	}

	encodedResponse := make([]map[string]interface{}, 0)
	decodedResponse := new(CollectionCreateResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Post("collections").BodyJSON(bodyJSON).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Collections); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// Delete deletes an existing collection.
func (s *CollectionService) Delete(collectionID string) (*http.Response, error) {
	apiError := new(APIError)

	resp, err := s.sling.New().Delete("collections/"+url.PathEscape(collectionID)).Receive(nil, apiError)

	return resp, relevantError(err, apiError)
}

// CollectionDeleteFeedOptionalParams are the optional parameters for CollectionService.DeleteFeed.
type CollectionDeleteFeedOptionalParams struct {
	KeepOrphanFeeds *bool `url:"keepOrphanFeeds,omitempty"`
}

// DeleteFeed deletes a feed from an existing collection.
func (s *CollectionService) DeleteFeed(collectionID string, feedID string, optionalParams *CollectionDeleteFeedOptionalParams) (*http.Response, error) {
	if optionalParams == nil {
		optionalParams = &CollectionDeleteFeedOptionalParams{}
	}

	apiError := new(APIError)

	resp, err := s.sling.New().Delete("collections/"+url.PathEscape(collectionID)+"/feeds/"+url.PathEscape(feedID)).QueryStruct(optionalParams).Receive(nil, apiError)

	return resp, relevantError(err, apiError)
}

// CollectionDeleteMultipleFeedsOptionalParams are the optional parameters for CollectionService.DeleteMultipleFeeds.
type CollectionDeleteMultipleFeedsOptionalParams CollectionDeleteFeedOptionalParams

// DeleteMultipleFeeds removes one or more feeds from an existing collection.
func (s *CollectionService) DeleteMultipleFeeds(collectionID string, feedIDs []string, optionalParams *CollectionDeleteFeedOptionalParams) (*http.Response, error) {
	if optionalParams == nil {
		optionalParams = &CollectionDeleteFeedOptionalParams{}
	}

	apiError := new(APIError)

	resp, err := s.sling.New().Delete("collections/"+url.PathEscape(collectionID)+"/feeds/.mdelete").QueryStruct(optionalParams).BodyJSON(feedIDs).Receive(nil, apiError)

	return resp, relevantError(err, apiError)
}

// CollectionDetailResponse represents the response from CollectionService.Details.
type CollectionDetailResponse struct {
	Collections []Collection `json:"collections"`
}

// Details returns details about a collection.
func (s *CollectionService) Details(collectionID string) (*CollectionDetailResponse, *http.Response, error) {
	encodedResponse := make([]map[string]interface{}, 0)
	decodedResponse := new(CollectionDetailResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("collections/"+url.PathEscape(collectionID)).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Collections); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// CollectionListOptionalParams are the optional parameters for CollectionService.List.
type CollectionListOptionalParams struct {
	WithEnterprise *bool `url:"withEnterprise,omitempty"`
	WithStatus     *bool `url:"withStats,omitempty"`
}

// CollectionListResponse represents the response from CollectionService.List.
type CollectionListResponse struct {
	Collections []Collection `json:"collections"`
}

// List returns the list of collections.
func (s *CollectionService) List(optionalParams *CollectionListOptionalParams) (*CollectionListResponse, *http.Response, error) {
	if optionalParams == nil {
		optionalParams = &CollectionListOptionalParams{}
	}

	encodedResponse := make([]map[string]interface{}, 0)
	decodedResponse := new(CollectionListResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("collections").QueryStruct(optionalParams).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Collections); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// CollectionUpdateOptionalParams are the optional parameters for CollectionService.Update.
type CollectionUpdateOptionalParams struct {
	DeleteCover *bool   `json:"deleteCover,omitempty"`
	Description *string `json:"description,omitempty"`
	Feeds       []Feed  `json:"feeds,omitempty"`
	Label       *string `json:"label,omitempty"`
}

// CollectionUpdateResponse represents the response from CollectionService.Update.
type CollectionUpdateResponse struct {
	Collections []Collection `json:"collections"`
}

// Update updates an existing collection.
func (s *CollectionService) Update(collectionID string, optionalParams *CollectionUpdateOptionalParams) (*CollectionUpdateResponse, *http.Response, error) {
	if optionalParams == nil {
		optionalParams = &CollectionUpdateOptionalParams{}
	}

	bodyJSON := &struct {
		*CollectionUpdateOptionalParams
		ID string `json:"id,omitempty"`
	}{
		CollectionUpdateOptionalParams: optionalParams,
		ID:                             collectionID,
	}

	encodedResponse := make([]map[string]interface{}, 0)
	decodedResponse := new(CollectionUpdateResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Post("collections").BodyJSON(bodyJSON).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Collections); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// CollectionUploadCoverImageResponse represents the response from CollectionService.UploadCoverImage.
type CollectionUploadCoverImageResponse struct {
	Collections []Collection `json:"collections"`
}

// UploadCoverImage uploads a new cover image for an existing collection.
func (s *CollectionService) UploadCoverImage(collectionID string, coverImage io.Reader) (*CollectionUploadCoverImageResponse, *http.Response, error) {
	body, contentType, err := mime.CreateMultipartMIMEAttachment(coverImage)
	if err != nil {
		return nil, nil, err
	}

	encodedResponse := make([]map[string]interface{}, 0)
	decodedResponse := new(CollectionUploadCoverImageResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Post("collections/"+url.PathEscape(collectionID)).Body(body).Set("Content-Type", contentType).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Collections); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}
