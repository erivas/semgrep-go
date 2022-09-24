package main

import (
	"errors"
	"io/ioutil"
	"os"

	"semgrep-go/internal/semgrep"
)

func createTempFile(v string) (string, error) {
	bytes := []byte(v)

	file, err := ioutil.TempFile("", "semgrep_input")
	if err != nil {
		return "", err
	}

	n, err := file.Write(bytes)
	if err != nil {
		return file.Name(), nil
	}
	if n != len(bytes) {
		return file.Name(), errors.New("expected to write more bytes")
	}

	err = file.Close()
	if err != nil {
		return file.Name(), err
	}

	return file.Name(), nil
}

func main() {
	println("Calling semgrep...")
	fn, _ := createTempFile("let x = 42")
	defer os.Remove(fn)
	println(semgrep.Run("ocaml", "rule.yaml", fn))
}
