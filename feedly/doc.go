/*
Package feedly provides a Go client for the Feedly API.


For more information about the Feedly API, see the documentation:
https://developer.feedly.com/

Authentication

By design, the feedly Client accepts any http.Client so OAuth2 requests can be
made by using the appropriate authenticated client.
Use the https://github.com/golang/oauth2 package to obtain an http.Client which
transparently authorizes requests.

Usage

You use the library by creating a Client and invoking its methods. The client
can be created manually with NewClient.

The below example illustrates how to:

- Fetch a persisted OAuth2 token from a file in JSON format.

- Create a feedly.Client.

- Retrieve a list of collections.

The example assumes you have a Feedly OAuth2 token persisted in a file in JSON format such as the following file.

	{
	  "access_token": "<AccessToken>",
	  "token_type": "Bearer",
	  "refresh_token": "<RefreshToken>",
	  "expiry": "2020-01-31T23:59:59.9999999-04:00"
	}

You can easily obtain a Feedly developer token by following the instructions: https://developer.feedly.com/v3/developer/

	package main

	import (
		"context"
		"encoding/json"
		"flag"
		"fmt"
		"io/ioutil"

		"github.com/sfanous/go-feedly/feedly"
		"golang.org/x/oauth2"
	)

	var filename string

	func fetchToken() (*oauth2.Token, error) {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		oauth2Token := oauth2.Token{}

		if err := json.Unmarshal(b, &oauth2Token); err != nil {
			return nil, err
		}

		return &oauth2Token, nil
	}

	func main() {
		flag.StringVar(&filename, "file", "", "token persistent store path")
		flag.Parse()

		oauth2Token, err := fetchToken()
		if err != nil {
			fmt.Printf("Failed to fetch OAuth2 token: %v", err)

			return
		}

		f := feedly.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(oauth2Token)))

		listResponse, _, err := f.Collections.List(nil)
		if err != nil {
			fmt.Printf("Failed to list Collections: %v", err)

			return
		}

		b, err := json.MarshalIndent(listResponse, "", "    ")
		if err != nil {
			fmt.Printf("Failed to marshal listResponse: %v", err)
		}

		fmt.Println(string(b))
	}


*/
package feedly
