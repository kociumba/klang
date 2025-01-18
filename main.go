package main

import (
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
	"github.com/kociumba/klang/generator"
	"github.com/kociumba/klang/parser"
)

func main() {
	input, err := os.Open("test/grammar.k")
	if err != nil {
		panic(err)
	}

	program := parser.Parse(input)

	log.Infof("%+v", program)

	cgen := generator.NewCodeGen().Generate(program)

	log.Infof("%s", cgen)

	os.WriteFile("test/output.c", []byte(cgen), 0644)

	cmd := exec.Command("zig", "cc", "test/output.c", "-o", "test/build/test.exe", "-O3")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Compilation error: %s", err)
	}
}

// func main() {
// 	input, err := os.ReadFile("test/grammar.k")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Create lexer
// 	l := lexer.New(string(input))

// 	if err := l.Preprocess(); err != nil {
// 		log.Errorf("Preprocessor error: %s", err)
// 		return
// 	}

// 	// Create parser
// 	p := parser.New(l)

// 	// Parse program
// 	program := p.ParseProgram()
// 	if len(p.Errors()) > 0 {
// 		for _, err := range p.Errors() {
// 			log.Errorf("Parser error: %s", err)
// 		}
// 		return
// 	}

// 	// Generate C code
// 	gen := codegen.New()
// 	output := gen.Generate(program)

// 	// Write to file
// 	err = os.WriteFile("test/output.c", []byte(output), 0644)
// 	if err != nil {
// 		panic(err)
// 	}

// 	cmd := exec.Command("zig", "cc", "test/output.c")
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	if err := cmd.Run(); err != nil {
// 		log.Errorf("Compilation error: %s", err)
// 	}
// }
