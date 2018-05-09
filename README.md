## Golang Password Generator

[![Build Status](https://travis-ci.org/sethvargo/go-password.svg?branch=master)](https://travis-ci.org/sethvargo/go-password)
[![GoDoc](https://godoc.org/github.com/sethvargo/go-password?status.svg)](https://godoc.org/github.com/sethvargo/go-password)

This library implements generation of random passwords with provided
requirements as described by  [AgileBits
1Password](https://discussions.agilebits.com/discussion/23842/how-random-are-the-generated-passwords)
in pure Golang. The algorithm is commonly used when generating website
passwords.

The library uses crypto/rand for added randomness.

Sample example passwords this library may generate:

```text
0N[k9PhDqmmfaO`p_XHjVv`HTq|zsH4XiH8umjg9JAGJ#\Qm6lZ,28XF4{X?3sHj
7@90|0H7!4p\,c<!32:)0.9N
UlYuRtgqyWEivlXnLeBpZvIQ
Q795Im1VR5h363s48oZGaLDa
wpvbxlsc
```

## Installation

```sh
$ go get -u github.com/sethvargo/go-password/password
```

## Usage

```golang
package main

import (
  "log"

  "github.com/sethvargo/go-password/password"
)

func main() {
  // Generate a password that is 64 characters long with 10 digits, 10 symbols,
  // allowing upper and lower case letters, disallowing repeat characters.
  res, err := password.Generate(64, 10, 10, false, false)
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
