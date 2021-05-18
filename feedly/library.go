package feedly

import (
	"net/http"
	"net/url"

	"github.com/dghubble/sling"
	"github.com/sfanous/go-feedly/internal/mapstructure"
)

// LibraryService provides methods for managing a user's public library.
type LibraryService struct {
	sling *sling.Sling
}

// newFeedService returns a new LibraryService.
func newLibraryService(sling *sling.Sling) *LibraryService {
	return &LibraryService{
		sling: sling,
	}
}

// Cover is a Feedly library cover.
type Cover struct {
	About           *string                `json:"about,omitempty"`
	Alias           *string                `json:"alias,omitempty"`
	BackgroundImage *string                `json:"backgroundImage,omitempty"`
	FullName        *string                `json:"fullName,omitempty"`
	LinkedIn        *string                `json:"linkedIn,omitempty"`
	Picture         *string                `json:"picture,omitempty"`
	Twitter         *string                `json:"twitter,omitempty"`
	UnmappedFields  map[string]interface{} `json:"-" mapstructure:",remain"`
}

// Library is a Feedly library.
type Library struct {
	Collections    []Collection           `json:"collections,omitempty"`
	Cover          *Cover                 `json:"cover,omitempty"`
	Tags           []Board                `json:"tags,omitempty"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// LibraryAliasAvailableResponse represents the response from LibraryService.AliasAvailable.
type LibraryAliasAvailableResponse struct {
	Available      *bool                  `json:"available,omitempty"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// AliasAvailable checks if an alias is available to be used.
func (s *LibraryService) AliasAvailable(alias string) (*LibraryAliasAvailableResponse, *http.Response, error) {
	encodedResponse := make(map[string]interface{})
	decodedResponse := new(LibraryAliasAvailableResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("alias/"+url.PathEscape(alias)).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// LibraryCoverResponse represents the response from LibraryService.Cover.
type LibraryCoverResponse struct {
	Cover          *Cover                 `json:"cover"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// Cover returns the library cover.
func (s *LibraryService) Cover() (*LibraryCoverResponse, *http.Response, error) {
	encodedResponse := make(map[string]interface{})
	decodedResponse := new(LibraryCoverResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("library/cover").Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Cover); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// Delete deletes a library.
func (s *LibraryService) Delete() (*http.Response, error) {
	apiError := new(APIError)

	resp, err := s.sling.New().Delete("library/cover").Receive(nil, apiError)

	return resp, relevantError(err, apiError)
}

// LibraryDetailsResponse represents the response from LibraryService.Details.
type LibraryDetailsResponse struct {
	Library        *Library               `json:"library"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// Details returns the library details.
func (s *LibraryService) Details(alias string) (*LibraryDetailsResponse, *http.Response, error) {
	encodedResponse := make(map[string]interface{})
	decodedResponse := new(LibraryDetailsResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("library/"+alias).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Library); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// LibraryLeoIndustriesResponse represents the response from LibraryService.LeoIndustries.
type LibraryLeoIndustriesResponse struct {
	Cover          *Cover                 `json:"cover"`
	Collections    []Collection           `json:"collections"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// LeoIndustries returns the Leo industries.
func (s *LibraryService) LeoIndustries() (*LibraryLeoIndustriesResponse, *http.Response, error) {
	encodedResponse := make(map[string]interface{})
	decodedResponse := new(LibraryLeoIndustriesResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("library/leoIndustries").Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, decodedResponse); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// LibraryListSharedResourcesResponse represents the response from LibraryService.ListSharedResources.
type LibraryListSharedResourcesResponse struct {
	SharedResources map[string]struct {
		Scope          *string                `json:"scope,omitempty"`
		Target         *string                `json:"target,omitempty"`
		UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
	} `json:"sharedResources,omitempty"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// ListSharedResources returns the list of shared resources.
func (s *LibraryService) ListSharedResources() (*LibraryListSharedResourcesResponse, *http.Response, error) {
	encodedResponse := make(map[string]interface{})
	decodedResponse := new(LibraryListSharedResourcesResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("library/acl").Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.SharedResources); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// ShareResource shares a resource.
func (s *LibraryService) ShareResource(collectionID string) (*http.Response, error) {
	bodyJSON := &struct {
		Scope string `json:"scope,omitempty"`
	}{
		Scope: "view",
	}

	apiError := new(APIError)

	resp, err := s.sling.New().Get("library/acl/"+url.PathEscape(collectionID)+"/"+url.PathEscape("global.public")).BodyJSON(bodyJSON).Receive(nil, apiError)

	return resp, relevantError(err, apiError)
}

// UnshareResource unshares a resource.
func (s *LibraryService) UnshareResource(collectionID string) (*http.Response, error) {
	apiError := new(APIError)

	resp, err := s.sling.New().Delete("library/acl/"+url.PathEscape(collectionID)+"/"+url.PathEscape("global.public")).Receive(nil, apiError)

	return resp, relevantError(err, apiError)
}

// LibraryUpdateCoverResponse represents the response from LibraryService.UpdateCover.
type LibraryUpdateCoverResponse struct {
	Cover          *Cover                 `json:"cover"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// UpdateCover updates the library cover.
func (s *LibraryService) UpdateCover(cover *Cover) (*LibraryUpdateCoverResponse, *http.Response, error) {
	encodedResponse := make(map[string]interface{})
	decodedResponse := new(LibraryUpdateCoverResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Post("library/cover").BodyJSON(cover).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Cover); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}
