package infrastructure

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("TestMain is invoked")
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println(os.Getenv("HUBSPOT_ACCESS_TOKEN"))
	code := m.Run()
	os.Exit(code)
}

func TestHubspot_GetContact(t *testing.T) {
	hubspot := NewHubspot(os.Getenv("HUBSPOT_ACCESS_TOKEN"))
	if err := hubspot.GetContact("6695090231"); err != nil {
		t.Fatalf("%s", err.Error())
	}
}
