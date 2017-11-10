## Golang Password Generator

[![GoDoc](https://godoc.org/github.com/sethvargo/go-password?status.svg)](https://godoc.org/github.com/sethvargo/go-password)

This library implements generation of random passwords with provided
requirements as described by  [AgileBits
1Password](https://discussions.agilebits.com/discussion/23842/how-random-are-the-generated-passwords)
in pure Golang. The algorithm is commonly used when generating website
passwords.

The library uses crypto/rand for added randomness.

## Installation

```sh
$ go -u get github.com/sethvargo/password
```

## Usage

```golang
package main

import (
  "strings"

  "github.com/sethvargo/go-password/password"
)

func main() {
  // Generate a password that is 64 characters long with 10 digits, 10 symbols,
  // allowing upper and lower case letters, disallowing repeat characters.
  res, err := Generate(64, 10, 10, false, false)
  if err != nil {
    log.Fatal(err)
  }
  log.Printf(res)
}
```

See the [GoDoc](https://godoc.org/github.com/sethvargo/go-password) for more
information.

## License

This code is licensed under the MIT license.
