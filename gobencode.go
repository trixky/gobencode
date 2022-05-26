package gobencode

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

const (
	CHAR_INTEGER    = 'i'
	CHAR_LIST       = 'l'
	CHAR_DICTIONARY = 'd'
	CHAR_DOUBLE_DOT = ':'
	CHAR_END        = 'e'
	CHAR_NEGATIVE   = '-'
)

var END_ERROR = errors.New("end character [e]")

func getElement(bufioReader *bufio.Reader) (element interface{}, err error) {
	b, err := bufioReader.ReadByte()

	if err != nil {
		return nil, fmt.Errorf("failed to read byte: %v", err)
	}

	switch {
	case b >= '0' && b <= '9': // bytes
		return handleBytes(bufioReader, b)
	case b == CHAR_INTEGER: // integer
		return handleInteger(bufioReader)
	case b == CHAR_LIST: // list
		return handleList(bufioReader)
	case b == CHAR_DICTIONARY: // dict
		return handleDict(bufioReader)
	case b == CHAR_END: // end
		return nil, END_ERROR
	default:
		return nil, fmt.Errorf("invalid character to start element: [%c]", b)
	}
}

func ParseFromReader(reader io.Reader) (data interface{}, err error) {
	bufioReader, ok := reader.(*bufio.Reader)

	if !ok {
		bufioReader = bufio.NewReader(reader)
	}

	return getElement(bufioReader)
}
