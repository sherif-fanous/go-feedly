[![Build Status](https://travis-ci.com/sfanous/go-feedly.svg?branch=master)](https://travis-ci.com/sfanous/go-feedly)
[![Go Report Card](https://goreportcard.com/badge/github.com/sfanous/go-feedly)](https://goreportcard.com/report/github.com/sfanous/go-feedly)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/sfanous/go-feedly)
[![Release](https://img.shields.io/github/v/release/sfanous/go-feedly.svg?style=flat)](https://github.com/sfanous/go-feedly/releases/latest)

# Go client for the Feedly API

Currently, does not support any requests that require a Pro or Enterprise account.

The below example illustrates how to:

- Fetch a persisted OAuth2 token from a file in JSON format.
- Create a feedly.Client.
- Retrieve a list of collections.

The example assumes you have a Feedly OAuth2 token persisted in a file in JSON format such as the following file

```json
{
  "access_token": "<AccessToken>",
  "token_type": "Bearer",
  "refresh_token": "<RefreshToken>",
  "expiry": "2020-01-31T23:59:59.9999999-04:00"
}
```

You can easily obtain a Feedly developer token by following the instructions [here](https://developer.feedly.com/v3/developer/)

```go
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
```

[Full documentation is available on GoDoc.](https://godoc.org/github.com/sfanous/go-feedly/feedly)