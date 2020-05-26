package feedly_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/sfanous/go-feedly/feedly"
	"github.com/stretchr/testify/assert"
)

func testStreamServiceContent(t *testing.T, stream interface{}) {
	streamID := ""
	streamType := ""

	switch stream := stream.(type) {
	case *feedly.Board:
		streamID = *stream.ID
		streamType = "board"
	case *feedly.Collection:
		streamID = *stream.ID
		streamType = "collection"
	case *feedly.Feed:
		streamID = *stream.ID
		streamType = "feed"
	case string:
		streamID = stream
		streamType = "system"
	}

	contentResponse, resp, err := client.Streams.Content(streamID, &feedly.StreamContentOptionalParams{
		Count:      feedly.NewInt(100),
		Ranked:     feedly.NewContentRank(feedly.Newest),
		UnreadOnly: feedly.NewBool(true),
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, contentResponse) {
		assert.IsType(t, &feedly.StreamContentResponse{}, contentResponse)

		responseStreams[streamType][*contentResponse.Stream.ID] = contentResponse.Stream
		testUnmappedFields(t, contentResponse, "StreamContentResponse")

		if doLog {
			b, err := json.MarshalIndent(contentResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal contentResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testStreamServiceEntryIDs(t *testing.T, stream interface{}) {
	streamID := ""

	switch stream := stream.(type) {
	case *feedly.Board:
		streamID = *stream.ID
	case *feedly.Collection:
		streamID = *stream.ID
	case *feedly.Feed:
		streamID = *stream.ID
	case string:
		streamID = stream
	}

	entryIDsResponse, resp, err := client.Streams.EntryIDs(streamID, &feedly.StreamEntryIDsOptionalParams{
		Count:      feedly.NewInt(100),
		Ranked:     feedly.NewContentRank(feedly.Newest),
		UnreadOnly: feedly.NewBool(true),
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, entryIDsResponse) {
		assert.IsType(t, &feedly.StreamEntryIDsResponse{}, entryIDsResponse)

		testUnmappedFields(t, entryIDsResponse, "StreamEntryIDsResponse")

		if doLog {
			b, err := json.MarshalIndent(entryIDsResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal entryIDsResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}
