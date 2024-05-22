package infrastructure

import (
	"fmt"
	"os"
	"testing"
)

func TestChatGPT_Create(t *testing.T) {
	chatgpt := NewChatGPT(
		os.Getenv("AZURE_OPENAI_KEY"),
		os.Getenv("AZURE_OPENAI_ENDPOINT"))
	generated, err := chatgpt.Create("お世話になっております。\n山田でございます。\n\n御社ではLPの制作はされていらっしゃいましたか？")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	fmt.Println(generated)
}

func TestChatGPT_Generate(t *testing.T) {
	chatgpt := NewChatGPT(
		os.Getenv("AZURE_OPENAI_KEY"),
		os.Getenv("AZURE_OPENAI_ENDPOINT"))
	generated, err := chatgpt.Generate("お世話になっております。\n山田でございます。\n\n御社ではLPの制作はされていらっしゃいましたか？")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	fmt.Println(generated)
}
