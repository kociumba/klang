package main

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/k0kubun/pp/v3"
	"github.com/kociumba/klang/generator"
	"github.com/kociumba/klang/parser"
	"github.com/leaanthony/clir"
	"github.com/leaanthony/spinner"
)

type Actions struct {
	Build bool
	Run   bool
}

var (
	mainSRC string
	action  Actions
)

func main() {
	cli := clir.NewCli("klang", "compiler for the klang language", "v0.0.1")

	printAST := false
	cli.BoolFlag("ast", "print the abstract syntax tree", &printAST)
	printCCOut := false
	cli.BoolFlag("ccout", "print the output of the underlaying zig cc compiler", &printCCOut)
	src := "no path provided"
	cli.StringFlag("src", "main source file to compile", &src)

	build := cli.NewSubCommandInheritFlags("build", "build a klang source file to an executable")
	build.Action(func() error {
		action.Build = true
		return nil
	})
	run := cli.NewSubCommandInheritFlags("run", "run a klang source file")
	run.Action(func() error {
		action.Run = true
		return nil
	})

	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}

	spin := spinner.New("Finding source files")
	spin.Start()

	if src == "no path provided" {
		spin.Error("No source file provided\n\nHINT: Use the -src flag")
		os.Exit(0)
	}

	validatePath(src)

	// fmt.Printf("Main klang source file found at: %s", mainSRC)

	spin.UpdateMessage("Reading source files")

	mainSRC = filepath.Clean(mainSRC)

	input, err := os.Open(mainSRC)
	if err != nil {
		spin.Error("Failed to open source file: " + err.Error())
		os.Exit(0)
	}

	spin.UpdateMessage("Applying replacements")

	replacedSRC, err := parser.GetReplacements(input)
	if err != nil {
		spin.Error("Failed to parse replacements: " + err.Error())
		os.Exit(0)
	}

	spin.UpdateMessage("Building AST")

	program := parser.Parse(replacedSRC, input.Name())

	if printAST {
		pp.Print(program)
	}

	spin.UpdateMessage("Generating C intermediary")

	cgen := generator.NewCodeGen().Generate(program)

	intermediary := strings.TrimSuffix(input.Name(), ".k") + ".c"

	os.WriteFile(intermediary, []byte(cgen), 0644)

	spin.UpdateMessage("Compiling C intermediary")

	var suffix string
	switch runtime.GOOS {
	case "windows":
		suffix = ".exe"
	default:
		suffix = ""
	}

	binary := strings.TrimSuffix(input.Name(), ".k") + suffix
	cmd := exec.Command("zig", "cc", intermediary, "-o", binary, "-O3")
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

	if action.Run {
		exec.Command(binary).Start()
	}
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
