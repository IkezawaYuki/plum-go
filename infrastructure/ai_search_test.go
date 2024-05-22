package infrastructure

import (
	"fmt"
	"os"
	"testing"
)

func TestAISearch_SearchDocuments(t *testing.T) {
	aiSearchSearch := NewAISearch(
		os.Getenv("AI_SEARCH_URL"),
		os.Getenv("AI_SEARCH_API_KEY2"))
	result, err := aiSearchSearch.SearchDocuments("LP")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(result)
}
