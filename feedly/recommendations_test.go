package feedly_test

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"testing"

	"github.com/sfanous/go-feedly/feedly"
	"github.com/stretchr/testify/assert"
)

func testRecommendationServiceTopic(t *testing.T) {
	topicResponse, resp, err := client.Recommendations.Topic(topics[rand.Intn(len(topics))], "en", &feedly.RecommendationTopicOptionalParams{
		Count: feedly.NewInt(100),
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, topicResponse) {
		assert.IsType(t, &feedly.RecommendationTopicResponse{}, topicResponse)

		testUnmappedFields(t, topicResponse, "RecommendationTopicResponse")

		if doLog {
			b, err := json.MarshalIndent(topicResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal topicResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}
