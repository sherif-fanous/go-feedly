package feedly_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/sfanous/go-feedly/feedly"
	"github.com/stretchr/testify/assert"
)

func testPreferenceServiceList(t *testing.T) {
	listResponse, resp, err := client.Preferences.List()
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, listResponse) {
		assert.IsType(t, &feedly.PreferenceListResponse{}, listResponse)
		assert.Equal(t, now.Format(time.RFC3339), listResponse.Preferences["currentDateTime"])

		if doLog {
			b, err := json.MarshalIndent(listResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal listResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testPreferenceServiceUpdate(t *testing.T, preferences map[string]string) {
	resp, err := client.Preferences.Update(preferences)
	if err != nil {
		t.Errorf("feedly_test: TestPreferenceService_Update: %v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func testPreferenceServiceUpdateDelete(t *testing.T) {
	testPreferenceServiceUpdate(t, map[string]string{
		"currentDateTime": "==DELETE==",
	})
}

func testPreferenceServiceUpdateSet(t *testing.T) {
	testPreferenceServiceUpdate(t, map[string]string{
		"currentDateTime": now.Format(time.RFC3339),
	})
}
