package gomatimport

import (
	"encoding/binary"
	"io"
)

type Header struct {
	Description string
	SubsysData  interface{}
	Version     uint16
	Endian      string
}

func DecodeHeader(f io.Reader) Header {
	hdbuffer := make([]byte, HeaderLength)
	f.Read(hdbuffer)

	var out Header

	out.Description = string(hdbuffer[:DescriptionLength])
	// out.SubsysData = ?
	out.Version = binary.LittleEndian.Uint16(hdbuffer[124:126])
	out.Endian = string(hdbuffer[126:128])

	return out
}
