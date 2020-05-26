package feedly_test

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/sfanous/go-feedly/feedly"
	"github.com/stretchr/testify/assert"
)

func testCollectionServiceAddFeed(t *testing.T, collection *feedly.Collection) {
	controlCollection := controlCollections[strings.ToLower(*collection.Label)]

	addFeedResponse, resp, err := client.Collections.AddFeed(*collection.ID, &feedly.Feed{
		ID: controlCollection.Feeds[0].ID,
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, addFeedResponse) {
		expectedNumberOfFeeds := len(collection.Feeds) + 1

		assert.IsType(t, &feedly.CollectionAddFeedResponse{}, addFeedResponse)
		assert.Equal(t, expectedNumberOfFeeds, len(addFeedResponse.Feeds))

		responseCollections[strings.ToLower(*collection.Label)].Feeds = addFeedResponse.Feeds

		controlCollection.Feeds = append(controlCollection.Feeds[1:], controlCollection.Feeds[0])

		testUnmappedFields(t, addFeedResponse, "CollectionAddFeedResponse")

		if doLog && addFeedResponse != nil {
			b, err := json.MarshalIndent(addFeedResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal addFeedResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testCollectionServiceAddMultipleFeeds(t *testing.T, collection *feedly.Collection) {
	controlCollection := controlCollections[strings.ToLower(*collection.Label)]

	numberOfFeedsToAdd := rand.Intn(numberOfFeeds) + 1

	addMultipleFeedsResponse, resp, err := client.Collections.AddMultipleFeeds(*collection.ID, controlCollection.Feeds[:numberOfFeedsToAdd])
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, addMultipleFeedsResponse) {
		expectedNumberOfFeeds := len(collection.Feeds) + numberOfFeedsToAdd

		assert.IsType(t, &feedly.CollectionAddMultipleFeedsResponse{}, addMultipleFeedsResponse)
		assert.Equal(t, expectedNumberOfFeeds, len(addMultipleFeedsResponse.Feeds))

		responseCollections[strings.ToLower(*collection.Label)].Feeds = addMultipleFeedsResponse.Feeds

		testUnmappedFields(t, addMultipleFeedsResponse, "CollectionAddMultipleFeedsResponse")

		if doLog && addMultipleFeedsResponse != nil {
			b, err := json.MarshalIndent(addMultipleFeedsResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal addMultipleFeedsResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testCollectionServiceCreate(t *testing.T, collection *feedly.Collection) {
	controlCollection := controlCollections[strings.ToLower(*collection.Label)]

	createResponse, resp, err := client.Collections.Create(strings.ToLower(*collection.Label), &feedly.CollectionCreateOptionalParams{
		Description: feedly.NewString(*collection.Label + " created by go-feedly for testing"),
		Feeds:       collection.Feeds[:numberOfFeeds],
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, createResponse) {
		assert.IsType(t, &feedly.CollectionCreateResponse{}, createResponse)
		assert.Equal(t, *collection.Label+" created by go-feedly for testing", *createResponse.Collections[0].Description)
		assert.Equal(t, strings.ToLower(*collection.Label), *createResponse.Collections[0].Label)
		assert.Equal(t, numberOfFeeds, len(createResponse.Collections[0].Feeds))

		controlCollection.Feeds = append(controlCollection.Feeds[numberOfFeeds:], controlCollection.Feeds[:numberOfFeeds]...)
		responseCollections[strings.ToLower(*createResponse.Collections[0].Label)] = &createResponse.Collections[0]

		testUnmappedFields(t, createResponse, "CollectionCreateResponse")

		if doLog && createResponse != nil {
			b, err := json.MarshalIndent(createResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal createResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testCollectionServiceDelete(t *testing.T, collection *feedly.Collection) {
	resp, err := client.Collections.Delete(*collection.ID)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	delete(responseCollections, strings.ToLower(*collection.Label))
}

func testCollectionServiceDeleteFeed(t *testing.T, collection *feedly.Collection) {
	resp, err := client.Collections.DeleteFeed(*collection.ID, *collection.Feeds[rand.Intn(len(collection.Feeds))].ID, nil)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testCollectionServiceDeleteMultipleFeeds(t *testing.T, collection *feedly.Collection) {
	numberOfFeedsToDelete := rand.Intn(numberOfFeeds) + 1

	feedsToDelete := make([]string, 0, numberOfFeedsToDelete)

	for i := 0; i < numberOfFeedsToDelete; i++ {
		feedsToDelete = append(feedsToDelete, *collection.Feeds[i].ID)
	}

	resp, err := client.Collections.DeleteMultipleFeeds(*collection.ID, feedsToDelete, nil)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testCollectionsServiceDetails(t *testing.T, collection *feedly.Collection) {
	detailResponse, resp, err := client.Collections.Details(*collection.ID)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, detailResponse) {
		assert.IsType(t, &feedly.CollectionDetailResponse{}, detailResponse)

		responseCollections[strings.ToLower(*detailResponse.Collections[0].Label)] = &detailResponse.Collections[0]

		testUnmappedFields(t, detailResponse, "CollectionDetailResponse")

		if doLog && detailResponse != nil {
			b, err := json.MarshalIndent(detailResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal detailResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testCollectionServiceList(t *testing.T) {
	listResponse, resp, err := client.Collections.List(&feedly.CollectionListOptionalParams{
		WithEnterprise: feedly.NewBool(false),
		WithStats:      feedly.NewBool(true),
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, listResponse) {
		assert.IsType(t, &feedly.CollectionListResponse{}, listResponse)
		assert.Equal(t, len(controlCollections), len(listResponse.Collections))

		for _, collection := range listResponse.Collections {
			collection := collection

			responseCollections[strings.ToLower(*collection.Label)] = &collection
		}

		testUnmappedFields(t, listResponse, "CollectionListResponse")

		if doLog && listResponse != nil {
			b, err := json.MarshalIndent(listResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal listResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testCollectionServiceUpdate(t *testing.T, collection *feedly.Collection) {
	controlCollection := controlCollections[strings.ToLower(*collection.Label)]

	updateResponse, resp, err := client.Collections.Update(*collection.ID, &feedly.CollectionUpdateOptionalParams{
		DeleteCover: feedly.NewBool(true),
		Description: feedly.NewString(*collection.Label + " updated by go-feedly for testing"),
		Label:       feedly.NewString(strings.ToUpper(*collection.Label)),
		Feeds:       controlCollection.Feeds[:numberOfFeeds],
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, updateResponse) {
		assert.IsType(t, &feedly.CollectionUpdateResponse{}, updateResponse)
		assert.Nil(t, updateResponse.Collections[0].Cover)
		assert.Equal(t, *collection.Label+" updated by go-feedly for testing", *updateResponse.Collections[0].Description)
		assert.Equal(t, strings.ToUpper(*collection.Label), *updateResponse.Collections[0].Label)
		assert.Equal(t, numberOfFeeds*2, len(updateResponse.Collections[0].Feeds))

		controlCollection.Feeds = append(controlCollection.Feeds[numberOfFeeds:], controlCollection.Feeds[:numberOfFeeds]...)
		responseCollections[strings.ToLower(*updateResponse.Collections[0].Label)] = &updateResponse.Collections[0]

		testUnmappedFields(t, updateResponse, "CollectionUpdateResponse")

		if doLog {
			b, err := json.MarshalIndent(updateResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal updateResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testCollectionServiceUploadCoverImage(t *testing.T, collection *feedly.Collection) {
	coverImage, err := os.OpenFile(coverImages[rand.Intn(len(coverImages))], os.O_RDONLY, 0444)
	if err != nil {
		t.Errorf("%v", err)
	}
	defer func() {
		if err := coverImage.Close(); err != nil {
			t.Logf("%v", err)
		}
	}()

	uploadCoverImageResponse, resp, err := client.Collections.UploadCoverImage(*collection.ID, coverImage)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, uploadCoverImageResponse) {
		assert.IsType(t, &feedly.CollectionUploadCoverImageResponse{}, uploadCoverImageResponse)
		assert.NotNil(t, uploadCoverImageResponse.Collections[0].Cover)

		responseCollections[strings.ToLower(*uploadCoverImageResponse.Collections[0].Label)] = &uploadCoverImageResponse.Collections[0]

		testUnmappedFields(t, uploadCoverImageResponse, "CollectionUploadCoverImageResponse")

		if doLog && uploadCoverImageResponse != nil {
			b, err := json.MarshalIndent(uploadCoverImageResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal uploadCoverImageResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}
