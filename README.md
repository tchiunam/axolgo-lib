# axolgo-lib, the Axolotl Library in Golang
[![Go](https://github.com/tchiunam/axolgo-lib/actions/workflows/go.yml/badge.svg)](https://github.com/tchiunam/axolgo-lib/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/tchiunam/axolgo-lib/branch/main/graph/badge.svg?token=B5DNGRMYUG)](https://codecov.io/gh/tchiunam/axolgo-lib)
[![CodeQL](https://github.com/tchiunam/axolgo-lib/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/tchiunam/axolgo-lib/actions/workflows/codeql-analysis.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

This is the library of the Axolotl series in Golang. It serves as the base of other axol libraries.

Go package: https://pkg.go.dev/github.com/tchiunam/axolgo-lib

## Use it with your Go module
To add as dependency for your package or upgrade to the latest version:
```
go get github.com/tchiunam/axolgo-lib
```

To upgrade or downgrade to a specific version:
```
go get github.com/tchiunam/axolgo-lib@v1.2.3
```

To remove dependency on your module and downgrade modules:
```
go get github.com/tchiunam/axolgo-lib@none
```

See 'go help get' or https://golang.org/ref/mod#go-get for details.

## Run test
To run test:
```
go test ./...
```

To run test with coverage result:
```
go test -coverpkg=./... ./...
```

## Test report
## Code Coverage graph
![Code Coverage graph](https://codecov.io/gh/tchiunam/axolgo-lib/branch/main/graphs/tree.svg?token=B5DNGRMYUG)

---
#### See more  
1. [axolgo-cloud](https://github.com/tchiunam/axolgo-cloud) for using cloud library (AWS SDK and GCP API)
2. [axolgo-cli](https://github.com/tchiunam/axolgo-cli) for using Axolgo in command line
