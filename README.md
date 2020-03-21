# go-jwkset

[![Build Status](https://img.shields.io/travis/nabeken/go-jwkset/master.svg)](https://travis-ci.org/nabeken/go-jwkset)
[![BSD License](http://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/nabeken/go-jwkset/blob/master/LICENSE)

go-jwkset is a library to fetch JSON Web Key Set on top of [square/go-jose.v2](https://gopkg.in/square/go-jose.v2) library.
go-jwkset allows you to build cache-ware custom fetcher for JWKSet.

# Built-in fetcher implementation

- Plain HTTP
- AWS S3
- [AWS Application Load Balancer](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/listener-authenticate-users.html)

# Versioning

This library follows [Semantic Versions](http://semver.org/) and we highly recommend to use some package manager such as `dep` or `glide`.
