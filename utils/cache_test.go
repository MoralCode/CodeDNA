package utils

import (
	"testing"
)

func TestCreation(t *testing.T) {

	cache := IdentityCache{
		Filename: "./cache.sqlite",
	}

	t.Error(cache.GetAll())
}

func TestHas(t *testing.T) {

	cache := IdentityCache{
		Filename: "./cache.sqlite",
	}
	cache.Add(IdentityValue{
		URL:       "https://example.com",
		Nickname:  "example",
		LineageID: "abcd1234",
	})

	if !cache.Has("https://example.com") {
		t.Errorf(`Cache doesnt Has() something that was just added`)
	}

	if !cache.Has("example.com") {
		t.Errorf(`Cache Has() fuzzy matching is not working`)
	}

}
