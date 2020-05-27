package feedly_test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sfanous/go-feedly/feedly"
	"golang.org/x/oauth2"
)

func ExampleNewClient_sandbox() {
	// Replace with a valid Feedly Sandbox OAuth2 Access Token
	oauth2AccessToken := ""
	oauth2Token := (*oauth2.Token)(nil)

	if err := json.Unmarshal([]byte(`{"access_token": "`+oauth2AccessToken+`","token_type": "Bearer"}`), &oauth2Token); err != nil {
		fmt.Printf("%v\n", err)

		return
	}

	client := feedly.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(oauth2Token)), feedly.WithAPIBaseURL("https://sandbox7.feedly.com"))

	profileListResponse, _, err := client.Profile.List()
	if err != nil {
		fmt.Printf("%v\n", err)

		return
	}

	fmt.Println(*profileListResponse.Profile.ID)
}
