package uuidcheck_test

import (
	"github.com/ashwingopalsamy/uuidcheck"
	"strings"
	"testing"
	"time"
)

func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect bool
	}{
		// Valid Cases
		{
			name:   "Lowercase valid UUID",
			input:  "f47ac10b-58cc-0372-8567-0e02b2c3d479",
			expect: true,
		},
		{
			name:   "Uppercase valid UUID",
			input:  "F47AC10B-58CC-0372-8567-0E02B2C3D479",
			expect: true,
		},
		{
			name:   "Mixed case valid UUID",
			input:  "f47Ac10B-58Cc-0372-8567-0E02b2C3D479",
			expect: true,
		},

		// Invalid Length
		{
			name:   "Empty string",
			input:  "",
			expect: false,
		},
		{
			name:   "Shorter than 36 chars",
			input:  "f47ac10b-58cc-0372-8567-0e02b2c3d47",
			expect: false,
		},
		{
			name:   "Longer than 36 chars",
			input:  "f47ac10b-58cc-0372-8567-0e02b2c3d479abc",
			expect: false,
		},

		// Invalid Hyphens
		{
			name:   "No hyphens at all",
			input:  "f47ac10b58cc037285670e02b2c3d479",
			expect: false,
		},
		{
			name:   "Hyphens in wrong places",
			input:  "f47ac10b-58cc0-372-8567-0e02b2c3d479",
			expect: false,
		},
		{
			name:   "Only hyphens",
			input:  "------------------------------------",
			expect: false,
		},

		// Invalid characters
		{
			name:   "Invalid character (g) in UUID",
			input:  "f47ac10b-58cc-0372-8567-0e02b2c3d47g",
			expect: false,
		},
		{
			name:   "Invalid character (Z) in UUID",
			input:  "Z47ac10b-58cc-0372-8567-0e02b2c3d479",
			expect: false,
		},

		// Edge Cases
		{
			name:   "All zeros but valid format",
			input:  "00000000-0000-0000-0000-000000000000",
			expect: true,
		},
		{
			name:   "All hyphens in correct positions but invalid",
			input:  "------------------------------------",
			expect: false,
		},
		{
			name:   "Almost valid but last char not hex",
			input:  "f47ac10b-58cc-0372-8567-0e02b2c3d47x",
			expect: false,
		},
		{
			name:   "Non-hex character in hex position",
			input:  "f47ac10b-58cc-0372-8567-0e02b2c3d47%",
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := uuidcheck.IsValidUUID(tt.input)
			if res != tt.expect {
				t.Errorf("IsValidUUID(%q) = %v; want %v", tt.input, res, tt.expect)
			}
		})
	}
}

func TestIsUUIDv7(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectV7    bool
		expectPanic bool
	}{
		{
			name:        "Short string",
			input:       "abcd",
			expectV7:    false,
			expectPanic: true, // accessing uuid[14] should panic
		},
		{
			name:        "Non-UUID string but long enough",
			input:       "abcdefghijklmnopqrstuvxyz0123456789abcd", // 36 chars but no hyphens
			expectV7:    false,
			expectPanic: false,
		},
		{
			name:        "Version 7 UUID all lowercase",
			input:       "00000000-0000-7000-0000-000000000000",
			expectV7:    true,
			expectPanic: false,
		},
		{
			name:        "Version 7 UUID mixed case",
			input:       "00000000-0000-7FFF-0000-000000000000",
			expectV7:    true,
			expectPanic: false,
		},
		{
			name:        "Version 4 UUID",
			input:       "f47ac10b-58cc-4372-8567-0e02b2c3d479",
			expectV7:    false,
			expectPanic: false,
		},
		{
			name:        "Version 1 UUID",
			input:       "f47ac10b-58cc-1372-8567-0e02b2c3d479",
			expectV7:    false,
			expectPanic: false,
		},
		{
			name:        "All zeros but version nibble not '7'",
			input:       "00000000-0000-4000-0000-000000000000", // version nibble = '4'
			expectV7:    false,
			expectPanic: false,
		},
		{
			name:        "Check upper nibble when version = '7'",
			input:       "00000000-0000-7abc-0000-000000000000",
			expectV7:    true,
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var got bool
			var didPanic bool

			func() {
				defer func() {
					if r := recover(); r != nil {
						didPanic = true
					}
				}()

				got = uuidcheck.IsUUIDv7(tt.input)
			}()

			if tt.expectPanic && !didPanic {
				t.Errorf("Expected panic but got none for input: %q", tt.input)
			}

			if !tt.expectPanic && didPanic {
				t.Errorf("Did not expect panic, but got one for input: %q", tt.input)
			}

			if !tt.expectPanic && got != tt.expectV7 {
				t.Errorf("IsUUIDv7(%q) = %v; want %v", tt.input, got, tt.expectV7)
			}
		})
	}
}

func TestUUIDv7ToTimestamp(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "Too short string",
			input:     "abcd",
			expectErr: true,
		},
		{
			name:      "Invalid character in first segment (time_low portion)",
			input:     "zzzzzzzz-0000-7000-0000-000000000000",
			expectErr: true,
		},
		{
			name:      "Invalid character in second segment",
			input:     "00000000-gggg-7000-0000-000000000000",
			expectErr: true,
		},
		{
			name:      "All zeros, valid length and hex - should return no error",
			input:     "00000000-0000-7000-0000-000000000000",
			expectErr: false,
		},
		{
			name:      "A valid random v7-like UUID with hex fields, no invalid chars",
			input:     "017f22e0-79b0-7cc0-98ac-2e7517f37f9f",
			expectErr: false,
		},
		{
			name:      "Maximal hex fields within allowed ranges",
			input:     "ffffffff-ffff-7fff-aaaa-ffffffffffff",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tm, err := uuidcheck.UUIDv7ToTimestamp(tt.input)
			if tt.expectErr && err == nil {
				t.Errorf("Expected an error but got none. time=%v", tm)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Did not expect an error, but got one: %v", err)
			}

			// If it's a valid case, we can make some sanity checks:
			if tt.input == "00000000-0000-7000-0000-000000000000" && err == nil {
				// According to the logic, this represents timestamp = 0
				// Which should map to Unix epoch start
				if !tm.Equal(time.UnixMilli(0).UTC()) {
					t.Errorf("Expected time to be Unix epoch start, got %v", tm)
				}
			}

			if strings.HasPrefix(tt.input, "017f22e0-79b0-7cc0") && err == nil {
				if tm.IsZero() {
					t.Error("Expected a non-zero time for this UUID")
				}
			}
		})
	}
}

func TestUUIDv7ToTime_Success(t *testing.T) {
	// Here, we take a known valid UUIDv7-like string.
	// In a real scenario, you'd generate a UUIDv7 using a proper UUIDv7 generator.
	uuid := "01939c67-06f5-7faf-ae43-6b450bff06af"

	tm, err := uuidcheck.UUIDv7ToTimestamp(uuid)
	if err != nil {
		t.Fatalf("Failed to convert UUID to time: %v", err)
	}

	// Log the time result. This demonstrates the conversion.
	t.Logf("UUID: %s converts to UTC time: %s", uuid, tm.Format(time.RFC3339Nano))

	// You might also add assertions if you know what timestamp to expect,
	// but here we just show that it produces a valid time.
	if tm.IsZero() {
		t.Error("Expected a valid non-zero time")
	}
}

// Additional note:
// If we had a known specification for what timestamp a certain UUIDv7 should produce,
// we could assert that exactly. For now, we're primarily testing error conditions and
// general correctness of parsing.
