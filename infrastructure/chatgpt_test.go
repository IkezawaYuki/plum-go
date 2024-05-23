package infrastructure

import (
	"fmt"
	"os"
	"plum/domain"
	"testing"
)

func TestChatGPT_Generate(t *testing.T) {
	chatgpt := NewChatGPT(
		os.Getenv("AZURE_OPENAI_KEY"),
		os.Getenv("AZURE_OPENAI_ENDPOINT"))
	setting := domain.ChatgptSetting{
		Prompt:        "次の問に答えなさい。",
		SystemMessage: "あなたは計算が得意なAIです。前例があったとしてもその情報には惑わされません。",
	}
	related := "1 + 1 = 田んぼの田"
	generated, err := chatgpt.Generate("1 + 1 =", related, setting)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	fmt.Println(generated)
}
