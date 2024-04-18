package password

import (
	"errors"
	"io"
	"strings"
	"sync/atomic"
	"testing"
)

type (
	MockReader struct {
		Counter int64
	}
)

const (
	N = 10000
)

func (mr *MockReader) Read(data []byte) (int, error) {
	for i := 0; i < len(data); i++ {
		data[i] = byte(atomic.AddInt64(&mr.Counter, 1))
	}
	return len(data), nil
}

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

func testGeneratorGenerate(t *testing.T, reader io.Reader) {
	t.Helper()

	gen, err := NewGenerator(nil)
	if reader != nil {
		gen.reader = reader
	}
	if err != nil {
		t.Fatal(err)
	}

	t.Run("exceeds_length", func(t *testing.T) {
		t.Parallel()

		if _, err := gen.Generate(0, 1, 0, false, false); !errors.Is(err, ErrExceedsTotalLength) {
			t.Errorf("expected %q to be %q", err, ErrExceedsTotalLength)
		}

		if _, err := gen.Generate(0, 0, 1, false, false); !errors.Is(err, ErrExceedsTotalLength) {
			t.Errorf("expected %q to be %q", err, ErrExceedsTotalLength)
		}
	})

	t.Run("exceeds_letters_available", func(t *testing.T) {
		t.Parallel()

		if _, err := gen.Generate(1000, 0, 0, false, false); !errors.Is(err, ErrLettersExceedsAvailable) {
			t.Errorf("expected %q to be %q", err, ErrLettersExceedsAvailable)
		}
	})

	t.Run("exceeds_digits_available", func(t *testing.T) {
		t.Parallel()

		if _, err := gen.Generate(52, 11, 0, false, false); !errors.Is(err, ErrDigitsExceedsAvailable) {
			t.Errorf("expected %q to be %q", err, ErrDigitsExceedsAvailable)
		}
	})

	t.Run("exceeds_symbols_available", func(t *testing.T) {
		t.Parallel()

		if _, err := gen.Generate(52, 0, 31, false, false); !errors.Is(err, ErrSymbolsExceedsAvailable) {
			t.Errorf("expected %q to be %q", err, ErrSymbolsExceedsAvailable)
		}
	})

	t.Run("gen_lowercase", func(t *testing.T) {
		t.Parallel()

		for i := 0; i < N; i++ {
			res, err := gen.Generate(i%len(LowerLetters), 0, 0, true, true)
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

		res, err := gen.Generate(1000, 0, 0, false, true)
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
			res, err := gen.Generate(52, 10, 30, false, false)
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
	testGeneratorGenerate(t, nil)
}

func TestGenerator_Reader_Generate(t *testing.T) {
	t.Parallel()
	testGeneratorGenerate(t, &MockReader{})
}

func testGeneratorGenerateCustom(t *testing.T, reader io.Reader) {
	t.Helper()

	gen, err := NewGenerator(&GeneratorInput{
		LowerLetters: "abcde",
		UpperLetters: "ABCDE",
		Symbols:      "!@#$%",
		Digits:       "01234",
		Reader:       reader,
	})
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < N; i++ {
		res, err := gen.Generate(52, 10, 10, false, true)
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
	testGeneratorGenerateCustom(t, nil)
}

func TestGenerator_Reader_Generate_Custom(t *testing.T) {
	t.Parallel()
	testGeneratorGenerateCustom(t, &MockReader{})
}
