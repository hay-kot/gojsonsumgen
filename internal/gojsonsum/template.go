package gojsonsum

import (
	_ "embed"
	"html/template"
	"strings"
)

//go:embed template.gotmpl
var CodeTemplate string
var t = template.Must(template.New("code").Parse(CodeTemplate))

// Render renders the code template
func Render(pkg string, args []SumTypeDef) (string, error) {
	builder := &strings.Builder{}

	err := t.Execute(builder, map[string]interface{}{
		"PackageName":  pkg,
		"SumTypes": args,
	})
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}
