package password

import (
	"errors"
	"strings"
	"testing"
)

const (
	N = 10000
)

func testHasDuplicates(tb testing.TB, s string) bool {
	tb.Helper()

	found := make(map[rune]struct{}, len(s))
	for _, ch := range s {
		if _, ok := found[ch]; ok {
			return true
		}
		found[ch] = struct{}{}
	}
	return false
}

func testGeneratorGenerate(t *testing.T) {
	t.Helper()

	gen := NewGenerator()
	t.Run("exceeds_length", func(t *testing.T) {
		t.Parallel()

		if _, err := gen.Generate(Input{
			Digits: 1,
		}); !errors.Is(err, ErrExceedsTotalLength) {
			t.Errorf("expected %q to be %q", err, ErrExceedsTotalLength)
		}

		if _, err := gen.Generate(Input{
			Symbols: 1,
		}); !errors.Is(err, ErrExceedsTotalLength) {
			t.Errorf("expected %q to be %q", err, ErrExceedsTotalLength)
		}
	})

	t.Run("exceeds_letters_available", func(t *testing.T) {
		t.Parallel()

		if _, err := gen.Generate(Input{
			Length: 1000,
		}); !errors.Is(err, ErrLettersExceedsAvailable) {
			t.Errorf("expected %q to be %q", err, ErrLettersExceedsAvailable)
		}
	})

	t.Run("exceeds_digits_available", func(t *testing.T) {
		t.Parallel()

		if _, err := gen.Generate(Input{
			Length: 52,
			Digits: 11,
		}); !errors.Is(err, ErrDigitsExceedsAvailable) {
			t.Errorf("expected %q to be %q", err, ErrDigitsExceedsAvailable)
		}
	})

	t.Run("exceeds_symbols_available", func(t *testing.T) {
		t.Parallel()

		if _, err := gen.Generate(Input{
			Length:  52,
			Symbols: 31,
		}); !errors.Is(err, ErrSymbolsExceedsAvailable) {
			t.Errorf("expected %q to be %q", err, ErrSymbolsExceedsAvailable)
		}
	})

	t.Run("gen_lowercase", func(t *testing.T) {
		t.Parallel()

		for i := 0; i < N; i++ {
			res, err := gen.Generate(Input{
				Length:      i % len(LowerLetters),
				NoUpper:     true,
				AllowRepeat: true,
			})
			if err != nil {
				t.Error(err)
			}

			if res != strings.ToLower(res) {
				t.Errorf("%q is not lowercase", res)
			}
		}
	})

	t.Run("gen_uppercase", func(t *testing.T) {
		t.Parallel()

		res, err := gen.Generate(Input{
			Length:      1000,
			AllowRepeat: true,
		})
		if err != nil {
			t.Error(err)
		}

		if res == strings.ToLower(res) {
			t.Errorf("%q does not include uppercase", res)
		}
	})

	t.Run("gen_no_repeats", func(t *testing.T) {
		t.Parallel()

		for i := 0; i < N; i++ {
			res, err := gen.Generate(Input{
				Length:  52,
				Digits:  10,
				Symbols: 30,
			})
			if err != nil {
				t.Error(err)
			}

			if testHasDuplicates(t, res) {
				t.Errorf("%q should not have duplicates", res)
			}
		}
	})
}

func TestGeneratorGenerate(t *testing.T) {
	t.Parallel()
	testGeneratorGenerate(t)
}

func TestGenerator_Reader_Generate(t *testing.T) {
	t.Parallel()
	testGeneratorGenerate(t)
}

func testGeneratorGenerateCustom(t *testing.T) {
	t.Helper()

	gen := NewGenerator().
		WithLowerLetters("abcde").
		WithUpperLetters("ABCDE").
		WithSymbols("!@#$%").
		WithDigits("01234")

	for i := 0; i < N; i++ {
		res, err := gen.Generate(Input{
			Length:      52,
			Digits:      10,
			Symbols:     10,
			AllowRepeat: true,
		})
		if err != nil {
			t.Error(err)
		}

		if strings.Contains(res, "f") {
			t.Errorf("%q should only contain lower letters abcde", res)
		}

		if strings.Contains(res, "F") {
			t.Errorf("%q should only contain upper letters ABCDE", res)
		}

		if strings.Contains(res, "&") {
			t.Errorf("%q should only include symbols !@#$%%", res)
		}

		if strings.Contains(res, "5") {
			t.Errorf("%q should only contain digits 01234", res)
		}
	}
}

func TestGeneratorGenerateCustom(t *testing.T) {
	t.Parallel()
	testGeneratorGenerateCustom(t)
}

func TestGenerator_Reader_Generate_Custom(t *testing.T) {
	t.Parallel()
	testGeneratorGenerateCustom(t)
}
