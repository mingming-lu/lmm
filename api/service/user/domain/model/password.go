package model

import (
	"fmt"
	"regexp"

	"lmm/api/domain/model"
	"lmm/api/service/user/domain"
	"lmm/api/util/mathutil"
)

const (
	minimumPasswordLength = 8
	maximumPasswordLength = 250
	weakPasswordThreshold = 20

	_digits    = `0-9`
	_lowerCase = `a-z`
	_upperCase = `A-Z`
	_symbols   = `~!@#$%^*(){}\[\]&\-_=+|\\\/:;"'<,>.?`
	_password  = _digits + _lowerCase + _upperCase + _symbols
)

var (
	digits           = regexp.MustCompile(fmt.Sprintf(`[%s]`, _digits))
	lowerCaseLetters = regexp.MustCompile(fmt.Sprintf(`[%s]`, _lowerCase))
	upperCaseLetters = regexp.MustCompile(fmt.Sprintf(`[%s]`, _upperCase))
	symbols          = regexp.MustCompile(fmt.Sprintf(`[%s]`, _symbols))
	patternPassword  = regexp.MustCompile(fmt.Sprintf(`[%s]{%d,%d}`, _password, minimumPasswordLength, maximumPasswordLength))
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
	if len(text) > maximumPasswordLength {
		return nil, domain.ErrUserPasswordTooLong
	}
	if !patternPassword.MatchString(text) {
		return nil, domain.ErrInvalidPassword
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
		total += mathutil.MinInt(2, mathutil.MinInt(countDigit, countLetter)) * 5
	}

	if countLower > 0 && countUpper > 0 {
		total += mathutil.MinInt(2, mathutil.MinInt(countLower, countUpper)) * 5
	}

	if countDigit > 0 && countLetter > 0 && countSymbol > 0 {
		total += mathutil.MinInt(2, mathutil.MinInt(countDigit, countLetter, countSymbol)) * 5
	}

	return total
}
