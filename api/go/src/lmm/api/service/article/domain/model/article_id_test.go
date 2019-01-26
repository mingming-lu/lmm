package model

import (
	"strings"

	"lmm/api/service/article/domain"
	"lmm/api/testing"
)

func TestNewArticleID(tt *testing.T) {

	minLength := 8
	maxLength := 80

	tt.Run("Valid", func(tt *testing.T) {
		type Case struct {
			Input  string
			Output string
		}

		cases := map[string]Case{
			"AllNumber": Case{Input: "12345678", Output: "12345678"},
			"AllLower":  Case{Input: "aaaaaaaa", Output: "aaaaaaaa"},
			"AllUpper":  Case{Input: "AAAAAAAA", Output: "aaaaaaaa"},
			"HasUpper":  Case{Input: "AaAaAAaa", Output: "aaaaaaaa"},
			"HasNumber": Case{Input: "aBcD1234", Output: "abcd1234"},
			"HasHypen":  Case{Input: "aaaa-aaaa", Output: "aaaa-aaaa"},
			"MaxLength": Case{
				Input:  strings.Repeat("z", maxLength),
				Output: strings.Repeat("z", maxLength),
			},
		}

		for testname, testcase := range cases {
			tt.Run(testname, func(tt *testing.T) {
				t := testing.NewTester(tt)

				id, err := NewArticleID(testcase.Input)
				t.NoError(err)
				t.NotNil(id)
				t.Is(testcase.Output, id.String())
			})
		}
	})

	tt.Run("Invalid", func(tt *testing.T) {
		type Case struct {
			Input string
			Error error
		}

		cases := map[string]Case{
			"TooShort":      Case{Input: strings.Repeat("v", minLength-1), Error: domain.ErrInvalidArticleID},
			"TooLong":       Case{Input: strings.Repeat("h", maxLength+1), Error: domain.ErrInvalidArticleID},
			"UnallowedChar": Case{Input: "not-avaiable-char!", Error: domain.ErrInvalidArticleID},
		}

		for testname, testcase := range cases {
			tt.Run(testname, func(tt *testing.T) {
				t := testing.NewTester(tt)

				id, err := NewArticleID(testcase.Input)
				t.IsError(testcase.Error, err)
				t.Nil(id)
			})
		}
	})
}
