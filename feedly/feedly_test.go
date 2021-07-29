package feedly_test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/sfanous/go-feedly/feedly"
	"golang.org/x/oauth2"
)

type stringSliceFlags []string

func (f *stringSliceFlags) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *stringSliceFlags) Set(flagValue string) error {
	*f = append(*f, flagValue)

	return nil
}

var apiBaseURL string
var apiBaseVersion string
var client *feedly.Client
var collectionsToMark []*feedly.Collection
var controlBoards = make(map[string]*feedly.Board)
var controlCollections = make(map[string]*feedly.Collection)
var controlBoardNames = make([]string, 0)
var controlCollectionNames = make([]string, 0)
var controlCollectionNamesAddedFeeds = make(map[string][]feedly.Feed)
var coverImages stringSliceFlags
var doLog bool
var entriesToMark []feedly.Entry
var feedsToMark []feedly.Feed
var librarySharedCollection string
var oauth2TokenFile string
var oauth2Token *oauth2.Token
var opml *feedly.OPML
var now time.Time
var numberOfCollections int
var numberOfFeeds int
var profile *feedly.Profile
var responseBoards = make(map[string]*feedly.Board)
var responseCollections = make(map[string]*feedly.Collection)
var responseStreams = make(map[string]map[string]*feedly.Stream)
var topics = []string{}

func prepareTestData() {
	leoIndustriesResponse, _, err := client.Library.LeoIndustries()
	if err != nil {
		fmt.Printf("feedly_test: prepareTestData: %v\n", err)

		os.Exit(1)
	}

	rand.Shuffle(len(leoIndustriesResponse.Collections), func(i int, j int) {
		leoIndustriesResponse.Collections[i], leoIndustriesResponse.Collections[j] = leoIndustriesResponse.Collections[j], leoIndustriesResponse.Collections[i]
	})

	uniqueTopics := make(map[string]struct{})
	numberOfCollectionsAdded := 0

	for i := 0; i < len(leoIndustriesResponse.Collections); i++ {
		collection := leoIndustriesResponse.Collections[i]

		if len(collection.Feeds) >= numberOfFeeds && numberOfCollectionsAdded < numberOfCollections {
			numberOfCollectionsAdded++

			rand.Shuffle(len(collection.Feeds), func(i int, j int) {
				collection.Feeds[i], collection.Feeds[j] = collection.Feeds[j], collection.Feeds[i]
			})

			controlBoards[strings.ToLower(*collection.Label)] = &feedly.Board{
				Label: collection.Label,
			}

			controlCollections[strings.ToLower(*collection.Label)] = &feedly.Collection{
				Label: collection.Label,
				Feeds: collection.Feeds,
			}

			controlBoardNames = append(controlBoardNames, strings.ToLower(*collection.Label))
			controlCollectionNames = append(controlCollectionNames, strings.ToLower(*collection.Label))
			controlCollectionNamesAddedFeeds[strings.ToLower(*collection.Label)] = make([]feedly.Feed, 0)
		}

		for _, topic := range collection.Topics {
			uniqueTopics[topic] = struct{}{}
		}
	}

	for topic := range uniqueTopics {
		topics = append(topics, topic)
	}

	responseStreams["board"] = make(map[string]*feedly.Stream)
	responseStreams["collection"] = make(map[string]*feedly.Stream)
	responseStreams["feed"] = make(map[string]*feedly.Stream)
	responseStreams["system"] = make(map[string]*feedly.Stream)
}

func sleep() {
	time.Sleep((time.Second * time.Duration(rand.Intn(3))) + (time.Millisecond * time.Duration(rand.Float64()*1000.0)))
}

