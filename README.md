# jsonpath

[![PkgGoDev](https://pkg.go.dev/badge/elliotcourant/jsonpath)](https://pkg.go.dev/elliotcourant/jsonpath)
[![Build Status](https://travis-ci.com/elliotcourant/jsonpath.svg?branch=main)](https://travis-ci.com/elliotcourant/jsonpath)
[![codecov](https://codecov.io/gh/elliotcourant/jsonpath/branch/main/graph/badge.svg)](https://codecov.io/gh/elliotcourant/jsonpath)

Go implementation of jsonpath

This is still a work in progress, updates to come.


## Example

```go
const jsonString = `{
  "firstName": "John",
  "lastName" : "doe",
  "age"      : 26,
  "address"  : {
    "streetAddress": "naist street",
    "city"         : "Nara",
    "postalCode"   : "630-0192"
  },
  "phoneNumbers": [
    {
      "type"  : "iPhone",
      "number": "0123-4567-8888"
    },
    {
      "type"  : "home",
      "number": "0123-4567-8910"
    },
    {
      "type"  : "mobile",
      "number": "0913-8532-8492"
    }
  ]
}`

result, err := Jsonpath([]byte(jsonString), "$.phoneNumbers[*].type")
if err != nil {
    log.Fatal(err)
}

// Output:
// iPhone
// home
// mobile
for _, item := range result {
    fmt.Println(item)
}
```