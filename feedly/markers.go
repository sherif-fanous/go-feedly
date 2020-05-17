package feedly

import (
	"net/http"

	"github.com/dghubble/sling"
	"github.com/sfanous/go-feedly/internal/mapstructure"
	"github.com/sfanous/go-feedly/pkg/time"
)

// MarkerService provides methods for managing markers.
type MarkerService struct {
	sling *sling.Sling
}

// newMarkerService returns a new MarkerService.
func newMarkerService(sling *sling.Sling) *MarkerService {
	return &MarkerService{
		sling: sling,
	}
}

type MarkAction string

const (
	KeepUnread     MarkAction = "keepUnread"
	MarkAsRead     MarkAction = "markAsRead"
	MarkAsSaved    MarkAction = "markAsSaved"
	MarkAsUnsaved  MarkAction = "markAsUnsaved"
	UndoMarkAsRead MarkAction = "undoMarkAsRead"
)

type MarkType string

const (
	Collections MarkType = "categories"
	Entries     MarkType = "entries"
	Feeds       MarkType = "feeds"
	Tags        MarkType = "tags"
)

// MarkerLatestReadOptionalParams are the optional parameters for MarkerService.LatestRead.
type MarkerLatestReadOptionalParams struct {
	NewerThan *time.Time `url:"newerThan,omitempty"`
}

// MarkerLatestReadResponse represents the response from MarkerService.LatestRead.
type MarkerLatestReadResponse struct {
	Entries []string `json:"entries,omitempty"`
	Feeds   []struct {
		AsOf *time.Time `json:"asOf,omitempty"`
		ID   *string    `json:"id,omitempty"`
	} `json:"feeds,omitempty"`
	Unread         []string               `json:"unread,omitempty"`
	Updated        *time.Time             `json:"updated,omitempty"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// LatestRead returns the latest read operations.
func (s *MarkerService) LatestRead(optionalParams *MarkerLatestReadOptionalParams) (*MarkerLatestReadResponse, *http.Response, error) {
	if optionalParams == nil {
		optionalParams = &MarkerLatestReadOptionalParams{}
	}

	encodedResponse := make(map[string]interface{})
	decodedResponse := new(MarkerLatestReadResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("markers/reads").QueryStruct(optionalParams).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, decodedResponse); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// MarkerLatestTaggedOptionalParams are the optional parameters for MarkerService.LatestTagged.
type MarkerLatestTaggedOptionalParams struct {
	NewerThan *time.Time `url:"newerThan,omitempty"`
}

// MarkerLatestTaggedResponse represents the response from MarkerService.LatestTagged.
type MarkerLatestTaggedResponse struct {
	TaggedEntries  map[string][]string    `json:"taggedEntries,omitempty"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// LatestTagged returns latest tagged entry ids.
func (s *MarkerService) LatestTagged(optionalParams *MarkerLatestTaggedOptionalParams) (*MarkerLatestTaggedResponse, *http.Response, error) {
	if optionalParams == nil {
		optionalParams = &MarkerLatestTaggedOptionalParams{}
	}

	encodedResponse := make(map[string]interface{})
	decodedResponse := new(MarkerLatestTaggedResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("markers/tags").QueryStruct(optionalParams).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, decodedResponse); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// MarkerMarkOptionalParams are the optional parameters for MarkerService.Mark.
type MarkerMarkOptionalParams struct {
	AsOf            *time.Time `json:"asOf,omitempty"`
	CollectionIDs   []string   `json:"categoryIds,omitempty"`
	EntryIDs        []string   `json:"entryIds,omitempty"`
	FeedIDs         []string   `json:"feedIds,omitempty"`
	LastReadEntryID *string    `json:"lastReadEntryId,omitempty"`
	TagIDs          []string   `json:"tagIds,omitempty"`
}

// Mark marks one or more collections, entries, feeds, or tags as read, saved, or unread.
func (s *MarkerService) Mark(markAction MarkAction, markType MarkType, optionalParams *MarkerMarkOptionalParams) (*http.Response, error) {
	bodyJSON := &struct {
		Action MarkAction `json:"action"`
		MarkerMarkOptionalParams
		Type MarkType `json:"type"`
	}{
		Action:                   markAction,
		MarkerMarkOptionalParams: MarkerMarkOptionalParams{},
		Type:                     markType,
	}

	switch markType {
	case Collections:
		switch markAction {
		case MarkAsRead:
			bodyJSON.AsOf = optionalParams.AsOf
			bodyJSON.CollectionIDs = optionalParams.CollectionIDs
			bodyJSON.LastReadEntryID = optionalParams.LastReadEntryID
		case UndoMarkAsRead:
			bodyJSON.CollectionIDs = optionalParams.CollectionIDs
		}
	case Entries:
		bodyJSON.EntryIDs = optionalParams.EntryIDs
	case Feeds:
		switch markAction {
		case MarkAsRead:
			bodyJSON.AsOf = optionalParams.AsOf
			bodyJSON.FeedIDs = optionalParams.FeedIDs
			bodyJSON.LastReadEntryID = optionalParams.LastReadEntryID
		case UndoMarkAsRead:
			bodyJSON.FeedIDs = optionalParams.FeedIDs
		}
	case Tags:
		bodyJSON.AsOf = optionalParams.AsOf
		bodyJSON.LastReadEntryID = optionalParams.LastReadEntryID
		bodyJSON.TagIDs = optionalParams.TagIDs
	}

	apiError := new(APIError)

	resp, err := s.sling.Post("markers").BodyJSON(bodyJSON).Receive(nil, apiError)

	return resp, relevantError(err, apiError)
}

// MarkerUnreadCountsOptionalParams are the optional parameters for MarkerService.UnreadCounts.
type MarkerUnreadCountsOptionalParams struct {
	AutoRefresh *bool      `url:"autorefresh,omitempty"`
	NewerThan   *time.Time `url:"newerThan,omitempty"`
	StreamID    *string    `url:"streamId,omitempty"`
}

// MarkerUnreadCountsResponse represents the response from MarkerService.UnreadCounts.
type MarkerUnreadCountsResponse struct {
	UnreadCounts []struct {
		ID             *string                `json:"id,omitempty"`
		Count          *int                   `json:"count,omitempty"`
		Updated        *time.Time             `json:"updated,omitempty"`
		UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
	} `json:"unreadcounts,omitempty"`
	Updated        *time.Time             `json:"updated,omitempty"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// UnreadCounts returns the list of unread counts.
func (s *MarkerService) UnreadCounts(optionalParams *MarkerUnreadCountsOptionalParams) (*MarkerUnreadCountsResponse, *http.Response, error) {
	if optionalParams == nil {
		optionalParams = &MarkerUnreadCountsOptionalParams{}
	}

	encodedResponse := make(map[string]interface{})
	decodedResponse := new(MarkerUnreadCountsResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("markers/counts").QueryStruct(optionalParams).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, decodedResponse); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}
