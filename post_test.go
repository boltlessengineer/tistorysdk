package tistorysdk_test

import (
	"encoding/json"
	"testing"

	"github.com/boltlessengineer/tistorysdk"
)

func TestList(t *testing.T) {
	t.Run("Parsing Post", func(t *testing.T) {
		const (
			ExampleJSON string = `
			{
				"id": "1",
				"title": "Some Title",
				"content": "Some Content",
				"categoryId": "0",
				"postUrl": "http://oath.tistory.com/1",
				"visibility": "0",
				"acceptComment":"1",
				"acceptTrackback":"1",
				"tags": {
					"tag": ["open", "api"]
				},
				"comments": "0",
				"trackbacks": "0",
				"date": "1303352668"
			}
			`
		)
		var post tistorysdk.Post
		if err := json.Unmarshal([]byte(ExampleJSON), &post); err != nil {
			t.Fatal(err)
		}
		//t.Log(post)
	})
}
