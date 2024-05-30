package helper_test

import (
	"fmt"
	"testing"
	"time"
	helper "youtube-scraper/helper"

	"github.com/stretchr/testify/assert"
)

const DEFAULT_ASSERTED_TEXT = "Expected value: %s - Received Value: %s"

func TestTrimWhiteSpaces(t *testing.T) {
	var emptyText string

	tests := []struct {
		Value   string
		Expects string
	}{
		{Value: emptyText, Expects: ""},
		{Value: "", Expects: ""},
		{Value: "Teste should take    empty spaces", Expects: "Teste should take empty spaces"},
		{Value: " Take initial and end spaces ", Expects: "Take initial and end spaces"},
		{Value: "     ", Expects: ""},
	}

	for _, test := range tests {
		result := helper.TrimWhiteSpaces(test.Value)

		assert.Equal(
			t,
			test.Expects,
			result,
			DEFAULT_ASSERTED_TEXT,
			test.Expects,
			result,
		)
	}
}

func TestTrimNonWords(t *testing.T) {
	var emptyText string

	tests := []struct {
		Value   string
		Expects string
	}{
		{Value: emptyText, Expects: ""},
		{Value: "", Expects: ""},
		{Value: " Take of Non word characters like ':' &*{}@3#zx.<./;\\¬4¹²³££¢¬´´``^~;", Expects: "Take of Non word characters like zx"},
		{Value: " ./\\()\"':,.;<>~!@#$%^&*|+=[]{}`~?- ", Expects: ""},
		{Value: "Remove😄Emojis💩💩💩💩test  🤖 random 👁👄👁💅🙌🧋🪬🪬🪬word ", Expects: "Remove Emojis test random word"},
	}

	for _, test := range tests {
		result := helper.TrimNonWords(test.Value)

		assert.Equal(
			t,
			test.Expects,
			result,
			DEFAULT_ASSERTED_TEXT,
			test.Expects,
			result,
		)
	}
}

func TestTrimLinks(t *testing.T) {
	var emptyText string

	tests := []struct {
		Value   string
		Expects string
	}{
		{Value: emptyText, Expects: ""},
		{Value: "", Expects: ""},
		{Value: " Take off Links: http://wwww.google.com.br ", Expects: "Take off Links:"},
		{Value: " linkonmiddlehttps://wwww.google.com.broftext ", Expects: "linkonmiddle"},
		{Value: " Shorted URL https://github.com/xeubaka random text  ", Expects: "Shorted URL random text"},
	}

	for _, test := range tests {
		result := helper.TrimLinks(test.Value)

		assert.Equal(
			t,
			test.Expects,
			result,
			DEFAULT_ASSERTED_TEXT,
			test.Expects,
			result,
		)
	}
}

func TestTrimEndLine(t *testing.T) {
	var emptyText string

	tests := []struct {
		Value   string
		Expects string
	}{
		{Value: emptyText, Expects: ""},
		{Value: "\r\n", Expects: ""},
		{Value: " Take off \n", Expects: " Take off "},
		{Value: " NoEndString\r", Expects: " NoEndString"},
	}

	for _, test := range tests {
		result := helper.TrimEndLine(test.Value)

		assert.Equal(
			t,
			test.Expects,
			result,
			DEFAULT_ASSERTED_TEXT,
			test.Expects,
			result,
		)
	}
}

func TestConvertStringToInt(t *testing.T) {
	var emptyText string

	tests := []struct {
		Value    string
		Expects  int
		printErr string
	}{
		{Value: emptyText, Expects: 0, printErr: "|  is not an number. |\n"},
		{Value: "1", Expects: 1, printErr: ""},
		{Value: "1@123", Expects: 0, printErr: "| 1@123 is not an number. |\n"},
	}

	for _, test := range tests {
		result, printErr := helper.ConvertStringToInt(test.Value)
		assert.Equal(
			t,
			test.Expects,
			result,
			DEFAULT_ASSERTED_TEXT,
			test.Expects,
			result,
		)
		assert.Equal(
			t,
			test.printErr,
			printErr,
			DEFAULT_ASSERTED_TEXT,
			test.printErr,
			printErr,
		)
	}
}

func TestConvertStringToTime(t *testing.T) {
	var emptyText string

	tests := []struct {
		Value   string
		Expects time.Duration
	}{
		{Value: emptyText, Expects: time.Duration(0)},
		{Value: "15M02S", Expects: time.Duration(902 * time.Second)},
		{Value: "PT13M", Expects: time.Duration(13 * time.Minute)},
		{Value: "P4DT", Expects: time.Duration(4 * 24 * time.Hour)},
	}

	recoverPanic := func() (recovered string) {
		if r := recover(); r != nil {
			fmt.Println("recovered")
		}
		return
	}

	for _, test := range tests {
		defer recoverPanic()
		result := helper.ConvertStringToTime(test.Value)
		assert.Equal(
			t,
			test.Expects,
			result,
			DEFAULT_ASSERTED_TEXT,
			test.Expects,
			result,
		)
	}
}

func TestConvertTimeDurationToString(t *testing.T) {
	tests := []struct {
		Value   time.Duration
		Expects string
	}{
		{Value: time.Duration(0), Expects: "0s"},
		{Value: time.Duration(902 * time.Second), Expects: "15m2s"},
		{Value: time.Duration(13 * time.Minute), Expects: "13m"},
		{Value: time.Duration(4 * 24 * time.Hour), Expects: "96h"},
	}

	recoverPanic := func() (recovered string) {
		if r := recover(); r != nil {
			fmt.Println("recovered")
		}
		return
	}

	for _, test := range tests {
		defer recoverPanic()
		result := helper.ConvertTimeDurationToString(test.Value)
		assert.Equal(t, test.Expects, result, DEFAULT_ASSERTED_TEXT, test.Expects, result)
	}
}
