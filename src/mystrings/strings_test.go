package mystrings

import (
	"testing"
	"strings"
)

func TestStrings(t *testing.T) {

	target := "search 日本語"

	if !strings.HasPrefix(target ,"search ") {
		t.Error("Error")
	}
	switch {
	case strings.HasPrefix(target ,"search ") :
		t.Log("入ったよ")
	}

	word := target[7:]
	if word == "日本語" {
		t.Log(word)
	}

	word = "テスト"
	t.Log(target)

}
