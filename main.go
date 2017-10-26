package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

// Pattern has the details of decoded drum machine patterns.
type Pattern struct {
	version string
	tempo   float32
}

// Pattern printer template
func (p *Pattern) String() string {
	return fmt.Sprintf(`Saved with HW Version: %s
Tempo: %v
`, p.version, p.tempo)
}

func main() {
}

// hexdump scans the given file and dumps the contents to Stdout.
func hexdump(fpath string) {
	f, err := os.Open(fpath)
	defer f.Close()
	check(err)

	stdoutDumper := hex.Dumper(os.Stdout)
	defer stdoutDumper.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		stdoutDumper.Write(scanner.Bytes())
	}
}

// DecodeFile decodes the drum machine files.
func DecodeFile(fpath string) (*Pattern, error) {
	p := Pattern{}
	f, err := os.Open(fpath)
	defer f.Close()
	check(err)

	contents, err := ioutil.ReadFile(fpath)
	check(err)

	// Check for the SPLICE header that starts each drum
	if headerInfo := string(contents[:6]); headerInfo != "SPLICE" {
		return &p, errors.New("SPLICE header not found")
	}

	// Find version and get as string
	version := strings.Trim(string(contents[14:27]), "\x00")
	p.version = version

	// Get tempo and decode it
	tempo := math.Float32frombits(binary.LittleEndian.Uint32(contents[46:50]))
	p.tempo = tempo

	// TODO: decode tracks info

	return &p, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
