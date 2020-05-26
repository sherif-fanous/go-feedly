package feedly_test

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"testing"

	"github.com/sfanous/go-feedly/feedly"
	"github.com/stretchr/testify/assert"
)

func testSearchServiceFeeds(t *testing.T) {
	feedsResponse, resp, err := client.Search.Feeds(topics[rand.Intn(len(topics))], &feedly.SearchFeedsOptionalParams{
		Count:  feedly.NewInt(10),
		Locale: feedly.NewString("en"),
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, feedsResponse) {
		assert.IsType(t, &feedly.SearchFeedsResponse{}, feedsResponse)

		testUnmappedFields(t, feedsResponse, "SearchFeedsResponse")

		if doLog {
			b, err := json.MarshalIndent(feedsResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal feedsResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testSearchServiceStream(t *testing.T) {
	streamResponse, resp, err := client.Search.Stream("topic/global.popular", topics[rand.Intn(len(topics))], &feedly.SearchStreamOptionalParams{
		Count:      feedly.NewInt(10),
		Engagement: feedly.NewEngagementFilter(feedly.High),
		Fields: &feedly.FieldFilter{
			Author: true,
			Title:  true,
		},
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, streamResponse) {
		assert.IsType(t, &feedly.SearchStreamResponse{}, streamResponse)

		testUnmappedFields(t, streamResponse, "SearchStreamResponse")

		if doLog {
			b, err := json.MarshalIndent(streamResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal streamResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}
