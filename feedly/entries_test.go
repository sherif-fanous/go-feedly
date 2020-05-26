package feedly_test

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"testing"

	"github.com/sfanous/go-feedly/feedly"
	"github.com/stretchr/testify/assert"
)

func testEntryServiceContent(t *testing.T) {
	contentResponse, resp, err := client.Entries.Content(*responseStreams["collection"][*responseCollections[controlCollectionNames[rand.Intn(len(controlCollectionNames))]].ID].Items[0].ID)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, contentResponse) {
		assert.IsType(t, &feedly.EntryContentResponse{}, contentResponse)

		testUnmappedFields(t, contentResponse, "EntryContentResponse")

		if doLog {
			b, err := json.MarshalIndent(contentResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal contentResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testEntryServiceMultipleContent(t *testing.T) {
	entryIDs := make([]string, 0)
	entries := responseStreams["collection"][*responseCollections[controlCollectionNames[rand.Intn(len(controlCollectionNames))]].ID].Items

	for i := 0; i < rand.Intn(len(entries)-5)+5; i++ {
		entryIDs = append(entryIDs, *entries[i].ID)
	}

	multipleContentResponse, resp, err := client.Entries.MultipleContent(entryIDs)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, multipleContentResponse) {
		assert.IsType(t, &feedly.EntryMultipleContentResponse{}, multipleContentResponse)

		testUnmappedFields(t, multipleContentResponse, "EntryMultipleContentResponse")

		if doLog {
			b, err := json.MarshalIndent(multipleContentResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal multipleContentResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}
