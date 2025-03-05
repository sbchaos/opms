package compiler

import (
	"bytes"
	"errors"
	"strings"
	"text/template"
	"time"
)

const (
	// ISODateFormat https://en.wikipedia.org/wiki/ISO_8601
	ISODateFormat = "2006-01-02"

	ISOTimeFormat = time.RFC3339
)

// Engine compiles a set of defined macros using the provided context
type Engine struct {
	baseTemplate *template.Template
}

func NewEngine() *Engine {
	baseTemplate := template.
		New("opms_template_engine").
		Funcs(OptimusFuncMap())

	return &Engine{
		baseTemplate: baseTemplate,
	}
}

func (e *Engine) CompileString(input string, context map[string]any) (string, error) {
	tmpl, err := e.baseTemplate.New("base").Parse(input)
	if err != nil {
		return "", errors.New("unable to parse string " + input)
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, context); err != nil {
		return "", errors.New("unable to render string " + input)
	}
	return strings.TrimSpace(buf.String()), nil
}