func testUnmappedFields(t *testing.T, v interface{}, namePrefix string) {
	interfaceValue := reflect.ValueOf(v)
	interfaceType := interfaceValue.Type()

	// interfaceValue is a slice.
	if interfaceValue.Kind() == reflect.Slice {
		for i := 0; i < interfaceValue.Len(); i++ {
			testUnmappedFields(t, interfaceValue.Index(i).Interface(), namePrefix+"["+strconv.Itoa(i)+"]")
		}
	}

	// interfaceValue is a pointer.
	if interfaceValue.Kind() == reflect.Ptr {
		if interfaceValue.IsNil() {
			return
		} else {
			// Dereference the non-nil pointer.
			interfaceValue = interfaceValue.Elem()
			interfaceType = interfaceValue.Type()

			// interfaceValue is not a struct.
			if interfaceValue.Kind() != reflect.Struct {
				return
			}
		}
	}

	// interfaceValue is neither a pointer nor a struct.
	if interfaceValue.Kind() != reflect.Struct {
		return
	}

	// Iterate over the fields of the struct.
	for i := 0; i < interfaceValue.NumField(); i++ {
		fieldValue := interfaceValue.Field(i)
		fieldType := interfaceType.Field(i)

		// fieldValue is a pointer.
		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				continue
			} else {
				// Dereference the non-nil pointer.
				fieldValue = fieldValue.Elem()

				// Skip unexported fields.
				if !fieldValue.CanInterface() {
					continue
				}

				if fieldValue.Kind() == reflect.Struct {
					testUnmappedFields(t, fieldValue.Interface(), namePrefix+"."+fieldType.Name)

					continue
				} else {
					continue
				}
			}
		}

		// Skip unexported fields
		if !fieldValue.CanInterface() {
			continue
		}

		if fieldValue.Kind() == reflect.Map && fieldType.Name == "UnmappedFields" {
			mapIter := fieldValue.MapRange()

			for mapIter.Next() {
				keyValue := mapIter.Key()
				valueValue := mapIter.Value()

				t.Logf("Missing field in struct %v\nField name: %q\nField type: %v\nField value: %#v", namePrefix, keyValue.Interface(), reflect.TypeOf(valueValue.Interface()).Kind(), valueValue.Interface())
			}
		}

		if fieldValue.Kind() == reflect.Struct {
			testUnmappedFields(t, fieldValue.Interface(), namePrefix+"."+fieldType.Name)

			continue
		}

		if fieldValue.Kind() == reflect.Slice {
			for i := 0; i < fieldValue.Len(); i++ {
				testUnmappedFields(t, fieldValue.Index(i).Interface(), namePrefix+"."+fieldType.Name+"["+strconv.Itoa(i)+"]")
			}

			continue
		}
	}
}

