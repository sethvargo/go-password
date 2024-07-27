// Package password provides a library for generating high-entropy random
// password strings via the crypto/rand package.
//
//	res, err := Generate(64, 10, 10, false, false)
//	if err != nil  {
//	  log.Fatal(err)
//	}
//	log.Printf(res)
//
// Most functions are safe for concurrent use.
package password

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

const (
	// LowerLetters is the list of lowercase letters.
	LowerLetters = "abcdefghijklmnopqrstuvwxyz"

	// UpperLetters is the list of uppercase letters.
	UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Digits is the list of permitted digits.
	Digits = "0123456789"

	// Symbols is the list of symbols.
	Symbols = "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"
)

var (
	// ErrExceedsTotalLength is the error returned with the number of digits and
	// symbols is greater than the total length.
	ErrExceedsTotalLength = errors.New("number of digits and symbols must be less than total length")

	// ErrLettersExceedsAvailable is the error returned with the number of letters
	// exceeds the number of available letters and repeats are not allowed.
	ErrLettersExceedsAvailable = errors.New("number of letters exceeds available letters and repeats are not allowed")

	// ErrDigitsExceedsAvailable is the error returned with the number of digits
	// exceeds the number of available digits and repeats are not allowed.
	ErrDigitsExceedsAvailable = errors.New("number of digits exceeds available digits and repeats are not allowed")

	// ErrSymbolsExceedsAvailable is the error returned with the number of symbols
	// exceeds the number of available symbols and repeats are not allowed.
	ErrSymbolsExceedsAvailable = errors.New("number of symbols exceeds available symbols and repeats are not allowed")
)

// Generator is the stateful generator which can be used to customize the list
// of letters, digits, and/or symbols.
type Generator struct {
	lowerLetters string
	upperLetters string
	digits       string
	symbols      string
}

// Input used to define input parameters for the generator
type Input struct {
	Length      int
	Digits      int
	Symbols     int
	NoUpper     bool
	AllowRepeat bool
	_           struct{}
}

// NewGenerator creates a new Generator from the specified configuration. If no
// input is given, all the default values are used. This function is safe for
// concurrent use.
func NewGenerator() Generator {
	return Generator{
		lowerLetters: LowerLetters,
		upperLetters: UpperLetters,
		digits:       Digits,
		symbols:      Symbols,
	}
}

// WithLowerLetters creates a new Generator from another Generator with specific
// LowerLetters
func (g Generator) WithLowerLetters(lowerLetters string) Generator {
	g.lowerLetters = lowerLetters
	return g
}

// WithUpperLetters creates a new Generator from another Generator with specific
// UpperLetters
func (g Generator) WithUpperLetters(upperLetters string) Generator {
	g.upperLetters = upperLetters
	return g
}

// WithDigits creates a new Generator from another Generator with specific
// Digits
func (g Generator) WithDigits(digits string) Generator {
	g.digits = digits
	return g
}

// WithSymbols creates a new Generator from another Generator with specific
// Symbols
func (g Generator) WithSymbols(symbols string) Generator {
	g.symbols = symbols
	return g
}

// Generate generates a password with the given requirements. length is the
// total number of characters in the password. numDigits is the number of digits
// to include in the result. numSymbols is the number of symbols to include in
// the result. noUpper excludes uppercase letters from the results. allowRepeat
// allows characters to repeat.
//
// The algorithm is fast, but it's not designed to be performant; it favors
// entropy over speed. This function is safe for concurrent use.
func (g Generator) Generate(input Input) (string, error) {
	letters := g.lowerLetters
	if !input.NoUpper {
		letters += g.upperLetters
	}

	chars := input.Length - input.Digits - input.Symbols
	if chars < 0 {
		return "", ErrExceedsTotalLength
	}

	if !input.AllowRepeat && chars > len(letters) {
		return "", ErrLettersExceedsAvailable
	}

	if !input.AllowRepeat && input.Digits > len(g.digits) {
		return "", ErrDigitsExceedsAvailable
	}

	if !input.AllowRepeat && input.Symbols > len(g.symbols) {
		return "", ErrSymbolsExceedsAvailable
	}

	var result string

	// Characters
	for i := 0; i < chars; i++ {
		ch, err := randomElement(letters)
		if err != nil {
			return "", err
		}

		if !input.AllowRepeat && strings.Contains(result, ch) {
			i--
			continue
		}

		result, err = randomInsert(result, ch)
		if err != nil {
			return "", err
		}
	}

	// Digits
	for i := 0; i < input.Digits; i++ {
		d, err := randomElement(g.digits)
		if err != nil {
			return "", err
		}

		if !input.AllowRepeat && strings.Contains(result, d) {
			i--
			continue
		}

		result, err = randomInsert(result, d)
		if err != nil {
			return "", err
		}
	}

	// Symbols
	for i := 0; i < input.Symbols; i++ {
		sym, err := randomElement(g.symbols)
		if err != nil {
			return "", err
		}

		if !input.AllowRepeat && strings.Contains(result, sym) {
			i--
			continue
		}

		result, err = randomInsert(result, sym)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}

// MustGenerate is the same as Generate, but panics on error.
func (g Generator) MustGenerate(input Input) string {
	res, err := g.Generate(input)
	if err != nil {
		panic(err)
	}
	return res
}

// Generate is the package shortcut for Generator.Generate.
func Generate(input Input) (string, error) {
	return NewGenerator().Generate(input)
}

// MustGenerate is the package shortcut for Generator.MustGenerate.
func MustGenerate(input Input) string {
	res, err := Generate(input)
	if err != nil {
		panic(err)
	}
	return res
}

// randomInsert randomly inserts the given value into the given string.
func randomInsert(s, val string) (string, error) {
	if s == "" {
		return val, nil
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(s)+1)))
	if err != nil {
		return "", fmt.Errorf("failed to generate random integer: %w", err)
	}
	i := n.Int64()
	return s[0:i] + val + s[i:], nil
}

// randomElement extracts a random element from the given string.
func randomElement(s string) (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(s))))
	if err != nil {
		return "", fmt.Errorf("failed to generate random integer: %w", err)
	}
	return string(s[n.Int64()]), nil
}
