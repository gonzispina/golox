package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"golox/lox"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lox [path to script]")
	} else if len(os.Args) == 2 {
		path, err := filepath.Abs(os.Args[1])
		if err != nil {
			fmt.Printf("Unable to find path %s", path)
			os.Exit(1)
		}
		err = runFile(path)
		if err != nil {
			os.Exit(2)
		}
	} else {
		err := runPrompt()
		if err != nil {
			os.Exit(3)
		}
	}
}

func runPrompt() error {
	fmt.Println("Welcome to lox command prompt!")
	defer func() {
		fmt.Println("Goodbye!")
	}()

	reader, err := readline.New("lox > ")

	if err != nil {
		return err
	}

	for {
		fmt.Print("lox > ")
		line, err := reader.Readline()
		if err != nil {
			if err == io.EOF || err == readline.ErrInterrupt {
				break
			}

			return err
		}

		fmt.Print(line)
	}

	return nil
}

func runFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return run(b)
}

func run(b []byte) error {
	tokens, err := lox.NewScanner(string(b)).ScanTokens()
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		return nil
	}

	e, err := lox.NewParser(tokens).Parse()
	if err != nil {
		return err
	}

	value, err := lox.NewInterpreter().Interpret(e)
	if err != nil {
		return err
	}

	fmt.Println(value)

	return nil
}