package model

import (
	"strings"

	"lmm/api/service/article/domain"
	"lmm/api/testing"

	"github.com/google/uuid"
)

func TestNewArticleID(tt *testing.T) {

	minLength := 8
	maxLength := 64

	tt.Run("Valid", func(tt *testing.T) {
		type Case struct {
			Input  string
			Output string
		}

		cases := map[string]Case{
			"AllNumber": Case{
				Input:  strings.Repeat("1234", 8),
				Output: strings.Repeat("1234", 8),
			},
			"AllLower": Case{
				Input:  strings.Repeat("abcd", 8),
				Output: strings.Repeat("abcd", 8),
			},
			"AllUpper": Case{
				Input:  strings.Repeat("LGTM", 8),
				Output: strings.Repeat("lgtm", 8),
			},
			"HasUpper": Case{
				Input:  strings.Repeat("AbCdE", 7),
				Output: strings.Repeat("abcde", 7),
			},
			"HasNumber": Case{
				Input:  strings.Repeat("plan9", 7),
				Output: strings.Repeat("plan9", 7),
			},
			"HasHypen": Case{
				Input:  strings.Repeat("b-tree", 6),
				Output: strings.Repeat("b-tree", 6),
			},
			"MinLength": Case{
				Input:  strings.Repeat("f", minLength),
				Output: strings.Repeat("f", minLength),
			},
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
			"UnallowedChar": Case{Input: strings.Repeat("!", minLength), Error: domain.ErrInvalidArticleID},
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

func TestSetArticleIDAlias(tt *testing.T) {
	t := testing.NewTester(tt)

	rawID := uuid.New().String()
	id, err := NewArticleID(rawID)

	t.NoError(err)
	t.NotNil(id)
	t.Is(rawID, id.String())

	t.NoError(id.SetAlias("awesome-article"))
	t.Is("awesome-article", id.String())

	t.Error(domain.ErrInvalidAliasArticleID, id.SetAlias("!!???"))
	t.Is("awesome-article", id.String())
}
