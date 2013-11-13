GonormCypher
============

Neo4j Go library based on Anorm in the Play Framework.

Cypher is an amazing query language and I don't think that it should be
abstracted away through incomplete "wrappers". This is my attempt at allowing
you to use Cypher to interact with Neo4j in your Go code in a clean, idiomatic
way.

Big ups to [AnormCypher](https://github.com/AnormCypher/AnormCypher) and
[Wes Freeman](https://twitter.com/wefreema), as this is obviously a rough port
of Wes' amazing work.

## Installation

Use `go get` to install GonormCypher:
```
go get github.com/marpaia/GonormCypher
```

## Unit tests

If you're testing this library on your local machine, just run `go test -v`,
otherwise, edit `gonorm_test.go` and change `http://localhost` to whatever you'd
like it to be, and then run `go test -v`.

## External dependencies

This project has no external dependencies other than the Go standard library.

## Documentation

Like most every other Golang project, this projects documentation can be found
on godoc at [godoc.org/github.com/marpaia/GonormCypher](http://godoc.org/github.com/marpaia/GonormCypher).

## Examples

The unit tests are probably the best source of working examples for this
project.

Consider the following example:

```go
package main

import (
    "fmt"
    "github.com/marpaia/GonormCypher"
)

var g *gonorm.Gonorm

func init() {
    g = gonorm.New("http://localhost", 7474)
}

func main() {
    result, err := g.Cypher(`
    MERGE (p1:Person{name:{name1}})
    MERGE (p2:Person{name:{name2}})
    CREATE UNIQUE p1-[:KNOWS]->p2
    RETURN p1.name
    `).On(map[string]interface{}{
        "name1": "Alice",
        "name2": "Bob",
    }).Execute().AsString()

    if err != nil {
        panic(err)
    }

    fmt.Println("The result is:", result)
}
```

This will print `The result is: Alice`

## Contributing

Please contribute and help improve this project!

- Fork the repo
- Make sure the tests pass
- Improve the code
- Make sure your feature has test coverage
- Make sure the tests pass
- Submit a pull request
