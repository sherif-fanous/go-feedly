package feedly_test

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"testing"

	"github.com/sfanous/go-feedly/feedly"
	"github.com/stretchr/testify/assert"
)

func testMarkerServiceLatestRead(t *testing.T) {
	latestReadResponse, resp, err := client.Markers.LatestRead(nil)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, latestReadResponse) {
		assert.IsType(t, &feedly.MarkerLatestReadResponse{}, latestReadResponse)

		testUnmappedFields(t, latestReadResponse, "MarkerLatestReadResponse")

		if doLog {
			b, err := json.MarshalIndent(latestReadResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal latestReadResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testMarkerServiceLatestTagged(t *testing.T) {
	latestTaggedResponse, resp, err := client.Markers.LatestTagged(nil)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, latestTaggedResponse) {
		assert.IsType(t, &feedly.MarkerLatestTaggedResponse{}, latestTaggedResponse)

		testUnmappedFields(t, latestTaggedResponse, "MarkerLatestTaggedResponse")

		if doLog {
			b, err := json.MarshalIndent(latestTaggedResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal latestTaggedResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testMarkerServiceMarkCollectionsAsRead(t *testing.T) {
	collectionIDs := make([]string, 0, len(collectionsToMark))

	for i := 0; i < len(collectionsToMark); i++ {
		collectionIDs = append(collectionIDs, *collectionsToMark[i].ID)
	}

	resp, err := client.Markers.Mark(feedly.MarkAsRead, feedly.Collections, &feedly.MarkerMarkOptionalParams{
		CollectionIDs: collectionIDs,
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func testMarkerServiceMarkCollectionsAsReadUndo(t *testing.T) {
	collectionIDs := make([]string, 0, len(collectionsToMark))

	for i := 0; i < len(collectionsToMark); i++ {
		collectionIDs = append(collectionIDs, *collectionsToMark[i].ID)
	}

	resp, err := client.Markers.Mark(feedly.UndoMarkAsRead, feedly.Collections, &feedly.MarkerMarkOptionalParams{
		CollectionIDs: collectionIDs,
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func testMarkerServiceMarkEntriesAsRead(t *testing.T) {
	entryIDs := make([]string, 0, len(entriesToMark))

	for i := 0; i < len(entriesToMark); i++ {
		entryIDs = append(entryIDs, *entriesToMark[i].ID)
	}

	resp, err := client.Markers.Mark(feedly.MarkAsRead, feedly.Entries, &feedly.MarkerMarkOptionalParams{
		EntryIDs: entryIDs,
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func testMarkerServiceMarkEntriesAsSaved(t *testing.T) {
	entryIDs := make([]string, 0, len(entriesToMark))

	for i := 0; i < len(entriesToMark); i++ {
		entryIDs = append(entryIDs, *entriesToMark[i].ID)
	}

	resp, err := client.Markers.Mark(feedly.MarkAsSaved, feedly.Entries, &feedly.MarkerMarkOptionalParams{
		EntryIDs: entryIDs,
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func testMarkerServiceMarkEntriesAsUnread(t *testing.T) {
	entryIDs := make([]string, 0, len(entriesToMark))

	for i := 0; i < len(entriesToMark); i++ {
		entryIDs = append(entryIDs, *entriesToMark[i].ID)
	}

	resp, err := client.Markers.Mark(feedly.KeepUnread, feedly.Entries, &feedly.MarkerMarkOptionalParams{
		EntryIDs: entryIDs,
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func testMarkerServiceMarkEntriesAsUnsaved(t *testing.T) {
	entryIDs := make([]string, 0, len(entriesToMark))

	for i := 0; i < len(entriesToMark); i++ {
		entryIDs = append(entryIDs, *entriesToMark[i].ID)
	}

	resp, err := client.Markers.Mark(feedly.MarkAsUnsaved, feedly.Entries, &feedly.MarkerMarkOptionalParams{
		EntryIDs: entryIDs,
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func testMarkerServiceMarkFeedsAsRead(t *testing.T) {
	feedIDs := make([]string, 0, len(feedsToMark))

	for i := 0; i < len(feedsToMark); i++ {
		feedIDs = append(feedIDs, *feedsToMark[i].ID)
	}

	resp, err := client.Markers.Mark(feedly.MarkAsRead, feedly.Feeds, &feedly.MarkerMarkOptionalParams{
		FeedIDs: feedIDs,
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func testMarkerServiceMarkFeedsAsReadUndo(t *testing.T) {
	feedIDs := make([]string, 0, len(feedsToMark))

	for i := 0; i < len(feedsToMark); i++ {
		feedIDs = append(feedIDs, *feedsToMark[i].ID)
	}

	resp, err := client.Markers.Mark(feedly.UndoMarkAsRead, feedly.Feeds, &feedly.MarkerMarkOptionalParams{
		FeedIDs: feedIDs,
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func testMarkerServiceMarkTagsAsRead(t *testing.T) {
	resp, err := client.Markers.Mark(feedly.MarkAsRead, feedly.Tags, &feedly.MarkerMarkOptionalParams{
		TagIDs: []string{
			*responseBoards[controlBoardNames[rand.Intn(len(responseBoards))]].ID,
		},
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func testMarkerServiceUnreadCounts(t *testing.T) {
	unreadCountsResponse, resp, err := client.Markers.UnreadCounts(nil)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, unreadCountsResponse) {
		assert.IsType(t, &feedly.MarkerUnreadCountsResponse{}, unreadCountsResponse)

		testUnmappedFields(t, unreadCountsResponse, "MarkerUnreadCountsResponse")

		if doLog {
			b, err := json.MarshalIndent(unreadCountsResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal unreadCountsResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}
