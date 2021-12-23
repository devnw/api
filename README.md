# `api` Implements a generic API client

[![Build & Test Action Status](https://github.com/devnw/api/actions/workflows/build.yml/badge.svg)](https://github.com/devnw/api/actions)
[![Go Report Card](https://goreportcard.com/badge/go.devnw.com/api)](https://goreportcard.com/report/go.devnw.com/api)
[![codecov](https://codecov.io/gh/devnw/api/branch/main/graph/badge.svg)](https://codecov.io/gh/devnw/api)
[![Go Reference](https://pkg.go.dev/badge/go.devnw.com/api.svg)](#documentation)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

## Installation

To install the package, run:

```bash
    go get -u go.devnw.com/api@latest
```

## Importing

It is recommended to use the package via the following import:

`import . "go.devnw.com/api"`

Using the `.` import allows for functions to be called directly as if the
functions were in the same namespace without the need to append the package
name.

## Documentation

## Benchmarks

To execute the benchmarks, run the following command:

```bash
    go test -bench=. ./...
```

To view benchmarks over time for the `main` branch of the repository they can
be seen on our [Benchmark Report Card].

[Benchmark Report Card]: https://go.devnw.com/api/dev/bench/
