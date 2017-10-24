package main

import (
	"bufio"
	"encoding/hex"
	"os"
	"path"
)

func main() {
	hexdump(path.Join("fixtures", "pattern_5.splice"))
}

// hexdump scans the given file and dumps the contents to Stdout.
func hexdump(fpath string) {
	f, err := os.Open(fpath)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	stdoutDumper := hex.Dumper(os.Stdout)
	defer stdoutDumper.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		stdoutDumper.Write([]byte(scanner.Text()))
	}
}

// DecodeFile decodes the drum machine files.
func DecodeFile(path string) (string, error) {
	return "", nil
}
