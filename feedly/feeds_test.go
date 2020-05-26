package feedly_test

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"testing"

	"github.com/sfanous/go-feedly/feedly"
	"github.com/stretchr/testify/assert"
)

func testFeedServiceMetadata(t *testing.T) {
	metadataResponse, resp, err := client.Feeds.Metadata(*responseCollections[controlCollectionNames[rand.Intn(len(responseCollections))]].Feeds[0].ID)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, metadataResponse) {
		assert.IsType(t, &feedly.FeedMetadataResponse{}, metadataResponse)

		testUnmappedFields(t, metadataResponse, "FeedMetadataResponse")

		if doLog {
			b, err := json.MarshalIndent(metadataResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal metadataResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testFeedServiceMultipleMetadata(t *testing.T) {
	feeds := responseCollections[controlCollectionNames[rand.Intn(len(responseCollections))]].Feeds
	feedIDs := make([]string, 0, len(feeds))

	for i := 0; i < len(feeds); i++ {
		feedIDs = append(feedIDs, *feeds[i].ID)
	}

	multipleMetadataResponse, resp, err := client.Feeds.MultipleMetadata(feedIDs)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, multipleMetadataResponse) {
		assert.IsType(t, &feedly.FeedMultipleMetadataResponse{}, multipleMetadataResponse)
		assert.Equal(t, len(feedIDs), len(multipleMetadataResponse.Feeds))

		testUnmappedFields(t, multipleMetadataResponse, "FeedMultipleMetadataResponse")

		if doLog {
			b, err := json.MarshalIndent(multipleMetadataResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal multipleMetadataResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}
