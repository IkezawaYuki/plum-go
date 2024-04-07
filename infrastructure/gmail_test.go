package infrastructure

import (
	"fmt"
	"os"
	"testing"
)

func TestGmailService_CreateDraft(t *testing.T) {
	b, err := os.ReadFile("../credentials.json")
	if err != nil {
		t.Fatalf("%v", err)
	}
	gmailSrv := NewGmailService(b)
	fmt.Println(gmailSrv)
}
