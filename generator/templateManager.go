package generator

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"
)

//go:embed templates/*.tmpl
var embeddedTemplates embed.FS

type TemplateManager struct {
	templates *template.Template
	buffers   map[string]*bytes.Buffer
}

func NewTemplateManager() (*TemplateManager, error) {
	// pp.Print(embeddedTemplates)
	tmpl, err := template.ParseFS(embeddedTemplates, "templates/*.go.tmpl")
	// pp.Print(tmpl.DefinedTemplates())
	if err != nil {
		return nil, err
	}

	return &TemplateManager{
		templates: tmpl,
		buffers: map[string]*bytes.Buffer{
			"headers":             new(bytes.Buffer),
			"typedefs":            new(bytes.Buffer),
			"structs":             new(bytes.Buffer),
			"variables":           new(bytes.Buffer),
			"macros":              new(bytes.Buffer),
			"function_prototypes": new(bytes.Buffer),
			"functions":           new(bytes.Buffer),
		},
	}, nil
}

func (tm *TemplateManager) GenerateToBuffer(templateName string, data interface{}, bufferKey string) error {
	buf, exists := tm.buffers[bufferKey]
	if !exists {
		return fmt.Errorf("buffer %q not found", bufferKey)
	}

	if err := tm.templates.ExecuteTemplate(buf, templateName, data); err != nil {
		return fmt.Errorf("error generating code for template %q: %w", templateName, err)
	}

	buf.WriteString("\n")
	return nil
}

func (tm *TemplateManager) WriteBuffersToBuilder(sb *strings.Builder) error {
	// Maintain a strict order of code sections
	order := []string{"headers", "typedefs", "structs", "variables", "macros", "function_prototypes", "functions"}
	for _, key := range order {
		_, err := sb.Write(tm.buffers[key].Bytes())
		if err != nil {
			return fmt.Errorf("error writing buffer %q to builder: %w", key, err)
		}
	}

	return nil
}
