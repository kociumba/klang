package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/k0kubun/pp/v3"
	"github.com/kociumba/klang/parser"
	"github.com/leaanthony/clir"
	"github.com/leaanthony/spinner"
)

var (
	mainSRC string
)

func main() {
	cli := clir.NewCli("klang", "compiler for the klang language", "v0.0.1")
	cli.DefaultCommand() // use this to run the build command by default, or run

	printAST := false
	cli.BoolFlag("ast", "print the abstract syntax tree", &printAST)
	printCCOut := false
	cli.BoolFlag("ccout", "print the output if the underlaying zig cc compiler", &printCCOut)

	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}

	if !validatePath(strings.Join(cli.OtherArgs(), "")) {
		os.Exit(1)
	}

	fmt.Printf("Main klang source file found at: %s", mainSRC)
	spin := spinner.New("Compiling")
	spin.Start()

	input, err := os.Open(mainSRC)
	if err != nil {
		panic(err)
	}

	replacedSRC, err := parser.GetReplacements(input)
	if err != nil {
		spin.Error("Failed to parse replacements: " + err.Error())
		// log.Fatal(err)
	}

	program := parser.Parse(replacedSRC, input.Name())

	if printAST {
		pp.Print(program)
	}

	// log.Infof("%+v", program)

	// cgen := generator.NewCodeGen().Generate(program)

	// log.Infof("%s", cgen)

	// os.WriteFile("test/output.c", []byte(cgen), 0644)

	cmd := exec.Command("zig", "cc", "test/output.c", "-o", "test/build/test.exe", "-O3")
	if printCCOut {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		spin.Error("Compilation error: " + err.Error())
		// log.Errorf("Compilation error: %s", err)
		os.Exit(0)
	}

	spin.Success("Compilation successful!")
}

func validatePath(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Errorf("Failed to resolve path: %v", err)
		return false
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Errorf("Path does not exist: %v", absPath)
		} else {
			log.Errorf("Error checking path: %v", err)
		}
		return false
	}

	if !info.IsDir() {
		if strings.HasSuffix(strings.ToLower(info.Name()), ".k") {
			mainSRC = absPath
			return true
		}
		log.Errorf("Path points to a file, but it's not a klang source file: %v\n\nDid you mean to point to a directory?", absPath)
		return false
	}

	hasKFile := false
	err = filepath.WalkDir(absPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".k") {
			if !hasKFile {
				mainSRC = filepath.Join(absPath, d.Name())
			}
			hasKFile = true
			return nil
		}
		return nil
	})

	if err != nil {
		log.Errorf("Error searching for klang source files: %v", err)
		return false
	}

	if !hasKFile {
		log.Errorf("Directory does not contain any klang source files: %v", absPath)
		return false
	}

	return true
}
