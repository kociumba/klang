package generator

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/kociumba/klang/parser"
)

//go:embed templates/*.tmpl
var embeddedTemplates embed.FS

type TemplateManager struct {
	templates *template.Template
	buffers   map[string]*bytes.Buffer
}

func NewTemplateManager() (*TemplateManager, error) {
	tmpl := template.New("")

	funcs := template.FuncMap{
		"reverseMods": func(mods []parser.Modifier) []parser.Modifier {
			reversed := make([]parser.Modifier, len(mods))
			copy(reversed, mods)
			for i := 0; i < len(reversed)/2; i++ {
				j := len(reversed) - i - 1
				reversed[i], reversed[j] = reversed[j], reversed[i]
			}
			return reversed
		},
		"renderSize": func(size *parser.Expression) string {
			if size == nil {
				return ""
			}
			buf := new(bytes.Buffer)
			if err := tmpl.ExecuteTemplate(buf, "expression", size); err != nil {
				return "<ERROR>"
			}
			return buf.String()
		},
		"needsParens": func(mods []parser.Modifier, idx int) bool {
			return idx > 0 && mods[idx-1].Array != nil
		},
		"repeat": func(count int, str string) string {
			return strings.Repeat(str, count)
		},
	}

	// pp.Print(embeddedTemplates)
	tmpl.Funcs(funcs)
	tmpl = template.Must(tmpl.ParseFS(embeddedTemplates, "templates/*.go.tmpl"))
	// pp.Print(tmpl.DefinedTemplates())

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
