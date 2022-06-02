package utils

import (
	"bufio"
	"fmt"
	"io"
)

var (
	ErrorIsNotAList             = fmt.Errorf("is not a list")
	ErrorIsNotAStringList       = fmt.Errorf("is not a string list")
	ErrorIsNotAListOfList       = fmt.Errorf("is not a list of list")
	ErrorIsNotAListOfStringList = fmt.Errorf("is not a list of string list")
)

// ByteToInteger convert numeric ascii byte to integer
func ByteToInteger(b byte) (int, bool) {
	n := int(b - '0')

	return n, n >= 0 && n <= 9
}

// ReadNBytes read n bytes from a reader
func ReadNBytes(bufioReader *bufio.Reader, len int) (string, error) {
	buffer := make([]byte, len)

	n, err := io.ReadFull(bufioReader, buffer)

	if err != nil {
		return "", err
	}

	if n < len {
		return "", fmt.Errorf("fewer byte(s) read than expected (%d/%d)", n, len)
	}

	return string(buffer), nil
}

// ToListOfStringList read n bytes from a reader
func ToListOfStringList(data interface{}) (list_of_list_of_string [][]string, err error) {
	list_of_interface, ok := data.([]interface{})

	if !ok {
		return nil, ErrorIsNotAList
	}

	for _, sub_list_of_interface := range list_of_interface {
		sub_list_of_list_of_interface, ok := sub_list_of_interface.([]interface{})

		if !ok {
			return nil, ErrorIsNotAListOfList
		}

		sub_list_of_list_of_string := []string{}

		for _, interface_element := range sub_list_of_list_of_interface {
			interface_string, ok := interface_element.(string)

			if !ok {
				return nil, ErrorIsNotAListOfStringList
			}

			sub_list_of_list_of_string = append(sub_list_of_list_of_string, interface_string)
		}

		list_of_list_of_string = append(list_of_list_of_string, sub_list_of_list_of_string)
	}

	return list_of_list_of_string, nil
}

// ToStringList convert an interface to a string list
func ToStringList(data interface{}) (string_list []string, err error) {
	interface_list, ok := data.([]interface{})

	if !ok {
		return nil, ErrorIsNotAList
	}

	for _, interface_element := range interface_list {
		string_element, ok := interface_element.(string)

		if !ok {
			return nil, ErrorIsNotAStringList
		}

		string_list = append(string_list, string_element)
	}

	return
}
