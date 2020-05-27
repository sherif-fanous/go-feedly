package feedly_test

import (
	"encoding/json"
	"github.com/sfanous/go-feedly/feedly"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
)

func testBoardServiceAddEntry(t *testing.T, boardIDs []string, entryID string) {
	resp, err := client.Boards.AddEntry(boardIDs, entryID)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func testBoardServiceAddMultipleEntries(t *testing.T, boardIDs []string, entryIDs []string) {
	resp, err := client.Boards.AddMultipleEntries(boardIDs, entryIDs)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func testBoardServiceCreate(t *testing.T, board *feedly.Board) {
	createResponse, resp, err := client.Boards.Create(strings.ToLower(*board.Label), &feedly.BoardCreateOptionalParams{
		Description: feedly.NewString(*board.Label + " created by go-feedly for testing"),
		IsPublic:    feedly.NewBool(false),
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, createResponse) {
		assert.IsType(t, &feedly.BoardCreateResponse{}, createResponse)
		assert.Equal(t, *board.Label+" created by go-feedly for testing", *createResponse.Boards[0].Description)
		assert.Equal(t, strings.ToLower(*board.Label), *createResponse.Boards[0].Label)

		responseBoards[strings.ToLower(*createResponse.Boards[0].Label)] = &createResponse.Boards[0]

		testUnmappedFields(t, createResponse, "BoardCreateResponse")

		if doLog {
			b, err := json.MarshalIndent(createResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal createResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testBoardServiceDelete(t *testing.T, board *feedly.Board) {
	resp, err := client.Boards.Delete([]string{
		*board.ID,
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	delete(responseBoards, strings.ToLower(*board.Label))
}

func testBoardServiceDeleteEntry(t *testing.T, boardIDs []string, entryID string) {
	resp, err := client.Boards.DeleteEntry(boardIDs, entryID)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func testBoardServiceDeleteMultipleEntries(t *testing.T, boardIDs []string, entryIDs []string) {
	resp, err := client.Boards.DeleteMultipleEntries(boardIDs, entryIDs)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func testBoardsServiceDetails(t *testing.T, board *feedly.Board) {
	detailResponse, resp, err := client.Boards.Details(*board.ID)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, detailResponse) {
		assert.IsType(t, &feedly.BoardDetailResponse{}, detailResponse)

		responseBoards[strings.ToLower(*detailResponse.Boards[0].Label)] = &detailResponse.Boards[0]

		testUnmappedFields(t, detailResponse, "BoardDetailResponse")

		if doLog && detailResponse != nil {
			b, err := json.MarshalIndent(detailResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal detailResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testBoardServiceList(t *testing.T) {
	listResponse, resp, err := client.Boards.List(&feedly.BoardListOptionalParams{
		WithEnterprise: feedly.NewBool(false),
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, listResponse) {
		assert.IsType(t, &feedly.BoardListResponse{}, listResponse)

		for _, board := range listResponse.Boards {
			board := board

			if _, ok := responseBoards[strings.ToLower(*board.Label)]; ok {
				responseBoards[strings.ToLower(*board.Label)] = &board
			}
		}

		testUnmappedFields(t, listResponse, "BoardListResponse")

		if doLog {
			b, err := json.MarshalIndent(listResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal listResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testBoardServiceUpdate(t *testing.T, board *feedly.Board) {
	updateResponse, resp, err := client.Boards.Update(*board.ID, &feedly.BoardUpdateOptionalParams{
		DeleteCover: feedly.NewBool(true),
		Description: feedly.NewString(*board.Label + " updated by go-feedly for testing"),
		Label:       feedly.NewString(strings.ToUpper(*board.Label)),
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, updateResponse) {
		assert.IsType(t, &feedly.BoardUpdateResponse{}, updateResponse)
		assert.Nil(t, updateResponse.Boards[0].Cover)
		assert.Equal(t, *board.Label+" updated by go-feedly for testing", *updateResponse.Boards[0].Description)
		assert.Equal(t, strings.ToUpper(*board.Label), *updateResponse.Boards[0].Label)

		responseBoards[strings.ToLower(*updateResponse.Boards[0].Label)] = &updateResponse.Boards[0]

		testUnmappedFields(t, updateResponse, "BoardUpdateResponse")

		if doLog {
			b, err := json.MarshalIndent(updateResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal updateResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}

func testBoardServiceUploadCoverImage(t *testing.T, board *feedly.Board) {
	coverImage, err := os.OpenFile(coverImages[rand.Intn(len(coverImages))], os.O_RDONLY, 0444)
	if err != nil {
		t.Errorf("%v", err)
	}
	defer func() {
		if err := coverImage.Close(); err != nil {
			t.Logf("%v", err)
		}
	}()

	uploadCoverImageResponse, resp, err := client.Boards.UploadCoverImage(*board.ID, coverImage)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if assert.NotNil(t, uploadCoverImageResponse) {
		assert.IsType(t, &feedly.BoardUploadCoverImageResponse{}, uploadCoverImageResponse)
		assert.NotNil(t, uploadCoverImageResponse.Boards[0].Cover)

		responseBoards[strings.ToLower(*uploadCoverImageResponse.Boards[0].Label)] = &uploadCoverImageResponse.Boards[0]

		testUnmappedFields(t, uploadCoverImageResponse, "BoardUploadCoverImageResponse")

		if doLog {
			b, err := json.MarshalIndent(uploadCoverImageResponse, "", "    ")
			if err != nil {
				t.Logf("Failed to marshal uploadCoverImageResponse: %v", err)
			}

			t.Log(string(b))
		}
	}
}
