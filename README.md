[![Go Reference](https://pkg.go.dev/badge/github.com/tommzn/go-config.svg)](https://pkg.go.dev/github.com/tommzn/go-config)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/go-config)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/tommzn/go-config)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommzn/go-config)](https://goreportcard.com/report/github.com/tommzn/go-config)
[![Actions Status](https://github.com/tommzn/go-config/actions/workflows/go.pkg.auto-ci.yml/badge.svg)](https://github.com/tommzn/go-config/actions)

# go-config

A Go library for loading and accessing YAML configuration from multiple sources through a single, unified interface. Built on top of [Viper](https://github.com/spf13/viper).

## Features

- Load configuration from local YAML files, in-memory strings, or AWS S3
- Uniform `Config` interface regardless of the source
- Typed accessors: string, int, int slice, bool, duration, slice of maps
- Unmarshal configuration directly into structs
- Automatic file discovery across standard config paths
- Dot-notation access for nested keys (e.g. `"namespace.key"`)
- Pointer-based return values with default value fallback

## Installation

```bash
go get github.com/tommzn/go-config
```

## Quick Start

```go
// Load from the default file (config.yml, searched in standard locations)
source := config.NewConfigSource()
cfg, err := source.Load()
if err != nil {
    log.Fatal(err)
}

value := cfg.Get("mykey", config.AsStringPtr("default"))
fmt.Println(*value)
```

## Configuration Sources

### File

Loads a YAML file from a given path. If no path is provided, it searches for `config.yml` in the following locations (in order):

1. `./`
2. `$HOME/`
3. `$HOME/go_config/`
4. `/etc/go_config/`

```go
// Auto-discover config.yml
source := config.NewFileConfigSource(nil)

// Explicit path
path := "./configs/app.yml"
source := config.NewFileConfigSource(&path)

cfg, err := source.Load()
```

### Static

Loads configuration from an in-memory YAML string. Useful for tests or embedded defaults.

```go
yaml := `
server:
  host: localhost
  port: 8080
`
source := config.NewStaticConfigSource(yaml)
cfg, err := source.Load()
```

### AWS S3

Downloads a YAML file from an S3 bucket. AWS credentials are resolved via the standard AWS SDK credential chain.

```go
// With explicit region
region := "eu-central-1"
source, err := config.NewS3ConfigSource("my-bucket", "configs/app.yml", &region)

// From environment variables: AWS_REGION, GO_CONFIG_S3_BUCKET, GO_CONFIG_S3_KEY
source, err := config.NewS3ConfigSourceFromEnv()

cfg, err := source.Load()
```

**Required environment variables for `NewS3ConfigSourceFromEnv`:**

| Variable | Description |
|---|---|
| `AWS_REGION` | AWS region of the S3 bucket |
| `GO_CONFIG_S3_BUCKET` | Name of the S3 bucket |
| `GO_CONFIG_S3_KEY` | Path and filename of the config file in the bucket |

## Accessing Configuration Values

All accessor methods accept a key and a default value (pointer). If the key is not found, or type conversion fails, the default is returned. All methods return pointers — a `nil` return means the key was missing and no default was given.

Helper functions `AsStringPtr`, `AsIntPtr`, `AsBoolPtr`, and `AsDurationPtr` are provided for convenience when passing defaults.

### String

```go
val := cfg.Get("app.name", config.AsStringPtr("my-app"))
fmt.Println(*val) // "my-app" if key not found
```

### Int

```go
val := cfg.GetAsInt("server.port", config.AsIntPtr(8080))
fmt.Println(*val)
```

### Int Slice

```go
val := cfg.GetAsIntSlice("allowed.ports", nil)
if val != nil {
    fmt.Println(*val) // []int{...}
}
```

### Bool

```go
val := cfg.GetAsBool("feature.enabled", config.AsBoolPtr(false))
fmt.Println(*val)
```

### Duration

Duration values support the following formats:

| Format | Example | Result |
|---|---|---|
| Seconds suffix | `"30s"` | 30 seconds |
| Minutes suffix | `"5m"` | 5 minutes |
| Hours suffix | `"2h"` | 2 hours |
| Plain integer | `"45"` | 45 seconds |

```go
val := cfg.GetAsDuration("cache.ttl", config.AsDurationPtr(60*time.Second))
fmt.Println(*val)
```

### Slice of Maps

Returns a `[]map[string]string` for list-of-object structures in YAML.

```yaml
endpoints:
  - host: service-a
    port: "8080"
  - host: service-b
    port: "9090"
```

```go
entries := cfg.GetAsSliceOfMaps("endpoints")
for _, e := range entries {
    fmt.Println(e["host"], e["port"])
}
```

### Nested Keys

Use dot notation to access nested values:

```yaml
database:
  host: localhost
  port: 5432
```

```go
host := cfg.Get("database.host", nil)
```

### Unmarshal into Struct

Decode the full configuration (or a subtree) into a struct using `mapstructure` tags:

```go
type ServerConfig struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
}

var srv ServerConfig
if err := cfg.Unmarshal(&srv); err != nil {
    log.Fatal(err)
}
fmt.Println(srv.Host, srv.Port)
```

## Example YAML Configuration

```yaml
namespace1:
  key1: value1

key2: value2
key3: 12345
boolval: true

durations:
  seconds: 43s
  minutes: 21m
  hours: 5h
  defaultvalue: 22   # plain int, interpreted as seconds

sliceofmaps:
  - key1_1: val1_1
    key1_2: val1_2
  - key2_1: val2_1
    key2_2: val2_2

intslice:
  - 342543545
  - 3465567
  - 547657
```

## Interfaces

```go
type ConfigSource interface {
    Load() (Config, error)
}

type Config interface {
    Get(key string, defaultValue *string) *string
    GetAsInt(key string, defaultValue *int) *int
    GetAsIntSlice(key string, defaultValue *[]int) *[]int
    GetAsBool(key string, defaultValue *bool) *bool
    GetAsDuration(key string, defaultValue *time.Duration) *time.Duration
    GetAsSliceOfMaps(key string) []map[string]string
    Unmarshal(rawVal any) error
}
```

## Helper Functions

| Function | Description |
|---|---|
| `AsStringPtr(v string) *string` | Returns a pointer to the given string |
| `AsIntPtr(v int) *int` | Returns a pointer to the given int |
| `AsBoolPtr(v bool) *bool` | Returns a pointer to the given bool |
| `AsDurationPtr(v time.Duration) *time.Duration` | Returns a pointer to the given duration |
| `AsDuration(value string) *time.Duration` | Parses a duration string (`"5s"`, `"3m"`, `"2h"`, or plain int) |

## Requirements

- Go 1.25+
- AWS credentials configured (only required for S3 source)
