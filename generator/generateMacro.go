package generator

import (
	"fmt"
	"strings"

	"github.com/kociumba/klang/parser"
)

func (cg *CodeGen) generateMacro(fn *parser.FunctionDef) string {
	if fn == nil {
		return ""
	}

	var output strings.Builder

	// First generate the type checking macro helpers
	output.WriteString(fmt.Sprintf("// Type checking helper for %s\n", fn.Name))
	output.WriteString(fmt.Sprintf("#define %s_TYPE_CHECK(", fn.Name))

	// Add parameter names to type check macro
	paramNames := make([]string, 0, len(fn.Params))
	for i, param := range fn.Params {
		if i > 0 {
			output.WriteString(", ")
		}
		paramNames = append(paramNames, param.Name)
		output.WriteString(param.Name)
	}
	output.WriteString(") \\\n")

	// Add type checks for each parameter
	for i, param := range fn.Params {
		if i > 0 {
			output.WriteString(" && ")
		}
		output.WriteString(fmt.Sprintf("_Generic((%s), %s: 1, default: 0)",
			param.Name, param.Type))
	}
	output.WriteString("\n\n")

	// Generate the static inline implementation
	output.WriteString("static inline ")
	if fn.ReturnType != nil {
		output.WriteString(*fn.ReturnType)
	} else {
		output.WriteString("void")
	}

	output.WriteString(fmt.Sprintf(" %s_impl(", fn.Name))

	// Parameters with types
	args := make([]string, 0, len(fn.Params))
	for _, param := range fn.Params {
		args = append(args, fmt.Sprintf("%s %s", param.Type, param.Name))
	}
	output.WriteString(strings.Join(args, ", "))
	output.WriteString(") {\n")

	// Generate function body
	if fn.Body != nil {
		output.WriteString(cg.generateBlock(fn.Body))
	}
	output.WriteString("}\n\n")

	// Create the final macro that includes type checking
	output.WriteString(fmt.Sprintf("#define %s(", fn.Name))
	output.WriteString(strings.Join(paramNames, ", "))
	output.WriteString(") \\\n")
	output.WriteString(fmt.Sprintf("    (_Static_assert(%s_TYPE_CHECK(", fn.Name))
	output.WriteString(strings.Join(paramNames, ", "))
	output.WriteString(fmt.Sprintf("), \"Type mismatch in %s\"), \\\n", fn.Name))
	output.WriteString(fmt.Sprintf("     %s_impl(", fn.Name))
	output.WriteString(strings.Join(paramNames, ", "))
	output.WriteString("))\n\n")

	return output.String()
}
