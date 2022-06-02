// Package parser provide a Bencode parser
package parser

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"

	"github.com/trixky/gobencode/utils"
)

const (
	char_integer    = 'i'
	char_list       = 'l'
	char_dictionary = 'd'
	char_double_dot = ':'
	char_end        = 'e'
	char_negative   = '-'
)

var (
	ErrorEnd                            = errors.New("end character [e]")
	ErrorInvalidCharacterToStartElement = errors.New("invalid character to start element")
	ErrorFailedToReadByte               = errors.New("failed to read byte")
	ErrorFailedToReadByteContent        = errors.New("failed to read byte content")
	ErrorInvalidStringLengthCharacter   = errors.New("invalid character for string length")
	ErrorIntegerCorrupted               = errors.New("corrupted integer value")
	ErrorListElementCorrupted           = errors.New("list element corrupted")
	ErrorDictionaryKeyCorrupted         = errors.New("dictionary key corrupted")
	ErrorDictionaryElementCorrupted     = errors.New("dictionary element corrupted")
)

// parseBytes parses a byte array in the bencode format from a reader
func parseBytes(bufioReader *bufio.Reader, b byte) (element interface{}, err error) {
	len, _ := utils.ByteToInteger(b)

	for {
		b, err = bufioReader.ReadByte()

		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrorFailedToReadByteContent, err)
		}

		if b == char_double_dot {
			str, err := utils.ReadNBytes(bufioReader, len)
			if err != nil {
				return nil, fmt.Errorf("%w: %v", ErrorFailedToReadByteContent, err)
			}

			return str, nil
		} else {
			integer, ok := utils.ByteToInteger(b)

			if !ok {
				return nil, fmt.Errorf("%w: [%c]", ErrorInvalidStringLengthCharacter, b)
			}

			len *= 10
			len += integer
		}
	}
}

// parseInteger parses an integer in the bencode format from a reader
func parseInteger(bufioReader *bufio.Reader) (element interface{}, err error) {
	buffer, err := bufioReader.ReadBytes(char_end)

	if err != nil {
		return nil, err
	}

	buffer_str := string(buffer)[:len(buffer)-1]
	integer, err := strconv.Atoi(buffer_str)

	if err != nil {
		return nil, fmt.Errorf("%w [%s]: %v", ErrorIntegerCorrupted, buffer_str, err)
	}

	return integer, nil
}

// parseList parses a list in the bencode format from a reader
func parseList(bufioReader *bufio.Reader) (element interface{}, err error) {
	list := make([]interface{}, 0)

	for {
		element, err := ParseElement(bufioReader)

		if err != nil {
			if err == ErrorEnd {
				break
			}
			return nil, fmt.Errorf("%w: %v", ErrorListElementCorrupted, err)
		}

		list = append(list, element)
	}

	return list, nil
}

// parseDictionary parses a dictionary in the bencode format from a reader
func parseDictionary(bufioReader *bufio.Reader) (element interface{}, err error) {
	dictionary := make(map[string]interface{})

	for {
		key, err := ParseElement(bufioReader)

		if err != nil {
			if err == ErrorEnd {
				break
			}
			return nil, fmt.Errorf("%w: %v", ErrorDictionaryKeyCorrupted, err)
		}

		string_key, ok := key.(string)

		if !ok {
			return nil, fmt.Errorf("%w: bad type [%T], (expected string)", ErrorDictionaryKeyCorrupted, key)
		}

		element, err := ParseElement(bufioReader)

		if err != nil {
			if err == ErrorEnd {
				return nil, err
			}
			return nil, fmt.Errorf("%w: %v", ErrorDictionaryElementCorrupted, err)
		}

		dictionary[string_key] = element
	}

	return dictionary, nil
}

// ParseElement parses any type of element in the bencode format from a reader
func ParseElement(bufioReader *bufio.Reader) (element interface{}, err error) {
	b, err := bufioReader.ReadByte()

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrorFailedToReadByte, err)
	}

	switch {
	case b >= '0' && b <= '9': // bytes
		return parseBytes(bufioReader, b)
	case b == char_integer: // integer
		return parseInteger(bufioReader)
	case b == char_list: // list
		return parseList(bufioReader)
	case b == char_dictionary: // dict
		return parseDictionary(bufioReader)
	case b == char_end: // end
		return nil, ErrorEnd
	default:
		return nil, fmt.Errorf("%w: [%c]", ErrorInvalidCharacterToStartElement, b)
	}
}
