package gobencode

import (
	"bufio"
	"fmt"
)

func byteToInteger(b byte) (int, bool) {
	n := int(b - '0')

	return n, n >= 0 && n <= 9
}

func readNBytes(bufioReader *bufio.Reader, len int) (string, error) {
	buffer := make([]byte, len)

	n, err := bufioReader.Read(buffer)

	if err != nil {
		return "", err
	}

	if n < len {
		return "", fmt.Errorf("fewer byte(s) read than expected (%d/%d)", n, len)
	}

	return string(buffer), nil
}
