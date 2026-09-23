# JwtView CLI

JwtView is a command-line interface (CLI) tool for viewing and validating JSON Web Tokens (JWTs).

## Features

- View JWT payloads
- Validate JWT signatures
- Easy-to-use CLI commands

## Installation

- Download the latest release from the [releases page](https://github.com/jgfranco17/jwtview-cli/releases)
- Or build from source using Go:

  ```bash
  git clone https://github.com/jgfranco17/jwtview-cli.git
  cd jwtview-cli

  # Build directly
  go build -o jwtview

  # With Just, with environment detection
  just build
  ```

## Usage

```bash
jwtview [flags] <token>
```
