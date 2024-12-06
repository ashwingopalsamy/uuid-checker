// Package uuidcheck provides functions to validate UUID strings and extract timestamps from UUID version 7 values.
// It aims to be a tiny, zero-dependency library focused on correctness, performance, and clarity.
package uuidcheck

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// IsValidUUID checks if the provided string is a syntactically valid UUID according to RFC 4122 format.
//
// A valid UUID is a 36-character string in the form "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", where
// each 'x' is a valid hexadecimal character (0-9, a-f, A-F), and hyphens are strictly placed at
// positions 8, 13, 18, and 23.
//
// For example:
//   - "f47ac10b-58cc-0372-8567-0e02b2c3d479" is a valid UUID
//   - "f47ac10b58cc037285670e02b2c3d479" (no hyphens) is not
//   - "f47ac10b-58cc-0372-8567-0e02b2c3d47z" (invalid hex 'z') is not
func IsValidUUID(uuid string) bool {
	// UUID must be exactly 36 characters: 8-4-4-4-12
	if len(uuid) != 36 {
		return false
	}

	// Set of predefined hyphen positions (8, 13, 18, 23)
	hyphenPositions := map[int]struct{}{
		8: {}, 13: {}, 18: {}, 23: {},
	}

	for i, c := range uuid {
		// If the character is at a hyphen position, it must be '-'
		if _, isHyphen := hyphenPositions[i]; isHyphen {
			if c != '-' {
				return false
			}
			continue
		}

		// Otherwise, it must be a hexadecimal digit (0-9, a-f, A-F)
		if !isHexadecimal(c) {
			return false
		}
	}

	return true
}

// IsUUIDv7 checks if the given UUID is a UUID version 7.
//
// This function assumes the input is already validated by IsValidUUID. A UUIDv7 is identified by the
// version nibble in the time_hi_and_version field (the first hex character of the third UUID section)
// being '7'.
//
// For example:
//   - "01939c00-282d-7f2f-9cc2-887dc7b40629" should return true
//   - "f47ac10b-58cc-0372-8567-0e02b2c3d479" (which might be version 3) will return false
func IsUUIDv7(uuid string) bool {
	// The version nibble is at uuid[14].
	return uuid[14] == '7'
}

// UUIDv7ToTimestamp extracts the Unix timestamp (in milliseconds since epoch) embedded in a UUIDv7
// and returns it as a time.Time in UTC.
//
// UUIDv7 encodes a 60-bit Unix timestamp in the first 60 bits of the UUID. This function assumes a
// correctly formatted and valid UUIDv7 string:
//   - Exactly 36 characters: 8-4-4-4-12 (with hyphens)
//   - Hex digits in all non-hyphen positions
//   - The version nibble in time_hi_and_version set to 7
//
// It returns an error if parsing fails or if the UUID does not contain a valid timestamp.
//
// Example:
//   - Given a valid UUIDv7 "01939c00-282d-7f2f-9cc2-887dc7b40629", this function returns
//     a time.Time corresponding to the timestamp encoded.
//
// Note: The extracted timestamp corresponds to when the UUID was generated (or intended to be generated),
// providing a sortable and roughly chronological ordering of UUIDs.
func UUIDv7ToTimestamp(uuid string) (time.Time, error) {
	parts := strings.Split(uuid, "-")
	if len(parts) < 2 {
		return time.Time{}, fmt.Errorf("invalid UUID format")
	}

	// Concatenate parts[0] (8 hex chars) and the first 4 hex chars of parts[1], providing 12 hex chars total.
	highBitsHex := parts[0] + parts[1][:4]

	// Parse as a 48-bit hex number, representing milliseconds since epoch.
	timestamp, err := strconv.ParseUint(highBitsHex, 16, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	// Convert milliseconds since Unix epoch to time.Time in UTC.
	t := time.UnixMilli(int64(timestamp)).UTC()
	return t, nil
}

// isHexadecimal checks if a character is a valid hexadecimal character (0-9, a-f, A-F).
func isHexadecimal(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
