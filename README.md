# go-jwkset

[![Go](https://github.com/nabeken/go-jwkset/actions/workflows/go.yml/badge.svg)](https://github.com/nabeken/go-jwkset/actions/workflows/go.yml)
[![BSD License](http://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/nabeken/go-jwkset/blob/master/LICENSE)

go-jwkset is a library to fetch [JSON Web Key](https://datatracker.ietf.org/doc/html/rfc7517) ("JWK") Set on top of [go-jose/go-jose/v4](https://github.com/go-jose/go-jose) library.
go-jwkset allows you to build a cache-ware custom fetcher for JWKSet.

# v2

The `master` branch is for v2 version of the library. Please import the library with the following.

```go
import "github.com/nabeken/go-jwkset/v2"
```

# Built-in fetcher implementation

- Plain HTTP
- AWS S3
- [AWS Application Load Balancer](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/listener-authenticate-users.html)

# Versioning

This library follows [Semantic Versions](http://semver.org/).
