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
	gmailSrv := NewGmailService(b, "../token.json")
	fmt.Println(gmailSrv)
	if err := gmailSrv.CreateDraft("ご担当者様", "yuki.ikezawa@strategy-drive.jp"); err != nil {
		t.Fatalf("%s", err.Error())
	}
}
