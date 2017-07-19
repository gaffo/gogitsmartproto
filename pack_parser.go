package gogitsmartproto

import (
	"io"
	"errors"
	"fmt"
	"encoding/binary"
)

type Pack struct {
	CItems int
	Version int
}

func ParsePack(rdr io.Reader) (pack Pack, err error) {
	err = nil
	four := make([]byte, 4, 4)
	rdr.Read(four)
	header := string(four)
	fmt.Println(four)
	fmt.Println("hi [", header, "]")
	if header != "PACK" {
		err = errors.New("Invalid Header")
		return
	}

	rdr.Read(four)
	pack.Version = int(binary.BigEndian.Uint32(four))

	rdr.Read(four)
	pack.CItems = int(binary.BigEndian.Uint32(four))

	

	return
}
