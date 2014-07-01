go-license
==========

A license management utility for programs written in Golang.

This program handles identifying software licenses and standardizing on a short,
abbreviated name for each known license type.

## Enforcement

License identifier enforcement is not strict. This makes it possible to warn
when an unrecognized license type is used, encouraging either conformance or an
update to the list of known licenses. There is no way we can know all types of
licenses.

## License guessing

This program also provides naive license guessing based on the license body
(text). This makes it easy to just throw a blob of text in and get a
standardized license identifier string out.

## Example

```go
package main

import (
    "fmt"
    "github.com/ryanuber/go-license"
)

func main() {
    // This case will work if there is a guessable license file in the
    // current working directory.
    l, err := license.NewFromDir(".")
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Println(l.Type)

    // This case will do the exact same thing as above, but uses an explicitly
    // set license file name instead of searching for one.
    l, err = license.NewFromFile("./LICENSE")
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Println(l.Type)

    // This case will work in all cases. The license type and the license data
    // are both being set explicitly. This enables one to use any license.
    l = license.New("MyLicense", "My terms go here")
    fmt.Println(l.Type)

    // This call determines if the license in use is recognized by go-license.
    fmt.Println(l.Recognized())
}
```
