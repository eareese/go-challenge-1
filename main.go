package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"path"
)

// Pattern has the details of decoded drum machine patterns.
type Pattern struct {
	version string
}

func (p *Pattern) String() string {
	return fmt.Sprintf("Saved with HW Version: %s\n", p.version)
}

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
func DecodeFile(path string) (*Pattern, error) {
	return &Pattern{}, nil
}
