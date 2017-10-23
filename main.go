package main

import (
	"bufio"
	"encoding/hex"
	"os"
	"path"
)

func main() {
	f, err := os.Open(path.Join("fixtures", "pattern_1.splice"))
	if err != nil {
		panic(err)
	}

	stdoutDumper := hex.Dumper(os.Stdout)
	defer stdoutDumper.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		stdoutDumper.Write([]byte(scanner.Text()))
	}

	f.Close()
}

// DecodeFile decodes the drum machine files.
func DecodeFile(path string) (string, error) {
	return "", nil
}
