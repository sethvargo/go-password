package password_test

import (
	"fmt"
	"log"

	"github.com/sethvargo/go-password/password"
)

func ExampleGenerate() {
	res, err := password.Generate(64, 10, 10, false, false, "")
	if err != nil {
		log.Fatal(err)
	}
	log.Print(res)
}

func ExampleMustGenerate() {
	// Will panic on error
	res := password.MustGenerate(64, 10, 10, false, false, "")
	log.Print(res)
}

func ExampleGenerator_Generate() {
	gen, err := password.NewGenerator(nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := gen.Generate(64, 10, 10, false, false, "")
	if err != nil {
		log.Fatal(err)
	}
	log.Print(res)
}

func ExampleNewGenerator_nil() {
	// This is exactly the same as calling "Generate" directly. It will use all
	// the default values.
	gen, err := password.NewGenerator(nil)
	if err != nil {
		log.Fatal(err)
	}

	_ = gen // gen.Generate(...)
}

func ExampleNewGenerator_custom() {
	// Customize the list of symbols.
	gen, err := password.NewGenerator(&password.GeneratorInput{
		Symbols: "!@#$%^()",
	})
	if err != nil {
		log.Fatal(err)
	}

	_ = gen // gen.Generate(...)
}

func ExampleNewMockGenerator_testing() {
	// Accept a password.PasswordGenerator interface instead of a
	// password.Generator struct.
	f := func(g password.PasswordGenerator) string {
		// These values don't matter
		return g.MustGenerate(1, 2, 3, false, false, "")
	}

	// In tests
	gen := password.NewMockGenerator("canned-response", nil)

	fmt.Print(f(gen))
	// Output: canned-response
}
