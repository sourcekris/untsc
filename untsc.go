package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/JoshVarga/blast"
)

const (
	name       = "TSComp"
	ext        = "TSC"
	fHeaderLen = 13 // fileID + 7 unknown bytes
)

var (
	fileID  = []byte{0x65, 0x5D, 0x13, 0x8C, 0x08, 0x01}
	fset    = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	arcFile = fset.String("e", "", fmt.Sprintf("The %s file to extract.", ext))
	dstPath = fset.String("d", "", "Optional output directory to extract to.")
)

func errpanic(e error) {
	if e != nil {
		panic(e)
	}
}

type header struct {
	id     []byte // Stores the file sig
	cSize  int    // Compressed file size
	fnSize int    // Filename length
	fSize  int    // Total file size
	fn     string // Filename for this header.
}

func main() {
	fset.Parse(os.Args[1:])

	if *arcFile == "" {
		fset.Usage()
		os.Exit(0)
	}

	if *dstPath != "" {
		if _, err := os.Stat(*dstPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "destination folder (%s) doesn't exist: %v", *dstPath, err)
			os.Exit(1)
		}
	}

	f, err := os.Open(*arcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening %s: %v", *arcFile, err)
		os.Exit(1)
	}
	defer f.Close()

	s, err := f.Stat()
	errpanic(err)

	h := &header{
		id:    make([]byte, len(fileID)),
		fSize: int(s.Size()),
	}

	// Read the file header.
	_, err = f.Read(h.id)
	errpanic(err)

	if !reflect.DeepEqual(h.id, fileID) {
		fmt.Fprintf(os.Stderr, "file is not a %s file", ext)
		os.Exit(1)
	}

	// See to the first entry header.
	_, err = f.Seek(0xd, io.SeekStart)
	errpanic(err)

	// Read archive members in a loop until done.
	for {
		// Read metadate from file.
		metadata := make([]byte, 16)
		_, err = f.Read(metadata)
		errpanic(err)
		h.cSize = int(binary.LittleEndian.Uint32(metadata[1:5]))
		h.fnSize = int(metadata[15])

		if h.fnSize > 12 {
			fmt.Fprintf(os.Stderr, "filename length is > 12: %d", h.fnSize)
			os.Exit(1)
		}

		fn := make([]byte, h.fnSize)
		_, err = f.Read(fn)
		errpanic(err)
		h.fn = string(fn)

		fmt.Printf("Extracting: %s (%d compressed bytes)\n", h.fn, h.cSize)

		_, err = f.Seek(1, io.SeekCurrent) // Seek passed the null string terminator.
		errpanic(err)
		// Read all the DCL compressed data.
		dcl := make([]byte, h.cSize)
		_, err = f.Read(dcl)
		errpanic(err)

		b := bytes.NewReader(dcl)
		r, err := blast.NewReader(b)
		errpanic(err)

		if *dstPath != "" {
			h.fn = *dstPath + "/" + h.fn
		}

		o, err := os.Create(h.fn)
		errpanic(err)

		_, err = io.Copy(o, r)
		errpanic(err)
		o.Close()

		cur, err := f.Seek(0, io.SeekCurrent)
		errpanic(err)
		if cur == int64(h.fSize) {
			break
		}
	}
}
