package model

import (
	"regexp"

	"lmm/api/domain/model"
	"lmm/api/service/user/domain"
	"lmm/api/util/mathutil"
)

const (
	weakPasswordThreshold = 20
	minimumPasswordLength = 8
)

var (
	digits           = regexp.MustCompile(`\d`)
	lowerCaseLetters = regexp.MustCompile(`[a-z]`)
	upperCaseLetters = regexp.MustCompile(`[A-Z]`)
	symbols          = regexp.MustCompile(`[%]`)
)

// Password domain value object model
type Password struct {
	model.ValueObject
	text string
}

// NewPassword creates a new password value object
// would returns error if password is empty or contains invalid characters
func NewPassword(text string) (*Password, error) {
	if text == "" {
		return nil, domain.ErrUserPasswordEmpty
	}
	if len(text) < minimumPasswordLength {
		return nil, domain.ErrUserPasswordTooShort
	}
	return &Password{text: text}, nil
}

func (pw Password) String() string {
	return pw.text
}

// IsWeak returns true if the password's strength is weak
func (pw Password) IsWeak() bool {
	return pw.calculateStrength() <= weakPasswordThreshold
}

func (pw Password) calculateStrength() int {
	countDigit := len(digits.FindAllString(pw.text, -1))
	countLower := len(lowerCaseLetters.FindAllString(pw.text, -1))
	countUpper := len(upperCaseLetters.FindAllString(pw.text, -1))
	countSymbol := len(symbols.FindAllString(pw.text, -1))

	countLetter := countLower + countUpper
	total := countDigit + countLetter + countSymbol

	if countDigit > 0 && countLetter > 0 {
		total += mathutil.MinInt(countDigit, countLetter) * 2
	}

	if countLower > 0 && countUpper > 0 {
		total += mathutil.MinInt(countLower, countUpper) * 5
	}

	if countDigit > 0 && countLetter > 0 && countSymbol > 0 {
		total += mathutil.MinInt(countDigit, countLetter, countSymbol) * 5
	}

	return total
}
