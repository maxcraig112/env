package env_test

import (
	"errors"
	"testing"
	"time"

	"github.com/maxcraig112/env"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name     string
		envValue *string // nil means unset
		defaults []string
		want     string
	}{
		{
			name:     "set value returned regardless of defaults",
			envValue: strPtr("hello"),
			defaults: []string{"fallback"},
			want:     "hello",
		},
		{
			name:     "unset with single default",
			envValue: nil,
			defaults: []string{"fallback"},
			want:     "fallback",
		},
		{
			name:     "unset with multiple defaults uses first",
			envValue: nil,
			defaults: []string{"first", "second"},
			want:     "first",
		},
		{
			name:     "unset with no defaults returns empty string",
			envValue: nil,
			defaults: nil,
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const key = "ENV_TEST_GET"
			if tt.envValue != nil {
				t.Setenv(key, *tt.envValue)
			}

			got := env.Get(key, tt.defaults...)
			if got != tt.want {
				t.Errorf("Get() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRequire(t *testing.T) {
	tests := []struct {
		name       string
		envValue   *string
		want       string
		checkError func(*testing.T, error)
	}{
		{
			name:       "set value returned with no error",
			envValue:   strPtr("value"),
			want:       "value",
			checkError: requireNoError,
		},
		{
			name:       "unset returns NotSetError",
			checkError: requireNotSetError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const key = "ENV_TEST_REQUIRE"
			if tt.envValue != nil {
				t.Setenv(key, *tt.envValue)
			}

			got, err := env.Require(key)
			tt.checkError(t, err)
			if err == nil && got != tt.want {
				t.Errorf("Require() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	tests := []struct {
		name       string
		envValue   *string
		defaults   []int
		want       int
		checkError func(*testing.T, error)
	}{
		{
			name:       "set value parsed",
			envValue:   strPtr("7"),
			want:       7,
			checkError: requireNoError,
		},
		{
			name:       "unset with default",
			defaults:   []int{42},
			want:       42,
			checkError: requireNoError,
		},
		{
			name:       "unset with no default returns zero value",
			want:       0,
			checkError: requireNoError,
		},
		{
			name:       "malformed value returns ParseError",
			envValue:   strPtr("not-a-number"),
			checkError: requireParseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const key = "ENV_TEST_INT"
			if tt.envValue != nil {
				t.Setenv(key, *tt.envValue)
			}

			got, err := env.GetInt(key, tt.defaults...)
			tt.checkError(t, err)
			if err == nil && got != tt.want {
				t.Errorf("GetInt() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestRequireInt(t *testing.T) {
	tests := []struct {
		name       string
		envValue   *string
		want       int
		checkError func(*testing.T, error)
	}{
		{
			name:       "set value parsed",
			envValue:   strPtr("99"),
			want:       99,
			checkError: requireNoError,
		},
		{
			name:       "unset returns NotSetError",
			checkError: requireNotSetError,
		},
		{
			name:       "malformed value returns ParseError",
			envValue:   strPtr("nope"),
			checkError: requireParseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const key = "ENV_TEST_REQUIRE_INT"
			if tt.envValue != nil {
				t.Setenv(key, *tt.envValue)
			}

			got, err := env.RequireInt(key)
			tt.checkError(t, err)
			if err == nil && got != tt.want {
				t.Errorf("RequireInt() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestGetBool(t *testing.T) {
	tests := []struct {
		name       string
		envValue   *string
		defaults   []bool
		want       bool
		checkError func(*testing.T, error)
	}{
		{name: "set true", envValue: strPtr("true"), want: true, checkError: requireNoError},
		{name: "set false", envValue: strPtr("false"), want: false, checkError: requireNoError},
		{name: "unset with default", defaults: []bool{true}, want: true, checkError: requireNoError},
		{name: "malformed value", envValue: strPtr("nope"), checkError: requireParseError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const key = "ENV_TEST_BOOL"
			if tt.envValue != nil {
				t.Setenv(key, *tt.envValue)
			}

			got, err := env.GetBool(key, tt.defaults...)
			tt.checkError(t, err)
			if err == nil && got != tt.want {
				t.Errorf("GetBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDuration(t *testing.T) {
	tests := []struct {
		name       string
		envValue   *string
		defaults   []time.Duration
		want       time.Duration
		checkError func(*testing.T, error)
	}{
		{
			name:       "set value parsed",
			envValue:   strPtr("1h30m"),
			want:       90 * time.Minute,
			checkError: requireNoError,
		},
		{
			name:       "unset with default",
			defaults:   []time.Duration{5 * time.Second},
			want:       5 * time.Second,
			checkError: requireNoError,
		},
		{
			name:       "malformed value",
			envValue:   strPtr("not-a-duration"),
			checkError: requireParseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const key = "ENV_TEST_DURATION"
			if tt.envValue != nil {
				t.Setenv(key, *tt.envValue)
			}

			got, err := env.GetDuration(key, tt.defaults...)
			tt.checkError(t, err)
			if err == nil && got != tt.want {
				t.Errorf("GetDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func strPtr(s string) *string { return &s }

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireNotSetError(t *testing.T, err error) {
	t.Helper()
	var notSet *env.NotSetError
	if !errors.As(err, &notSet) {
		t.Fatalf("err = %v, want *env.NotSetError", err)
	}
}

func requireParseError(t *testing.T, err error) {
	t.Helper()
	var parseErr *env.ParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("err = %v, want *env.ParseError", err)
	}
}
