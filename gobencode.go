// Package gobencode allows to parse and unmarshall the bencode format with some utilies
package gobencode

import (
	"bufio"
	"io"

	"github.com/trixky/gobencode/bencode"
	"github.com/trixky/gobencode/parser"
)

// ParseFromReader parses the bencode format from reader in interface
func ParseFromReader(reader io.Reader) (data interface{}, err error) {
	bufioReader, ok := reader.(*bufio.Reader)

	if !ok {
		bufioReader = bufio.NewReader(reader)
	}

	data, err = parser.ParseElement(bufioReader)

	return
}

// UnmarshallFromReader parses and unmarshall the bencode format from reader in a Bencode structre
func UnmarshallFromReader(reader io.Reader) (bc bencode.Bencode, err error) {
	data, err := ParseFromReader(reader)

	bc.Data = data

	if err != nil {
		return bc, err
	}

	err = bc.UnmarshallAll()

	return
}
