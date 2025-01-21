package parser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Replacement struct {
	Rule *ReplaceRule `"replace" @@`
}

type ReplaceRule struct {
	From string `@Ident`
	To   string `"->" @(String|Ident)`
}

type Replacements struct {
	Rules []Replacement `(@@)*`
}

func extractReplacements(input io.Reader) (string, error) {
	var replacements []string

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "replace") {
			replacements = append(replacements, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	return strings.Join(replacements, "\n"), nil
}

// need to to modify the src code in memory to replace the replaced keywords
// since go tags are static and can't use variables, not a big problem tho
// just some simple string replacement
func GetReplacements(input *os.File) (string, error) {
	filtered, err := extractReplacements(input)
	if err != nil {
		return "", err
	}

	parser := participle.MustBuild[Replacements](
		participle.Lexer(lexer.MustSimple([]lexer.SimpleRule{
			{Name: "Ident", Pattern: `[a-zA-Z_]\w*`},
			{Name: "String", Pattern: `"[^"]*"`},
			{Name: "Arrow", Pattern: `->`},
			{Name: "Comment", Pattern: `//[^\n]*`},
			{Name: "Whitespace", Pattern: `\s+`},
			{Name: "EOL", Pattern: `[\r\n]+`},
		})),
		participle.Elide("Whitespace", "Comment", "EOL"),
	)

	replacements := Replacements{}
	r, err := parser.ParseString(input.Name(), filtered)
	if err != nil {
		return "", fmt.Errorf("failed to parse replacements: %w", err)
	}

	replacements = *r

	replaceMap := make(map[string]string)
	for _, rule := range replacements.Rules {
		replaceMap[rule.Rule.To] = rule.Rule.From
	}

	if _, err := input.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("failed to reset file cursor: %w", err)
	}

	data, err := io.ReadAll(input)
	if err != nil {
		return "", fmt.Errorf("failed to read input file: %w", err)
	}

	content := string(data)
	for userReplacement, originalKeyword := range replaceMap {
		content = strings.ReplaceAll(content, userReplacement, originalKeyword)
	}

	return content, nil
}
