package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
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
		stdoutDumper.Write([]byte(scanner.Text()))
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

	// spliceInfo := string(contents[:6])
	version := strings.Trim(string(contents[14:27]), "\x00")
	p.version = version

	tempo := math.Float32frombits(binary.LittleEndian.Uint32(contents[46:50]))
	p.tempo = tempo

	return &p, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
