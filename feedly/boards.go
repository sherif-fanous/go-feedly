package feedly

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/dghubble/sling"
	"github.com/sfanous/go-feedly/internal/mapstructure"
	"github.com/sfanous/go-feedly/internal/mime"
	"github.com/sfanous/go-feedly/pkg/time"
)

// BoardService provides methods for managing personal boards, aka tags.
type BoardService struct {
	sling *sling.Sling
}

// newFeedService returns a new BoardService.
func newBoardService(sling *sling.Sling) *BoardService {
	return &BoardService{
		sling: sling,
	}
}

// Board is a Feedly board.
type Board struct {
	Cover          *string                `json:"cover,omitempty"`
	Created        *time.Time             `json:"created,omitempty"`
	Customizable   *bool                  `json:"customizable,omitempty"`
	Description    *string                `json:"description,omitempty"`
	Enterprise     *bool                  `json:"enterprise,omitempty"`
	HTMLURL        *string                `json:"htmlUrl,omitempty"`
	ID             *string                `json:"id,omitempty"`
	IsPublic       *bool                  `json:"isPublic,omitempty"`
	Label          *string                `json:"label,omitempty"`
	ShowHighlights *bool                  `json:"showHighlights,omitempty"`
	ShowNotes      *bool                  `json:"showNotes,omitempty"`
	StreamID       *string                `json:"streamId,omitempty"`
	UnmappedFields map[string]interface{} `json:"-" mapstructure:",remain"`
}

// AddEntry adds one entry to one or more existing boards.
func (s *BoardService) AddEntry(BoardIDs []string, entryID string) (*http.Response, error) {
	bodyJSON := &struct {
		EntryID string `json:"entryId,omitempty"`
	}{
		EntryID: entryID,
	}

	apiError := new(APIError)

	resp, err := s.sling.New().Put("tags/"+url.PathEscape(strings.Join(BoardIDs, ","))).BodyJSON(bodyJSON).Receive(nil, apiError)

	return resp, relevantError(err, apiError)
}

// AddEntries adds one or more entries to one or more existing boards.
func (s *BoardService) AddEntries(BoardIDs []string, entryIDs []string) (*http.Response, error) {
	bodyJSON := &struct {
		EntryIds []string `json:"entryIds,omitempty"`
	}{
		EntryIds: entryIDs,
	}

	apiError := new(APIError)

	resp, err := s.sling.New().Put("tags/"+url.PathEscape(strings.Join(BoardIDs, ","))).BodyJSON(bodyJSON).Receive(nil, apiError)

	return resp, relevantError(err, apiError)
}

// BoardCreateOptionalParams are the optional parameters for BoardService.Create.
type BoardCreateOptionalParams struct {
	Description    *string `json:"description,omitempty"`
	Feeds          []Feed  `json:"feeds,omitempty"`
	ID             *string `json:"id,omitempty"`
	IsPublic       *bool   `json:"isPublic,omitempty"`
	ShowHighlights *bool   `json:"showHighlights,omitempty"`
	ShowNotes      *bool   `json:"showNotes,omitempty"`
}

// BoardCreateResponse represents the response from BoardService.Create.
type BoardCreateResponse struct {
	Boards []Board `json:"boards"`
}

