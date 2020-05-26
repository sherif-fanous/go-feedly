package feedly_test

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/sfanous/go-feedly/feedly"
	"github.com/stretchr/testify/assert"
)

func testOPMLServiceExport(t *testing.T) {
	opmlExportResponse, resp, err := client.OPML.Export()
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, opmlExportResponse) {
		assert.IsType(t, &feedly.OPMLExportResponse{}, opmlExportResponse)

		opml = opmlExportResponse.OPML

		if doLog {
			b, err := xml.MarshalIndent(opmlExportResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal opmlExportResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testOPMLServiceImport(t *testing.T) {
	b, err := xml.MarshalIndent(opml, "", "    ")
	if err != nil {
		t.Logf("Failed to marshal opml: %v", err)
	}

	resp, err := client.OPML.Import(bytes.NewBuffer(b))
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}
