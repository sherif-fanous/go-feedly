package feedly_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/sfanous/go-feedly/feedly"
	"github.com/stretchr/testify/assert"
)

func testProfileServiceList(t *testing.T) {
	listResponse, resp, err := client.Profile.List()
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, listResponse) {
		assert.IsType(t, &feedly.ProfileListResponse{}, listResponse)
		assert.Equal(t, "Feedly_"+now.Format(time.RFC3339), *listResponse.Profile.FamilyName)
		assert.Equal(t, "Go_"+now.Format(time.RFC3339), *listResponse.Profile.GivenName)

		profile = listResponse.Profile

		testUnmappedFields(t, listResponse, "ProfileListResponse")

		if doLog {
			b, err := json.MarshalIndent(listResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal listResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testProfileServiceUpdate(t *testing.T) {
	resp, err := client.Profile.Update(&feedly.Profile{
		FamilyName: feedly.NewString("Feedly_" + now.Format(time.RFC3339)),
		GivenName:  feedly.NewString("Go_" + now.Format(time.RFC3339)),
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}