// Create creates a new board.
func (s *BoardService) Create(label string, optionalParams *BoardCreateOptionalParams) (*BoardCreateResponse, *http.Response, error) {
	if optionalParams == nil {
		optionalParams = &BoardCreateOptionalParams{}
	}

	bodyJSON := &struct {
		*BoardCreateOptionalParams
		Label string `json:"label"`
	}{
		BoardCreateOptionalParams: optionalParams,
		Label:                     label,
	}

	encodedResponse := make([]map[string]interface{}, 0)
	decodedResponse := new(BoardCreateResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Post("boards").BodyJSON(bodyJSON).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Boards); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// Delete deletes one or more existing boards.
func (s *BoardService) Delete(boardIDs []string) (*http.Response, error) {
	apiError := new(APIError)

	resp, err := s.sling.New().Delete("tags/"+url.PathEscape(strings.Join(boardIDs, ","))).Receive(nil, apiError)

	return resp, relevantError(err, apiError)
}

// BoardListOptionalParams are the optional parameters for BoardService.List.
type BoardListOptionalParams struct {
	WithEnterprise *bool `url:"withEnterprise,omitempty"`
}

// BoardListResponse represents the response from BoardService.List.
type BoardListResponse struct {
	Boards []Board `json:"boards"`
}

// List returns the list of boards.
func (s *BoardService) List(optionalParams *BoardListOptionalParams) (*BoardListResponse, *http.Response, error) {
	if optionalParams == nil {
		optionalParams = &BoardListOptionalParams{}
	}

	encodedResponse := make([]map[string]interface{}, 0)
	decodedResponse := new(BoardListResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("boards").QueryStruct(optionalParams).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Boards); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// RemoveEntry removes one entry from one or more existing boards.
func (s *BoardService) RemoveEntry(BoardIDs []string, entryID string) (*http.Response, error) {
	bodyJSON := &struct {
		EntryID string `json:"entryId,omitempty"`
	}{
		EntryID: entryID,
	}

	apiError := new(APIError)

	resp, err := s.sling.New().Delete("tags/"+url.PathEscape(strings.Join(BoardIDs, ","))).BodyJSON(bodyJSON).Receive(nil, apiError)

	return resp, relevantError(err, apiError)
}

// RemoveEntries removes one or more entries from one or more existing boards.
func (s *BoardService) RemoveEntries(BoardIDs []string, entryIDs []string) (*http.Response, error) {
	bodyJSON := &struct {
		EntryIds []string `json:"entryIds,omitempty"`
	}{
		EntryIds: entryIDs,
	}

	apiError := new(APIError)

	resp, err := s.sling.New().Delete("tags/"+url.PathEscape(strings.Join(BoardIDs, ","))).BodyJSON(bodyJSON).Receive(nil, apiError)

	return resp, relevantError(err, apiError)
}

// BoardUpdateOptionalParams are the optional parameters for BoardService.Update.
type BoardUpdateOptionalParams struct {
	DeleteCover    *bool   `json:"deleteCover,omitempty"`
	Description    *string `json:"description,omitempty"`
	IsPublic       *bool   `json:"isPublic,omitempty"`
	Label          *string `json:"label,omitempty"`
	ShowHighlights *bool   `json:"showHighlights,omitempty"`
	ShowNotes      *bool   `json:"showNotes,omitempty"`
}

// BoardUpdateResponse represents the response from BoardService.Update.
type BoardUpdateResponse struct {
	Boards []Board `json:"boards"`
}

// Update updates an existing board.
func (s *BoardService) Update(boardID string, optionalParams *BoardUpdateOptionalParams) (*BoardUpdateResponse, *http.Response, error) {
	if optionalParams == nil {
		optionalParams = &BoardUpdateOptionalParams{}
	}

	bodyJSON := &struct {
		*BoardUpdateOptionalParams
		ID string `json:"id"`
	}{
		BoardUpdateOptionalParams: optionalParams,
		ID:                        boardID,
	}

	encodedResponse := make([]map[string]interface{}, 0)
	decodedResponse := new(BoardUpdateResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Post("boards").BodyJSON(bodyJSON).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Boards); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// BoardUploadCoverImageResponse represents the response from BoardService.UploadCoverImage.
type BoardUploadCoverImageResponse struct {
	Boards []Board `json:"boards"`
}

// UploadCoverImage uploads a new cover image for an existing board.
func (s *BoardService) UploadCoverImage(boardID string, coverImage io.Reader) (*BoardUploadCoverImageResponse, *http.Response, error) {
	body, contentType, err := mime.CreateMultipartMIMEAttachment(coverImage)
	if err != nil {
		return nil, nil, err
	}

	encodedResponse := make([]map[string]interface{}, 0)
	decodedResponse := new(BoardUploadCoverImageResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Post("boards/"+url.PathEscape(boardID)).Body(body).Set("Content-Type", contentType).Receive(&encodedResponse, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	if err := mapstructure.Decode(encodedResponse, &decodedResponse.Boards); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}
