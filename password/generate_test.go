package password

import (
	"log"
	"strings"
	"testing"
)

func testHasDuplicates(tb testing.TB, s string) bool {
	found := make(map[rune]struct{}, len(s))
	for _, ch := range s {
		if _, ok := found[ch]; ok {
			return true
		}
		found[ch] = struct{}{}
	}
	return false
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	t.Run("exceeds_length", func(t *testing.T) {
		t.Parallel()

		if _, err := Generate(0, 1, 0, false, false); err != ErrExceedsTotalLength {
			t.Errorf("expected %q to be %q", err, ErrExceedsTotalLength)
		}

		if _, err := Generate(0, 0, 1, false, false); err != ErrExceedsTotalLength {
			t.Errorf("expected %q to be %q", err, ErrExceedsTotalLength)
		}
	})

	t.Run("exceeds_letters_available", func(t *testing.T) {
		t.Parallel()

		if _, err := Generate(1000, 0, 0, false, false); err != ErrLettersExceedsAvailable {
			t.Errorf("expected %q to be %q", err, ErrLettersExceedsAvailable)
		}
	})

	t.Run("exceeds_digits_available", func(t *testing.T) {
		t.Parallel()

		if _, err := Generate(52, 11, 0, false, false); err != ErrDigitsExceedsAvailable {
			t.Errorf("expected %q to be %q", err, ErrDigitsExceedsAvailable)
		}
	})

	t.Run("exceeds_symbols_available", func(t *testing.T) {
		t.Parallel()

		if _, err := Generate(52, 0, 31, false, false); err != ErrSymbolsExceedsAvailable {
			t.Errorf("expected %q to be %q", err, ErrSymbolsExceedsAvailable)
		}
	})

	t.Run("gen_lowercase", func(t *testing.T) {
		t.Parallel()

		for i := 0; i < 10000; i++ {
			res, err := Generate(i%len(LowerLetters), 0, 0, true, true)
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

		res, err := Generate(1000, 0, 0, false, true)
		if err != nil {
			t.Error(err)
		}

		if res == strings.ToLower(res) {
			t.Errorf("%q does not include uppercase", res)
		}
	})

	t.Run("gen_no_repeats", func(t *testing.T) {
		t.Parallel()

		for i := 0; i < 10000; i++ {
			res, err := Generate(52, 10, 30, false, false)
			if err != nil {
				t.Error(err)
			}

			if testHasDuplicates(t, res) {
				t.Errorf("%q should not have duplicates", res)
			}
		}
	})
}

func ExampleGenerate() {
	res, err := Generate(64, 10, 10, false, false)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(res)
}