func TestFeedly(t *testing.T) {
	t.Run("PreferenceServiceUpdateSet", func(t *testing.T) {
		testPreferenceServiceUpdateSet(t)
	})
	sleep()

	t.Run("PreferenceServiceList", func(t *testing.T) {
		testPreferenceServiceList(t)
	})
	sleep()

	t.Run("PreferenceServiceUpdateDelete", func(t *testing.T) {
		testPreferenceServiceUpdateDelete(t)
	})
	sleep()

	t.Run("ProfileServiceUpdate", func(t *testing.T) {
		testProfileServiceUpdate(t)
	})
	sleep()

	t.Run("ProfileServiceList", func(t *testing.T) {
		testProfileServiceList(t)
	})
	sleep()

	t.Run("RecommendationServiceTopic", func(t *testing.T) {
		testRecommendationServiceTopic(t)
	})
	sleep()

	t.Run("SearchServiceFeeds", func(t *testing.T) {
		testSearchServiceFeeds(t)
	})
	sleep()

	t.Run("SearchServiceStream", func(t *testing.T) {
		testSearchServiceStream(t)
	})
	sleep()

	for collectionLabel, collection := range controlCollections {
		t.Run(fmt.Sprintf("CollectionServiceCreate %s", strings.Title(collectionLabel)), func(t *testing.T) {
			testCollectionServiceCreate(t, collection)
		})

		sleep()
	}

	t.Run("CollectionServiceList", func(t *testing.T) {
		testCollectionServiceList(t)
	})
	sleep()

	for _, collectionLabel := range controlCollectionNames {
		collectionLabel := collectionLabel

		t.Run(fmt.Sprintf("CollectionServiceUploadCoverImage %s", strings.Title(collectionLabel)), func(t *testing.T) {
			testCollectionServiceUploadCoverImage(t, responseCollections[collectionLabel])
		})
		sleep()

		t.Run(fmt.Sprintf("CollectionServiceDetails %s", strings.Title(collectionLabel)), func(t *testing.T) {
			testCollectionsServiceDetails(t, responseCollections[collectionLabel])
		})
		sleep()

		t.Run(fmt.Sprintf("CollectionServiceUpdate %s", strings.Title(collectionLabel)), func(t *testing.T) {
			testCollectionServiceUpdate(t, responseCollections[collectionLabel])
		})
		sleep()

		t.Run(fmt.Sprintf("CollectionServiceAddFeed %s", strings.Title(collectionLabel)), func(t *testing.T) {
			testCollectionServiceAddFeed(t, responseCollections[collectionLabel])
		})
		sleep()

		t.Run(fmt.Sprintf("CollectionServiceAddMultipleFeeds %s", strings.Title(collectionLabel)), func(t *testing.T) {
			testCollectionServiceAddMultipleFeeds(t, responseCollections[collectionLabel])
		})
		sleep()

		t.Run(fmt.Sprintf("CollectionServiceDeleteFeed %s", strings.Title(collectionLabel)), func(t *testing.T) {
			testCollectionServiceDeleteFeed(t, responseCollections[collectionLabel])
		})
		sleep()

		t.Run(fmt.Sprintf("CollectionServiceDeleteMultipleFeeds %s", strings.Title(collectionLabel)), func(t *testing.T) {
			testCollectionServiceDeleteMultipleFeeds(t, responseCollections[collectionLabel])
		})
		sleep()

		t.Run(fmt.Sprintf("StreamServiceContent_Collection %s", strings.Title(collectionLabel)), func(t *testing.T) {
			testStreamServiceContent(t, responseCollections[collectionLabel])
		})
		sleep()

		t.Run(fmt.Sprintf("StreamServiceEntryIDs %s", strings.Title(collectionLabel)), func(t *testing.T) {
			testStreamServiceEntryIDs(t, responseCollections[collectionLabel])
		})
		sleep()
	}

	for boardLabel, board := range controlBoards {
		t.Run(fmt.Sprintf("BoardServiceCreate %s", strings.Title(boardLabel)), func(t *testing.T) {
			testBoardServiceCreate(t, board)
		})

		sleep()
	}

	t.Run("BoardServiceList", func(t *testing.T) {
		testBoardServiceList(t)
	})
	sleep()

	for _, boardLabel := range controlBoardNames {
		collection := responseCollections[boardLabel]
		collectionStreamedEntries := responseStreams["collection"][*collection.ID].Items

		board := responseBoards[boardLabel]
		t.Run(fmt.Sprintf("BoardServiceUploadCoverImage %s", strings.Title(boardLabel)), func(t *testing.T) {
			testBoardServiceUploadCoverImage(t, board)
		})
		sleep()

		board = responseBoards[boardLabel]
		t.Run(fmt.Sprintf("BoardServiceDetails %s", strings.Title(boardLabel)), func(t *testing.T) {
			testBoardsServiceDetails(t, board)
		})
		sleep()

		board = responseBoards[boardLabel]
		t.Run(fmt.Sprintf("BoardServiceUpdate %s", strings.Title(boardLabel)), func(t *testing.T) {
			testBoardServiceUpdate(t, board)
		})
		sleep()

		entryIDs := make([]string, 0)
		randomizedEntryIndexes := rand.Perm(len(collectionStreamedEntries))

		for i := 0; i < rand.Intn(len(collectionStreamedEntries)-(len(collectionStreamedEntries)/2))+(len(collectionStreamedEntries)/2); i++ {
			entryIDs = append(entryIDs, *collectionStreamedEntries[randomizedEntryIndexes[i]].ID)
		}

		board = responseBoards[boardLabel]
		t.Run(fmt.Sprintf("BoardServiceAddEntry %s", strings.Title(boardLabel)), func(t *testing.T) {
			testBoardServiceAddEntry(t, []string{*board.ID}, entryIDs[0])
		})
		sleep()

		board = responseBoards[boardLabel]
		t.Run(fmt.Sprintf("BoardServiceAddMultipleEntries %s", strings.Title(boardLabel)), func(t *testing.T) {
			testBoardServiceAddMultipleEntries(t, []string{*board.ID}, entryIDs[1:])
		})
		sleep()

		board = responseBoards[boardLabel]
		t.Run(fmt.Sprintf("StreamServiceContent_Board %s", strings.Title(boardLabel)), func(t *testing.T) {
			testStreamServiceContent(t, board)
		})
		sleep()

		board = responseBoards[boardLabel]
		t.Run(fmt.Sprintf("StreamServiceEntryIDs %s", strings.Title(boardLabel)), func(t *testing.T) {
			testStreamServiceEntryIDs(t, board)
		})
		sleep()

		board = responseBoards[boardLabel]
		t.Run(fmt.Sprintf("BoardServiceDeleteEntry %s", strings.Title(boardLabel)), func(t *testing.T) {
			testBoardServiceDeleteEntry(t, []string{*board.ID}, entryIDs[0])
		})
		sleep()

		board = responseBoards[boardLabel]
		t.Run(fmt.Sprintf("BoardServiceDeleteMultipleEntries %s", strings.Title(boardLabel)), func(t *testing.T) {
			testBoardServiceDeleteMultipleEntries(t, []string{*board.ID}, entryIDs[1:])
		})
		sleep()
	}

	t.Run("EntryServiceContent", func(t *testing.T) {
		testEntryServiceContent(t)
	})
	sleep()

	t.Run("EntryServiceMultipleContent", func(t *testing.T) {
		testEntryServiceMultipleContent(t)
	})
	sleep()

	t.Run("FeedServiceMetadata", func(t *testing.T) {
		testFeedServiceMetadata(t)
	})
	sleep()

	t.Run("FeedServiceMultipleMetadata", func(t *testing.T) {
		testFeedServiceMultipleMetadata(t)
	})
	sleep()

	t.Run("LibraryServiceAliasAvailable", func(t *testing.T) {
		testLibraryServiceAliasAvailable(t)
	})
	sleep()

	t.Run("LibraryServiceUpdateCover", func(t *testing.T) {
		testLibraryServiceUpdateCover(t)
	})
	sleep()

	t.Run("LibraryServiceCover", func(t *testing.T) {
		testLibraryServiceCover(t)
	})
	sleep()

	t.Run("LibraryServiceShareResource", func(t *testing.T) {
		testLibraryServiceShareResource(t)
	})
	sleep()

	t.Run("LibraryServiceListSharedResources", func(t *testing.T) {
		testLibraryServiceListSharedResources(t)
	})
	sleep()

	t.Run("LibraryServiceUnshareResource", func(t *testing.T) {
		testLibraryServiceUnshareResource(t)
	})
	sleep()

	t.Run("LibraryServiceListSharedResources", func(t *testing.T) {
		testLibraryServiceListSharedResources(t)
	})
	sleep()

	t.Run("LibraryServiceDetails", func(t *testing.T) {
		testLibraryServiceDetails(t)
	})
	sleep()

	t.Run("LibraryServiceDelete", func(t *testing.T) {
		testLibraryServiceDelete(t)
	})
	sleep()

	t.Run("MixServiceMostEngagingInCollection", func(t *testing.T) {
		testMixServiceMostEngagingInCollection(t)
	})
	sleep()

	t.Run("MixServiceMostEngagingInFeed", func(t *testing.T) {
		testMixServiceMostEngagingInFeed(t)
	})
	sleep()

	t.Run("MixServiceMostEngagingInSystemCategory", func(t *testing.T) {
		testMixServiceMostEngagingInSystemCategory(t)
	})
	sleep()

	t.Run("MixServiceMostEngagingInTopic", func(t *testing.T) {
		testMixServiceMostEngagingInTopic(t)
	})
	sleep()

	numberOfCollectionsToMark := rand.Intn(len(controlCollectionNames)-(len(controlCollectionNames)/2)) + (len(controlCollectionNames) / 2)
	collectionsToMark = make([]*feedly.Collection, 0, numberOfCollectionsToMark)
	randomizedCollectionNamesIndexes := rand.Perm(len(controlCollectionNames))

	for i := 0; i < numberOfCollectionsToMark; i++ {
		collectionsToMark = append(collectionsToMark, responseCollections[controlCollectionNames[randomizedCollectionNamesIndexes[i]]])
	}

	numberOfFeedsToMark := rand.Intn(len(controlCollectionNamesAddedFeeds[strings.ToLower(*collectionsToMark[0].Label)])-(len(controlCollectionNamesAddedFeeds[strings.ToLower(*collectionsToMark[0].Label)])/2)) + (len(controlCollectionNamesAddedFeeds[strings.ToLower(*collectionsToMark[0].Label)]) / 2)
	feedsToMark = controlCollectionNamesAddedFeeds[strings.ToLower(*collectionsToMark[0].Label)][0:numberOfFeedsToMark]

	numberOfEntriesToMark := rand.Intn(len(responseStreams["collection"][*responseCollections[controlCollectionNames[randomizedCollectionNamesIndexes[0]]].ID].Items)-(len(responseStreams["collection"][*responseCollections[controlCollectionNames[randomizedCollectionNamesIndexes[0]]].ID].Items)/2)) + (len(responseStreams["collection"][*responseCollections[controlCollectionNames[randomizedCollectionNamesIndexes[0]]].ID].Items) / 2)
	entriesToMark = responseStreams["collection"][*responseCollections[controlCollectionNames[randomizedCollectionNamesIndexes[0]]].ID].Items[0:numberOfEntriesToMark]

	t.Run("MarkerServiceUnreadCounts", func(t *testing.T) {
		testMarkerServiceUnreadCounts(t)
	})
	sleep()

	t.Run("MarkerServiceMarkEntriesAsSaved", func(t *testing.T) {
		testMarkerServiceMarkEntriesAsSaved(t)
	})
	sleep()

	t.Run("MarkerServiceLatestTagged", func(t *testing.T) {
		testMarkerServiceLatestTagged(t)
	})
	sleep()

	t.Run("MarkerServiceMarkEntriesAsUnsaved", func(t *testing.T) {
		testMarkerServiceMarkEntriesAsUnsaved(t)
	})
	sleep()

	t.Run("MarkerServiceMarkTagsAsRead", func(t *testing.T) {
		testMarkerServiceMarkTagsAsRead(t)
	})
	sleep()

	t.Run("MarkerServiceMarkEntriesAsRead", func(t *testing.T) {
		testMarkerServiceMarkEntriesAsRead(t)
	})
	sleep()

	t.Run("MarkerServiceLatestRead", func(t *testing.T) {
		testMarkerServiceLatestRead(t)
	})
	sleep()

	t.Run("MarkerServiceMarkEntriesAsUnread", func(t *testing.T) {
		testMarkerServiceMarkEntriesAsUnread(t)
	})
	sleep()

	t.Run("MarkerServiceMarkFeedsAsRead", func(t *testing.T) {
		testMarkerServiceMarkFeedsAsRead(t)
	})
	sleep()

	t.Run("MarkerServiceMarkFeedsAsReadUndo", func(t *testing.T) {
		testMarkerServiceMarkFeedsAsReadUndo(t)
	})
	sleep()

	t.Run("MarkerServiceMarkCollectionsAsRead", func(t *testing.T) {
		testMarkerServiceMarkCollectionsAsRead(t)
	})
	sleep()

	t.Run("MarkerServiceMarkCollectionsAsReadUndo", func(t *testing.T) {
		testMarkerServiceMarkCollectionsAsReadUndo(t)
	})
	sleep()

	t.Run("OPMLServiceExport", func(t *testing.T) {
		testOPMLServiceExport(t)
	})
	sleep()

	for collectionLabel, collection := range responseCollections {
		t.Run(fmt.Sprintf("CollectionServiceDelete %s", strings.Title(collectionLabel)), func(t *testing.T) {
			testCollectionServiceDelete(t, collection)
		})

		sleep()
	}

	for boardLabel, board := range responseBoards {
		t.Run(fmt.Sprintf("BoardServiceDelete %s", strings.Title(boardLabel)), func(t *testing.T) {
			testBoardServiceDelete(t, board)
		})

		sleep()
	}

	t.Run("OPMLServiceImport", func(t *testing.T) {
		testOPMLServiceImport(t)
	})
	sleep()

	t.Run("CollectionServiceList", func(t *testing.T) {
		testCollectionServiceList(t)
	})
	sleep()

	for collectionLabel, collection := range responseCollections {
		t.Run(fmt.Sprintf("CollectionServiceDelete %s", strings.Title(collectionLabel)), func(t *testing.T) {
			testCollectionServiceDelete(t, collection)
		})

		sleep()
	}
}

