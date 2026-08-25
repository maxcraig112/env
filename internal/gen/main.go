// Command gen generates the typed Get*/Require* functions for the env
// package (zz_generated.go). Run it via `go generate ./...` (or `make
// generate`) from the module root after changing the type list below or
// template.go.tmpl.
package main

import (
	"bytes"
	_ "embed"
	"go/format"
	"log"
	"os"
	"text/template"
)

//go:embed template.go.tmpl
var tmplText string

type typeSpec struct {
	// Suffix is appended to Get/Require, e.g. "Int" -> GetInt, RequireInt.
	Suffix string
	// GoType is the Go type returned by the generated functions.
	GoType string
	// Parse is a Go expression, using the variable "raw", that evaluates to
	// (GoType, error).
	Parse string
}

var types = []typeSpec{
	{Suffix: "Int", GoType: "int", Parse: "strconv.Atoi(raw)"},
	{Suffix: "Int64", GoType: "int64", Parse: "strconv.ParseInt(raw, 10, 64)"},
	{Suffix: "Float64", GoType: "float64", Parse: "strconv.ParseFloat(raw, 64)"},
	{Suffix: "Bool", GoType: "bool", Parse: "strconv.ParseBool(raw)"},
	{Suffix: "Duration", GoType: "time.Duration", Parse: "time.ParseDuration(raw)"},
}

func main() {
	tmpl := template.Must(template.New("gen").Parse(tmplText))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, types); err != nil {
		log.Fatalf("gen: executing template: %v", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("gen: formatting output: %v", err)
	}

	if err := os.WriteFile("zz_generated.go", formatted, 0o644); err != nil {
		log.Fatalf("gen: writing zz_generated.go: %v", err)
	}
}
