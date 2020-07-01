package password

import (
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
	found := make(map[rune]struct{}, len(s))
	for _, ch := range s {
		if _, ok := found[ch]; ok {
			return true
		}
		found[ch] = struct{}{}
	}
	return false
}

func testGenerator_Generate(t *testing.T, reader io.Reader) {
	t.Parallel()

	gen, err := NewGenerator(nil)
	if reader != nil {
		gen.reader = reader
	}
	if err != nil {
		t.Fatal(err)
	}

	t.Run("exceeds_length", func(t *testing.T) {
		t.Parallel()

		if _, err := gen.Generate(0, 1, 0, 0, false); err != ErrExceedsTotalLength {
			t.Errorf("expected %q to be %q", err, ErrExceedsTotalLength)
		}

		if _, err := gen.Generate(0, 0, 1, 0, false); err != ErrExceedsTotalLength {
			t.Errorf("expected %q to be %q", err, ErrExceedsTotalLength)
		}
	})

	t.Run("exceeds_lowerletters_available", func(t *testing.T) {
		t.Parallel()

		if _, err := gen.Generate(1000, 0, 0, 0, false); err != ErrLowerLettersExceedsAvailable {
			t.Errorf("expected %q to be %q", err, ErrLowerLettersExceedsAvailable)
		}
	})

	t.Run("exceeds_upperletters_available", func(t *testing.T) {
		t.Parallel()

		if _, err := gen.Generate(1000, 0, 0, 1000, false); err != ErrUpperLettersExceedsAvailable {
			t.Errorf("expected %q to be %q", err, ErrUpperLettersExceedsAvailable)
		}
	})

	t.Run("exceeds_digits_available", func(t *testing.T) {
		t.Parallel()

		if _, err := gen.Generate(52, 11, 0, 0, false); err != ErrDigitsExceedsAvailable {
			t.Errorf("expected %q to be %q", err, ErrDigitsExceedsAvailable)
		}
	})

	t.Run("exceeds_symbols_available", func(t *testing.T) {
		t.Parallel()

		if _, err := gen.Generate(52, 0, 31, 0, false); err != ErrSymbolsExceedsAvailable {
			t.Errorf("expected %q to be %q", err, ErrSymbolsExceedsAvailable)
		}
	})

	t.Run("gen_no_repeats", func(t *testing.T) {
		t.Parallel()

		for i := 0; i < N; i++ {
			res, err := gen.Generate(92, 10, 30, 26, false)
			if err != nil {
				t.Error(err)
			}

			if testHasDuplicates(t, res) {
				t.Errorf("%q should not have duplicates", res)
			}
		}
	})
}

func TestGenerator_Generate(t *testing.T) {
	testGenerator_Generate(t, nil)
}

func TestGenerator_Reader_Generate(t *testing.T) {
	testGenerator_Generate(t, &MockReader{})
}

func testGenerator_Generate_Custom(t *testing.T, reader io.Reader) {
	t.Parallel()

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
		res, err := gen.Generate(52, 10, 10, 26, true)
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

func TestGenerator_Generate_Custom(t *testing.T) {
	testGenerator_Generate_Custom(t, nil)
}

func TestGenerator_Reader_Generate_Custom(t *testing.T) {
	testGenerator_Generate_Custom(t, &MockReader{})
}
