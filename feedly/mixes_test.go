package feedly_test

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"testing"

	"github.com/sfanous/go-feedly/feedly"
	"github.com/stretchr/testify/assert"
)

func testMixServiceMostEngaging(t *testing.T, streamID string) {
	mostEngagingResponse, resp, err := client.Mixes.MostEngaging(streamID, &feedly.MixMostEngagingOptionalParams{
		Backfill:   feedly.NewBool(true),
		Count:      feedly.NewInt(10),
		Hours:      feedly.NewInt(24),
		Locale:     feedly.NewString("en"),
		UnreadOnly: feedly.NewBool(true),
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, mostEngagingResponse) {
		assert.IsType(t, &feedly.MixMostEngagingResponse{}, mostEngagingResponse)

		testUnmappedFields(t, mostEngagingResponse, "MixMostEngagingResponse")

		if doLog {
			b, err := json.MarshalIndent(mostEngagingResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal mostEngagingResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testMixServiceMostEngagingInCollection(t *testing.T) {
	testMixServiceMostEngaging(t, *responseCollections[controlCollectionNames[rand.Intn(len(responseCollections))]].ID)
}

func testMixServiceMostEngagingInFeed(t *testing.T) {
	testMixServiceMostEngaging(t, *responseCollections[controlCollectionNames[rand.Intn(len(responseCollections))]].Feeds[0].ID)
}

func testMixServiceMostEngagingInSystemCategory(t *testing.T) {
	testMixServiceMostEngaging(t, "user/"+*profile.ID+"/category/global.all")
}

func testMixServiceMostEngagingInTopic(t *testing.T) {
	testMixServiceMostEngaging(t, "topic/"+topics[rand.Intn(len(topics))])
}