func TestMain(m *testing.M) {
	now = time.Now()
	rand.Seed(now.UnixNano())

	flag.StringVar(&apiBaseURL, "api_url", feedly.APIBaseURL, fmt.Sprintf("The feedly API base URL. Default: %s", feedly.APIBaseURL))
	flag.StringVar(&apiBaseVersion, "api_version", feedly.APIBaseVersion, fmt.Sprintf("The feedly API base version. Default: %s", feedly.APIBaseVersion))
	flag.Var(&coverImages, "cover", "Path to a GIF/JPG/PNG file. Flag can be repeated for multiple files")
	flag.BoolVar(&doLog, "log", false, "Enable logging while executing tests. Default: Disabled")
	flag.StringVar(&oauth2TokenFile, "oauth2", "", "Absolute or relative path to a file persisting an OAuth2 token in JSON format")
	flag.IntVar(&numberOfCollections, "collections", 3, "Number of Collections to create. Default: 3")
	flag.IntVar(&numberOfFeeds, "feeds", 8, "Number of Feeds to subscribe to per Collection. Default: 8")
	flag.Parse()

	var b []byte

	if oauth2AccessToken, ok := os.LookupEnv("FEEDLY_ACCESS_TOKEN"); ok {
		b = []byte(`{"access_token": "` + oauth2AccessToken + `","token_type": "Bearer"}`)
	} else {
		var err error

		b, err = ioutil.ReadFile(oauth2TokenFile)
		if err != nil {
			fmt.Printf("feedly_test: TestMain: %v\n", err)

			os.Exit(1)
		}
	}

	if err := json.Unmarshal(b, &oauth2Token); err != nil {
		fmt.Printf("feedly_test: TestMain: %v\n", err)

		os.Exit(1)
	}

	client = feedly.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(oauth2Token)), feedly.WithAPIBaseURL(apiBaseURL), feedly.WithAPIBaseVersion(apiBaseVersion))

	prepareTestData()

	os.Exit(m.Run())
}
