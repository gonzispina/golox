package main

import (
	"fmt"
	"golox/lox"
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

	}
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
	s := lox.NewScanner(b)

	tokens, err := s.ScanTokens()
	if err != nil {
		return err
	}

	for _, t := range tokens {
		fmt.Printf("Token: %s \n", t)
	}

	return nil
}
