# env

[![CI](https://github.com/maxcraig112/env/actions/workflows/ci.yml/badge.svg)](https://github.com/maxcraig112/env/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/maxcraig112/env.svg)](https://pkg.go.dev/github.com/maxcraig112/env)

Small, dependency-free helpers around `os.Getenv` for reading and requiring
environment variables in Go, with optional typed defaults.

- `Get` / `Require` for plain strings
- Typed variants: `GetInt`, `GetBool`, `GetDuration`, etc.
- Defaults are passed as a variadic argument; the first one given is used
- Failures are reported as errors, never panics
- No third-party dependencies, no reflection, no struct tags

## Install

```sh
go get github.com/maxcraig112/env
```

## Usage

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/maxcraig112/env"
)

func main() {
	// Returns the value if set, otherwise "" (no default given).
	name := env.Get("APP_NAME")

	// Returns the value if set, otherwise the first default.
	logLevel := env.Get("LOG_LEVEL", "info")

	// Returns a *env.NotSetError if APP_SECRET is not set.
	secret, err := env.Require("APP_SECRET")
	if err != nil {
		log.Fatal(err)
	}

	// Typed helpers work the same way, with typed defaults. They return an
	// error if the value is set but cannot be parsed.
	port, err := env.GetInt("PORT", 8080)
	if err != nil {
		log.Fatal(err)
	}
	debug, err := env.GetBool("DEBUG", false)
	if err != nil {
		log.Fatal(err)
	}
	timeout, err := env.GetDuration("TIMEOUT", 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	// Typed Require* variants return an error on missing or malformed values.
	maxConns, err := env.RequireInt("MAX_CONNECTIONS")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(name, logLevel, secret, port, debug, timeout, maxConns)
}
```

## API

### Strings

```go
func Get(key string, defaults ...string) string
func Require(key string) (string, error)
```

`Get` returns the value of `key` if set. Otherwise it returns the first
value in `defaults`, or `""` if no defaults are given; it never errors.
`Require` returns the value of `key`, along with a `*env.NotSetError` if it
isn't set.

### Typed helpers

The same `Get` / `Require` pattern is available for other common types,
with identical behavior across types:

| Type            | Get           | Require            |
| --------------- | ------------- | ------------------- |
| `int`           | `GetInt`      | `RequireInt`         |
| `int64`         | `GetInt64`    | `RequireInt64`       |
| `float64`       | `GetFloat64`  | `RequireFloat64`     |
| `bool`          | `GetBool`     | `RequireBool`        |
| `time.Duration` | `GetDuration` | `RequireDuration`    |

Each `Get<Type>` has the signature
`func(key string, defaults ...Type) (Type, error)`, and each `Require<Type>`
has the signature `func(key string) (Type, error)`.

If the variable is unset, `Get<Type>` returns the first default (or the zero
value of `Type` if none is given) with a `nil` error; an unset variable
with a fallback is not an error. If the variable is set but cannot be parsed
as the target type, both the `Get*` and `Require*` variants return a
`*env.ParseError`. `Require<Type>` additionally returns a `*env.NotSetError`
when the variable is unset.

### Errors

```go
type NotSetError struct{ Key string }
type ParseError struct {
	Key   string
	Value string
	Type  string
	Err   error
}
```

Both implement `error`, and `ParseError` implements `Unwrap() error`, so
both work with `errors.As` / `errors.Is`:

```go
port, err := env.RequireInt("PORT")
var parseErr *env.ParseError
if errors.As(err, &parseErr) {
	log.Fatalf("PORT=%q is not a valid int", parseErr.Value)
}
```

## Why errors instead of panics?

Environment configuration is typically read once, at process startup. This
package reports every failure (a missing required variable, or a value
that doesn't parse as the requested type) as a returned `error`, so callers
can decide how to handle it (log and exit, wrap with more context, retry,
etc.) using ordinary Go error handling rather than `recover`.

## License

Apache License 2.0. See [LICENSE](LICENSE).
