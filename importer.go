package gomatimport

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const (
	HeaderLength      = 128
	DescriptionLength = 116
)

type Header struct {
	Description string
	SubsysData  interface{}
	Version     uint16
	Endian      string
}

type DataElement struct {
	Type          uint
	NumberOfBytes uint
	RawBytes      []byte
}

type Mat struct {
	H    Header
	Data interface{}
}

func DecodeHeader(f *os.File) Header {
	hdbuffer := make([]byte, HeaderLength)
	f.Read(hdbuffer)

	var out Header

	out.Description = string(hdbuffer[:DescriptionLength])
	// out.SubsysData = ?
	out.Version = binary.LittleEndian.Uint16(hdbuffer[124:126])
	out.Endian = string(hdbuffer[126:128])

	return out
}

func SplitData(f *os.File) []DataElement {
	delements := make([]DataElement, 0, 10)
	for {
		var delm DataElement

		// Read first 4 bytes and check if this element is small.
		tgbuffer1 := [4]byte{}
		readcnt, err := f.Read(tgbuffer1[:])
		if err != nil && err != io.EOF {
			panic(err)
		}
		if readcnt == 0 || err == io.EOF {
			break
		}
		// Check if first 2 bytes of the tag are zeroes.
		typ := binary.LittleEndian.Uint16(tgbuffer1[2:])
		// Number of bytes
		var nob uint
		if typ != 0 {
			// If first 2 bytes are not zeros, this is a small data element.
			nob = uint(binary.LittleEndian.Uint16(tgbuffer1[2:]))
		} else {
			// If first 2 bytes are zeros, this is a normal data element
			// which is not a small data element.

			// First 4 bytes are type.
			typ = binary.LittleEndian.Uint16(tgbuffer1[:])

			tgbuffer2 := [4]byte{}
			readcnt, err = f.Read(tgbuffer2[:])
			if err != nil && err != io.EOF {
				panic(err)
			}
			if readcnt == 0 || err == io.EOF {
				break
			}
			nob = uint(binary.LittleEndian.Uint32(tgbuffer2[:]))
		}
		dataBuffer := make([]byte, nob)
		readcnt, err = f.Read(dataBuffer) // maybe removed

		delm.Type = uint(typ)
		delm.NumberOfBytes = uint(nob)
		delm.RawBytes = dataBuffer

		delements = append(delements, delm)

		fmt.Println("type = ", typ)
		fmt.Println("number of bytes = ", nob)
		fmt.Println("read count = ", readcnt)
		// do something
	}

	return delements
}

func DecodeDataElement(delm DataElement) interface{} {
	return nil
}

func DecodeData(f *os.File) interface{} {
	dataElements := SplitData(f)
	for _, delm := range dataElements {
		DecodeDataElement(delm)
	}
	var out Mat
	return out
}

func Import(filename string, dst interface{}) (Mat, error) {
	var out Mat

	f, err := os.Open(filename)
	if err != nil {
		return out, fmt.Errorf("Failed to open file %s, error message is : %s\n", err.Error())
	}
	defer func() {
		f.Close()
	}()

	// Read Header
	out.H = DecodeHeader(f)
	out.Data = DecodeData(f)

	return out, nil
}
