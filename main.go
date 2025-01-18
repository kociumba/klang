package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/kociumba/klang/codegen"
	"github.com/kociumba/klang/lexer"
	"github.com/kociumba/klang/parser"
)

func main() {
	input, err := os.ReadFile("grammar.k")
	if err != nil {
		panic(err)
	}

	// Create lexer
	l := lexer.New(string(input))

	if err := l.Preprocess(); err != nil {
		log.Printf("Preprocessor error: %s", err)
		return
	}

	// Create parser
	p := parser.New(l)

	// Parse program
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			log.Printf("Parser error: %s", err)
		}
		return
	}

	// Generate C code
	gen := codegen.New()
	output := gen.Generate(program)

	// Write to file
	err = os.WriteFile("output.c", []byte(output), 0644)
	if err != nil {
		panic(err)
	}
}
