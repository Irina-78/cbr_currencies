package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type CbrDecoder struct {
	xml.Decoder
}

func newCbrDecoder(inStr string) *CbrDecoder {
	inStr = strings.Replace(inStr, ",", ".", -1)

	decoder := xml.NewDecoder(strings.NewReader(inStr))

	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}

	return &CbrDecoder{*decoder}
}
