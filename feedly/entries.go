package feedly

import (
	"net/http"
	"net/url"

	"github.com/dghubble/sling"
	"github.com/sfanous/go-feedly/internal/mapstructure"
	"github.com/sfanous/go-feedly/pkg/time"
)

// EntryService provides methods for managing entries.
type EntryService struct {
	sling *sling.Sling
}

// newEntryService returns a new EntryService.
func newEntryService(sling *sling.Sling) *EntryService {
	return &EntryService{
		sling: sling,
	}
}

// Entry is a Feedly entry.
type Entry struct {
	ActionTimestamp *time.Time `json:"actionTimestamp,omitempty"`
	Alternate       []struct {
		HRef           *string                `json:"href,omitempty"`
		Title          *string                `json:"title,omitempty"`
		Type           *string                `json:"type,omitempty"`
		UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
	} `json:"alternate,omitempty"`
	AMPURL                 *string `json:"ampUrl,omitempty"`
	AnalysisFeedbackPrompt *struct {
		Confidence     *float64               `json:"confidence,omitempty"`
		ID             *string                `json:"id,omitempty"`
		Label          *string                `json:"label,omitempty"`
		Type           *string                `json:"type,omitempty"`
		UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
	} `json:"analysisFeedbackPrompt,omitempty"`
	Author    *string `json:"author,omitempty"`
	CDNAmpURL *string `json:"cdnAmpUrl,omitempty"`
	Canonical []struct {
		HRef           *string                `json:"href,omitempty"`
		Type           *string                `json:"type,omitempty"`
		UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
	} `json:"canonical,omitempty"`
	CanonicalURL *string      `json:"canonicalUrl,omitempty"`
	Categories   []Collection `json:"categories,omitempty"`
	CommonTopics []struct {
		ID             *string                `json:"id,omitempty"`
		Label          *string                `json:"label,omitempty"`
		SalienceLevel  *string                `json:"salienceLevel,omitempty"`
		Score          *float64               `json:"score,omitempty"`
		UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
	} `json:"commonTopics,omitempty"`
	Content *struct {
		Content        *string                `json:"content,omitempty"`
		Direction      *string                `json:"direction,omitempty"`
		UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
	} `json:"content,omitempty"`
	Crawled *time.Time `json:"crawled,omitempty"`
	Created *struct {
		Application    *string                `json:"application,omitempty"`
		UserAgent      *string                `json:"userAgent,omitempty"`
		UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
	} `json:"created,omitempty"`
	CreatedBy *struct {
		Application    *string                `json:"application,omitempty"`
		Client         *string                `json:"client,omitempty"`
		UserAgent      *string                `json:"userAgent,omitempty"`
		UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
	} `json:"createdBy,omitempty"`
	Enclosure []struct {
		HRef           *string                `json:"href,omitempty"`
		Height         *int                   `json:"height,omitempty"`
		Length         *int                   `json:"Length,omitempty"`
		Title          *string                `json:"title,omitempty"`
		Type           *string                `json:"type,omitempty"`
		Width          *int                   `json:"width,omitempty"`
		UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
	} `json:"enclosure,omitempty"`
	Engagement     *int     `json:"EngagementFilter,omitempty"`
	EngagementRate *float64 `json:"engagementRate,omitempty"`
	Entities       []struct {
		ID       *string `json:"id,omitempty"`
		Label    *string `json:"label,omitempty"`
		Mentions []struct {
			Text           *string                `json:"text,omitempty"`
			UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
		} `json:"mentions,omitempty"`
		SalienceLevel  *string                `json:"salienceLevel,omitempty"`
		UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
	} `json:"entities,omitempty"`
	Fingerprint *string  `json:"fingerprint,omitempty"`
	ID          *string  `json:"id,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
	Language    *string  `json:"language,omitempty"`
	Memes       []struct {
		ID             *string                `json:"id,omitempty"`
		Label          *string                `json:"label,omitempty"`
		Score          *float64               `json:"score,omitempty"`
		UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
	} `json:"memes,omitempty"`
	Origin *struct {
		HTMLURL        *string                `json:"htmlUrl,omitempty"`
		StreamID       *string                `json:"streamId,omitempty"`
		Title          *string                `json:"title,omitempty"`
		UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
	} `json:"origin,omitempty"`
	OriginID   *string                  `json:"originId,omitempty"`
	Priorities []map[string]interface{} `json:"priorities,omitempty"`
	Published  *time.Time               `json:"published,omitempty"`
	Recrawled  *time.Time               `json:"recrawled,omitempty"`
	SID        *string                  `json:"sid,omitempty"`
	Summary    *struct {
		Content        *string                `json:"content,omitempty"`
		Direction      *string                `json:"direction,omitempty"`
		UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
	} `json:"summary,omitempty"`
	Tags      []Board `json:"tags,omitempty"`
	Thumbnail []struct {
		Height         *int                   `json:"height,omitempty"`
		URL            *string                `json:"url,omitempty"`
		Width          *int                   `json:"width,omitempty"`
		UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
	} `json:"thumbnail,omitempty"`
	Title       *string    `json:"title,omitempty"`
	Unread      *bool      `json:"unread,omitempty"`
	UpdateCount *int       `json:"updateCount,omitempty"`
	Updated     *time.Time `json:"updated,omitempty"`
	Visual      *struct {
		ContentType    *string                `json:"contentType,omitempty"`
		EdgeCacheURL   *string                `json:"edgeCacheUrl,omitempty"`
		ExpirationDate *time.Time             `json:"expirationDate,omitempty"`
		Height         *int                   `json:"height,omitempty"`
		Processor      *string                `json:"processor,omitempty"`
		URL            *string                `json:"url,omitempty"`
		Width          *int                   `json:"width,omitempty"`
		UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
	} `json:"visual,omitempty"`
	Webfeeds *struct {
		AccentColor     *string                `json:"accentColor,omitempty"`
		AnalyticsEngine *string                `json:"analyticsEngine,omitempty"`
		AnalyticsID     *string                `json:"analyticsId,omitempty"`
		CoverImage      *string                `json:"coverImage,omitempty"`
		Icon            *string                `json:"icon,omitempty"`
		Logo            *string                `json:"logo,omitempty"`
		Partial         *bool                  `json:"partial,omitempty"`
		Promotion       []string               `json:"promotion,omitempty"`
		RelatedLayout   *string                `json:"relatedLayout,omitempty"`
		RelatedTarget   *string                `json:"relatedTarget,omitempty"`
		Wordmark        *string                `json:"wordmark,omitempty"`
		UnmappedFields  map[string]interface{} `json:"-" mapstructure:",remain"`
	} `json:"webfeeds,omitempty"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// EntryContentResponse represents the response from EntryService.Content.
type EntryContentResponse struct {
	Entries        []Entry                `json:"entries"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// Content returns the content of an entry.
func (s *EntryService) Content(entryID string) (*EntryContentResponse, *http.Response, error) {
	encodedResponse := make([]map[string]interface{}, 0)
	decodedResponse := new(EntryContentResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("entries/"+url.PathEscape(entryID)).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Entries); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// EntryCreateResponse represents the response from EntryService.Create.
type EntryCreateResponse struct {
	EntryIDs       []string               `json:"entryIds"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// Create creates and tags an entry.
func (s *EntryService) Create(entry *Entry) (*EntryCreateResponse, *http.Response, error) {
	entryIDs := new(EntryCreateResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Post("entries").BodyJSON(entry).Receive(&entryIDs.EntryIDs, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	return entryIDs, resp, nil
}

// EntryMultipleContentResponse represents the response from EntryService.MultipleContent.
type EntryMultipleContentResponse struct {
	Entries        []Entry                `json:"entries"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// MultipleContent returns the content for one or more entries.
func (s *EntryService) MultipleContent(entryIDs []string) (*EntryMultipleContentResponse, *http.Response, error) {
	encodedResponse := make([]map[string]interface{}, 0)
	decodedResponse := new(EntryMultipleContentResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Post("entries/.mget").BodyJSON(entryIDs).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Entries); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}
