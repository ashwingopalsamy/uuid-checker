# uuidcheck

[![Go Reference](https://pkg.go.dev/badge/github.com/ashwingopalsamy/uuidcheck.svg)](https://pkg.go.dev/github.com/ashwingopalsamy/uuidcheck)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go) 
[![Go Report Card](https://goreportcard.com/badge/github.com/ashwingopalsamy/uuidcheck)](https://goreportcard.com/report/github.com/ashwingopalsamy/uuidcheck)
[![Coverage Status](https://codecov.io/gh/ashwingopalsamy/uuidcheck/branch/master/graph/badge.svg)](https://codecov.io/gh/ashwingopalsamy/uuidcheck)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)


A tiny, zero-dependency Go library that validates UUIDs against standard RFC 4122 formatting, converts UUIDv7() into timestamps by ensuring accuracy and light compute.

**Why?**
Sometimes you just need to know if that string is a proper UUID without dragging in heavy libraries or writing clunky checks. `uuidcheck` does one thing and does it well.

## Features

- **Simple & Light:** No regular expressions, no external dependencies and a single-pass check.
- **Strict RFC 4122 Format:** Ensures correct length, hyphen positions and valid hex characters.
- **UUIDv7 Support:** Extracts embedded timestamps from version 7 UUIDs.
- **Fully Tested:** Includes comprehensive unit tests, covering a range of edge cases.

## Getting Started

```bash
go get github.com/ashwingopalsamy/uuidcheck
```

## How It Works

`IsValidUUID` runs a quick series of checks:

- **Length Check:** Must be exactly 36 characters.
- **Hyphen Positions:** Hyphens must appear at positions 8, 13, 18, and 23.
- **Hex Digits:** All other characters must be valid hex (`0-9`, `A-F`, `a-f`).

`IsUUIDv7` checks the version nibble of the `time_hi_and_version` field, ensuring its '7'.

`UUIDv7ToTimestamp` extracts the first 48 bits from the UUID (the combination of `time_low` and part of `time_mid`) and interprets them as a Unix timestamp in milliseconds.

## Examples

**Valid:**

- `01939c00-282d-782f-9cc2-887dc7b40629`
- `01939C00-282D-782F-9CC2-887DC7B40629`

**Invalid:**

- `01939c-282d-782f-9cc2-887` (too short)
- `f01939c00-282d-782f-9cg2-887dc7b40629` (invalid hex char `g`)
- `01939c00282d782f9cc2887dc7b40629` (no hyphens)

## Testing

We believe in solid test coverage. Just run:

```bash
go test -v ./...
```

You'll find unit tests and edge case scenarios in `uuidcheck_test.go`.

## Contributing

Contributions are welcome!
Feel free to open issues, submit PRs, or propose features. Just keep it simple and aligned with the library’s goal: blazing-fast, straightforward UUID validation.

## License

This project is licensed under the [MIT License](LICENSE).
