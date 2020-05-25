package feedly

import (
	"net/http"

	"github.com/dghubble/sling"
	"github.com/sfanous/go-feedly/internal/mapstructure"
	"github.com/sfanous/go-feedly/pkg/time"
)

// ProfileService provides methods for managing profile information.
type ProfileService struct {
	sling *sling.Sling
}

// newProfileService returns a new ProfileService.
func newProfileService(sling *sling.Sling) *ProfileService {
	return &ProfileService{
		sling: sling,
	}
}

// Profile is a Feedly user profile.
type Profile struct {
	AnonymizedHash      *string    `json:"anonymizedHash,omitempty"`
	Client              *string    `json:"client,omitempty"`
	CohortGroups        []string   `json:"cohortGroups,omitempty"`
	Cohorts             []string   `json:"cohorts,omitempty"`
	Created             *time.Time `json:"created,omitempty"`
	CustomFamilyName    *string    `json:"customFamilyName,omitempty"`
	CustomGivenName     *string    `json:"customGivenName,omitempty"`
	DropboxConnected    *bool      `json:"dropboxConnected,omitempty"`
	Email               *string    `json:"email"`
	EvernoteConnected   *bool      `json:"evernoteConnected,omitempty"`
	Facebook            *string    `json:"facebook,omitempty"`
	FacebookConnected   *bool      `json:"facebookConnected,omitempty"`
	FamilyName          *string    `json:"familyName,omitempty"`
	FullName            *string    `json:"fullName,omitempty"`
	Gender              *string    `json:"gender,omitempty"`
	GivenName           *string    `json:"givenName,omitempty"`
	Google              *string    `json:"google,omitempty"`
	ID                  *string    `json:"id,omitempty"`
	InstapaperConnected *bool      `json:"instapaperConnected,omitempty"`
	LandingPage         *string    `json:"landingPage,omitempty"`
	Locale              *string    `json:"locale,omitempty"`
	Login               *string    `json:"login,omitempty"`
	Logins              []struct {
		ID         *string `json:"id,omitempty"`
		Provider   *string `json:"provider,omitempty"`
		ProviderID *string `json:"providerId,omitempty"`
		Verified   *bool   `json:"verified,omitempty"`
	} `json:"logins,omitempty"`
	Picture              *string                `json:"picture,omitempty"`
	PocketConnected      *bool                  `json:"pocketConnected,omitempty"`
	Reader               *string                `json:"reader,omitempty"`
	RefPage              *string                `json:"refPage,omitempty"`
	Source               *string                `json:"source,omitempty"`
	Twitter              *string                `json:"twitter,omitempty"`
	TwitterConnected     *bool                  `json:"twitterConnected,omitempty"`
	Verified             *bool                  `json:"verified,omitempty"`
	Wave                 *string                `json:"wave,omitempty"`
	WindowsLiveConnected *bool                  `json:"windowsLiveConnected,omitempty"`
	WordPressConnected   *bool                  `json:"wordPressConnected,omitempty"`
	UnmappedFields       map[string]interface{} `json:"-" mapstructure:",remain"`
}

// ProfileListResponse represents the response from ProfileService.List.
type ProfileListResponse struct {
	Profile *Profile `json:"profile"`
}

// List returns the profile of the user.
func (s *ProfileService) List() (*ProfileListResponse, *http.Response, error) {
	encodedResponse := make(map[string]interface{})
	decodedResponse := new(ProfileListResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("profile").Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Profile); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// Update updates the profile of the user.
func (s *ProfileService) Update(profile *Profile) (*http.Response, error) {
	apiError := new(APIError)

	resp, err := s.sling.New().Post("profile").BodyJSON(profile).Receive(nil, apiError)

	return resp, relevantError(err, apiError)
}
