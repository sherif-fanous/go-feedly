package feedly_test

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"testing"

	"github.com/sfanous/go-feedly/feedly"
	"github.com/stretchr/testify/assert"
)

func testLibraryServiceAliasAvailable(t *testing.T) {
	aliasAvailableResponse, resp, err := client.Library.AliasAvailable("yldeef_og")
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, aliasAvailableResponse) {
		assert.IsType(t, &feedly.LibraryAliasAvailableResponse{}, aliasAvailableResponse)

		testUnmappedFields(t, aliasAvailableResponse, "LibraryAliasAvailableResponse")

		if doLog {
			b, err := json.MarshalIndent(aliasAvailableResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal aliasAvailableResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testLibraryServiceCover(t *testing.T) {
	coverResponse, resp, err := client.Library.Cover()
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, coverResponse) {
		assert.IsType(t, &feedly.LibraryCoverResponse{}, coverResponse)

		testUnmappedFields(t, coverResponse, "LibraryCoverResponse")

		if doLog {
			b, err := json.MarshalIndent(coverResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal coverResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testLibraryServiceDelete(t *testing.T) {
	resp, err := client.Library.Delete()
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func testLibraryServiceDetails(t *testing.T) {
	detailsResponse, resp, err := client.Library.Details("yldeef_og")
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, detailsResponse) {
		assert.IsType(t, &feedly.LibraryDetailsResponse{}, detailsResponse)

		testUnmappedFields(t, detailsResponse, "LibraryDetailsResponse")

		if doLog {
			b, err := json.MarshalIndent(detailsResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal detailsResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testLibraryServiceListSharedResources(t *testing.T) {
	listSharedResourcesResponse, resp, err := client.Library.ListSharedResources()
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, listSharedResourcesResponse) {
		assert.IsType(t, &feedly.LibraryListSharedResourcesResponse{}, listSharedResourcesResponse)

		testUnmappedFields(t, listSharedResourcesResponse, "LibraryListSharedResourcesResponse")

		if doLog {
			b, err := json.MarshalIndent(listSharedResourcesResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal listSharedResourcesResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testLibraryServiceShareResource(t *testing.T) {
	librarySharedCollection = *responseCollections[controlCollectionNames[rand.Intn(len(controlCollectionNames))]].ID

	resp, err := client.Library.ShareResource(librarySharedCollection)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func testLibraryServiceUnshareResource(t *testing.T) {
	resp, err := client.Library.UnshareResource(librarySharedCollection)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func testLibraryServiceUpdateCover(t *testing.T) {
	updateCoverResponse, resp, err := client.Library.UpdateCover(&feedly.Cover{
		About:    feedly.NewString("Cover updated by go-feedly for testing"),
		Alias:    feedly.NewString("yldeef_og"),
		FullName: feedly.NewString(*profile.FullName),
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, updateCoverResponse) {
		assert.IsType(t, &feedly.LibraryUpdateCoverResponse{}, updateCoverResponse)

		testUnmappedFields(t, updateCoverResponse, "LibraryUpdateCoverResponse")

		if doLog {
			b, err := json.MarshalIndent(updateCoverResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal updateCoverResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}
