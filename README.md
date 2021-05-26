[![Go Reference](https://pkg.go.dev/badge/github.com/tommzn/go-config.svg)](https://pkg.go.dev/github.com/tommzn/go-config)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/go-config)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/tommzn/go-config)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommzn/go-config)](https://goreportcard.com/report/github.com/tommzn/go-config)

# Read & Access Configurations 
Provides different sources to read config and a generic interface to access configurations. Under the hood it used ![Viper Config]/https://github.com/spf13/viper) to load, parse and access configurations.

## Sources
Following sources are available:
- local/static config
- config from YAML files
- config from files stored in AWS S3 bucket


