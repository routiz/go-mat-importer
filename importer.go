package gomatimport

import (
	"encoding/binary"
	"fmt"
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
