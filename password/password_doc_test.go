package password_test

import (
	"log"

	"github.com/sethvargo/go-password/password"
)

func ExampleGenerate() {
	res, err := password.Generate(64, 10, 10, false, false)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(res)
}

func ExampleMustGenerate() {
	// Will panic on error
	res := password.MustGenerate(64, 10, 10, false, false)
	log.Print(res)
}

func ExampleGenerator_Generate() {
	gen := password.NewGenerator()
	res, err := gen.Generate(64, 10, 10, false, false)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(res)
}

func ExampleNewGenerator_nil() {
	// This is exactly the same as calling "Generate" directly. It will use all
	// the default values.
	gen := password.NewGenerator()
	_ = gen // gen.Generate(...)
}

func ExampleNewGenerator_custom() {
	// Customize the list of symbols.
	gen := password.NewGenerator().WithSymbols("!@#$%^()")
	_ = gen // gen.Generate(...)
}
