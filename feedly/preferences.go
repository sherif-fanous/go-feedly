package feedly

import (
	"net/http"

	"github.com/dghubble/sling"
)

// PreferenceService provides methods for managing preferences.
type PreferenceService struct {
	sling *sling.Sling
}

// newPreferenceService returns a new PreferenceService.
func newPreferenceService(sling *sling.Sling) *PreferenceService {
	return &PreferenceService{
		sling: sling,
	}
}

// PreferenceListResponse represents the response from PreferenceService.List.
type PreferenceListResponse struct {
	Preferences map[string]string `json:"preferences"`
}

// List returns the application specific preferences.
func (s *PreferenceService) List() (*PreferenceListResponse, *http.Response, error) {
	decodedResponse := new(PreferenceListResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("preferences").Receive(&decodedResponse.Preferences, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// Update update the preferences of the user.
func (s *PreferenceService) Update(preferences map[string]string) (*http.Response, error) {
	apiError := new(APIError)

	resp, err := s.sling.New().Post("preferences").BodyJSON(preferences).Receive(nil, apiError)

	return resp, relevantError(err, apiError)
}
