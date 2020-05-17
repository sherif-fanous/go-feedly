package feedly

import (
	"encoding/xml"
	"io"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/sfanous/go-feedly/pkg/decoders"
)

// OPMLService provides methods for managing OPML files.
type OPMLService struct {
	sling *sling.Sling
}

// newOPMLService returns a new OPMLService.
func newOPMLService(sling *sling.Sling) *OPMLService {
	return &OPMLService{
		sling: sling,
	}
}

// OPML describes the OPML format for storing outlines in XML 1.0 called Outline Processor Markup Language.
type OPML struct {
	XMLName *xml.Name `xml:"opml,omitempty"`
	Version *string   `xml:"version,attr,omitempty"`
	Head    *Head     `xml:"head,omitempty"`
	Body    *Body     `xml:"body,omitempty"`
}

// Head describes the head element in the opml element of OPML.
type Head struct {
	Title           *string `xml:"title,omitempty"`
	DateCreated     *string `xml:"dateCreated,omitempty"`
	DateModified    *string `xml:"dateModified,omitempty"`
	OwnerName       *string `xml:"ownerName,omitempty"`
	OwnerEmail      *string `xml:"ownerEmail,omitempty"`
	OwnerID         *string `xml:"ownerId,omitempty"`
	Docs            *string `xml:"docs,omitempty"`
	ExpansionState  *string `xml:"expansionState,omitempty"`
	VertScrollState *string `xml:"vertScrollState,omitempty"`
	WindowTop       *string `xml:"windowTop,omitempty"`
	WindowBottom    *string `xml:"windowBottom,omitempty"`
	WindowLeft      *string `xml:"windowLeft,omitempty"`
	WindowRight     *string `xml:"windowRight,omitempty"`
}

// Body describes the body element in the opml element of OPML.
type Body struct {
	Outlines []Outline `xml:"outline"`
}

// Outline describes the outline element in the body element of OPML.
type Outline struct {
	Outlines     []Outline `xml:"outline,omitempty"`
	Text         *string   `xml:"text,attr,omitempty"`
	Type         *string   `xml:"type,attr,omitempty"`
	IsComment    *string   `xml:"isComment,attr,omitempty"`
	IsBreakpoint *string   `xml:"isBreakpoint,attr,omitempty"`
	Created      *string   `xml:"created,attr,omitempty"`
	Category     *string   `xml:"category,attr,omitempty"`
	XMLURL       *string   `xml:"xmlUrl,attr,omitempty"`
	HTMLURL      *string   `xml:"htmlUrl,attr,omitempty"`
	URL          *string   `xml:"url,attr,omitempty"`
	Language     *string   `xml:"language,attr,omitempty"`
	Title        *string   `xml:"title,attr,omitempty"`
	Version      *string   `xml:"version,attr,omitempty"`
	Description  *string   `xml:"description,attr,omitempty"`
}

// OPMLExportResponse represents the response from OPMLService.Export.
type OPMLExportResponse struct {
	OPML *OPML `json:"opml"`
}

// Export exports the user’s subscriptions.
func (s *OPMLService) Export() (*OPMLExportResponse, *http.Response, error) {
	decodedResponse := new(OPMLExportResponse)
	apiError := new(APIError)

	resp, err := s.sling.New().Get("opml").ResponseDecoder(decoders.XMLDecoder{}).Receive(&decodedResponse.OPML, apiError)
	if err := relevantError(err, apiError); err != nil {
		return nil, resp, err
	}

	return decodedResponse, resp, nil
}

// Import imports the user's subscriptions.
func (s *OPMLService) Import(opml io.Reader) (*http.Response, error) {
	apiError := new(APIError)

	resp, err := s.sling.New().Post("opml").Body(opml).Set("Content-Type", "text/xml").Receive(nil, apiError)

	return resp, relevantError(err, apiError)
}
