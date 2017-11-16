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
	tracks  []string
}

// String printer for Patterns
func (p *Pattern) String() string {
	printed := fmt.Sprintf("Saved with HW Version: %v\n", p.version)
	printed += fmt.Sprintf("Tempo: %v\n", p.tempo)
	for _, tr := range p.tracks {
		printed += tr
	}
	return printed
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

	// Collect tracks info, exactly how depends on the file version
	tracksInfo := contents[50:]
	if version == "0.708-alpha" {
		p.tracks = parseTracksWithStopHeader(tracksInfo)
	} else {
		p.tracks = parseTracks(tracksInfo)
	}

	return &p, nil
}

// parseTracks reads track info (id, instrument, beat pattern) and formats it into strings like:
// (0) cowbell    |-x-x|-x-x|-x-x|-x-x|
func parseTracks(t []byte) []string {
	tracks := make([]string, 0)

	for len(t) > 0 {

		var trackDisplay = "("

		trackID := t[0]
		trackDisplay += fmt.Sprintf("%v) ", trackID)

		instrumentNameLength := t[4]
		instrumentNameEnd := 5 + instrumentNameLength
		instrumentName := string(t[5:instrumentNameEnd])
		trackDisplay += instrumentName
		trackDisplay += "\t|"

		trackPattern := t[instrumentNameEnd:(instrumentNameEnd + 16)]
		for i, x := range trackPattern {
			if x == 1 {
				trackDisplay += "x"
			} else if x == 0 {
				trackDisplay += "-"
			}
			if (i+1)%4 == 0 {
				trackDisplay += "|"
			}
		}
		trackDisplay += "\n"
		tracks = append(tracks, trackDisplay)
		t = t[(instrumentNameEnd + 16):]
	}

	return tracks
}

// parseTrackWithStopHeader does the same work as parseTracks, but for an earlier encoding version.
func parseTracksWithStopHeader(t []byte) []string {
	tracks := make([]string, 0)

	for len(t) > 0 {

		var trackDisplay = "("

		trackID := t[0]
		trackDisplay += fmt.Sprintf("%v) ", trackID)

		instrumentNameLength := t[4]
		instrumentNameEnd := 5 + instrumentNameLength
		instrumentName := string(t[5:instrumentNameEnd])
		trackDisplay += instrumentName
		trackDisplay += "\t|"

		trackPattern := t[instrumentNameEnd:(instrumentNameEnd + 16)]
		for i, x := range trackPattern {
			if x == 1 {
				trackDisplay += "x"
			} else if x == 0 {
				trackDisplay += "-"
			}
			if (i+1)%4 == 0 {
				trackDisplay += "|"
			}
		}
		trackDisplay += "\n"
		tracks = append(tracks, trackDisplay)

		potentialSpliceHeader := string(t[(instrumentNameEnd + 16):(instrumentNameEnd + 16 + 6)])
		if potentialSpliceHeader == "SPLICE" {
			t = make([]byte, 0)
		} else {
			t = t[(instrumentNameEnd + 16):]
		}
	}

	return tracks
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
