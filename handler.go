package gobencode

import (
	"bufio"
	"fmt"
	"strconv"
)

func handleBytes(bufioReader *bufio.Reader, b byte) (element interface{}, err error) {
	len, _ := byteToInteger(b)

	for {
		b, err = bufioReader.ReadByte()

		if err != nil {
			return nil, fmt.Errorf("failed to read bytes content: %v", err) // ameliorer
		}

		if b == CHAR_DOUBLE_DOT {
			str, err := readNBytes(bufioReader, len)
			if err != nil {
				return nil, fmt.Errorf("failed to read bytes content: %v", err) // ameliorer
			}

			return str, nil
		} else {
			integer, ok := byteToInteger(b)

			if !ok {
				return nil, fmt.Errorf("invalid character for string length: [%c]", b)
			}

			len *= 10
			len += integer
		}
	}
}

func handleInteger(bufioReader *bufio.Reader) (element interface{}, err error) {
	buffer, err := bufioReader.ReadBytes(CHAR_END)

	if err != nil {
		return nil, err
	}

	buffer_str := string(buffer)[:len(buffer)-1]
	integer, err := strconv.Atoi(buffer_str)

	if err != nil {
		if err == strconv.ErrRange {
			return nil, fmt.Errorf("integer value is to big: [%s]", buffer_str)
		}
		return nil, fmt.Errorf("corrupted integer value: [%s]", buffer_str)
	}

	return integer, nil
}

func handleList(bufioReader *bufio.Reader) (element interface{}, err error) {
	list := make([]interface{}, 0)

	for {
		element, err := getElement(bufioReader)

		if err != nil {
			if err == END_ERROR {
				break
			}
			return nil, fmt.Errorf("list element corrupted: %v", err)
		}

		list = append(list, element)
	}

	return list, nil
}

func handleDict(bufioReader *bufio.Reader) (element interface{}, err error) {
	dictionary := make(map[string]interface{})

	for {
		key, err := getElement(bufioReader)

		if err != nil {
			if err == END_ERROR {
				break
			}
			return nil, fmt.Errorf("dictionary key corrupted: %v", err)
		}

		if _, ok := key.(string); !ok {
			return nil, fmt.Errorf("dictionary key corrupted: bad type [%s], (expected string)", "ERADSDFASD")
		}

		element, err := getElement(bufioReader)

		if err != nil {
			if err == END_ERROR {
				return nil, fmt.Errorf("missing dictionary value for the key [%s]", key)
			}
			return nil, fmt.Errorf("dictionary element corrupted: %v", err)
		}

		dictionary[key.(string)] = element
	}

	return dictionary, nil
}
