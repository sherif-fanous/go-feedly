package feedly

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

const (
	// APIBaseURL is the base URL for the Feedly API.
	APIBaseURL = "https://cloud.feedly.com"
	// APIBaseVersion is the base version for the Feedly API.
	APIBaseVersion = "v3"
)

// A Client is a Feedly API client. Its zero value is not a usable Feedly client.
type Client struct {
	apiBaseURL     string
	apiBaseVersion string
	sling          *sling.Sling
	// Feedly API Services
	Boards          *BoardService
	Collections     *CollectionService
	Entries         *EntryService
	Feeds           *FeedService
	Library         *LibraryService
	Markers         *MarkerService
	Mixes           *MixService
	OPML            *OPMLService
	Preferences     *PreferenceService
	Profile         *ProfileService
	Recommendations *RecommendationService
	Search          *SearchService
	Streams         *StreamService
}

// WithAPIBaseURL returns a function that initializes a Client with an API base URL.
func WithAPIBaseURL(apiBaseURL string) func(*Client) {
	return func(c *Client) {
		c.apiBaseURL = apiBaseURL
	}
}

// WithAPIBaseVersion returns a function that initializes a Client with an API version.
func WithAPIBaseVersion(apiBaseVersion string) func(*Client) {
	return func(c *Client) {
		c.apiBaseVersion = apiBaseVersion
	}
}

// NewClient returns a new Client.
func NewClient(httpClient *http.Client, optionalParameters ...func(*Client)) *Client {
	client := &Client{
		apiBaseURL:     APIBaseURL,
		apiBaseVersion: APIBaseVersion,
	}

	for _, optionalParameter := range optionalParameters {
		optionalParameter(client)
	}

	base := sling.New().Client(httpClient).Base(fmt.Sprintf("%s/%s/", client.apiBaseURL, client.apiBaseVersion))

	client.sling = base
	client.Boards = newBoardService(base.New())
	client.Collections = newCollectionService(base.New())
	client.Entries = newEntryService(base.New())
	client.Feeds = newFeedService(base.New())
	client.Library = newLibraryService(base.New())
	client.Markers = newMarkerService(base.New())
	client.Mixes = newMixService(base.New())
	client.OPML = newOPMLService(base.New())
	client.Preferences = newPreferenceService(base.New())
	client.Profile = newProfileService(base.New())
	client.Search = newSearchService(base.New())
	client.Recommendations = newRecommendationService(base.New())
	client.Streams = newStreamService(base.New())

	return client
}
