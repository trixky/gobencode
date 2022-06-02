package bencode

import (
	"testing"
)

func TestEncodeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: "0:",
		},
		{
			input:    "oui",
			expected: "3:oui",
		},
		{
			input:    "chat",
			expected: "4:chat",
		},
		{
			input:    "supertest",
			expected: "9:supertest",
		},
	}

	for index, test := range tests {
		output := encodeString(test.input)

		if test.expected != output {
			t.Errorf("test %d: expected [%s] | [%s] output", index, test.expected, output)
			continue
		}
	}
}

func TestEncodeInteger(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{
			input:    -1,
			expected: "i-1e",
		},
		{
			input:    0,
			expected: "i0e",
		},
		{
			input:    1,
			expected: "i1e",
		},
		{
			input:    333,
			expected: "i333e",
		},
	}

	for index, test := range tests {
		output := encodeInteger(test.input)

		if test.expected != output {
			t.Errorf("test %d: expected [%s] | [%s] output", index, test.expected, output)
			continue
		}
	}
}

func TestEncodeList(t *testing.T) {
	tests := []struct {
		input    []interface{}
		expected string
	}{
		{
			input:    []interface{}{},
			expected: "le",
		},
		{
			input:    []interface{}{"oui", "non"},
			expected: "l3:oui3:none",
		},
		{
			input:    []interface{}{1, -2, 3},
			expected: "li1ei-2ei3ee",
		},
		{
			input:    []interface{}{1, []interface{}{[]interface{}{"chat"}, "2", -3}, -3},
			expected: "li1ell4:chate1:2i-3eei-3ee",
		},
		{
			input:    []interface{}{1, []interface{}{1, map[string]interface{}{"chat": 8}}, "3"},
			expected: "li1eli1ed4:chati8eee1:3e",
		},
	}

	for index, test := range tests {
		output, err := encodeList(test.input)

		if err != nil {
			t.Errorf("failed to encode list %d: %v", index, err)
			continue
		}

		if test.expected != output {
			t.Errorf("test %d: expected [%s] | [%s] output", index, test.expected, output)
			continue
		}
	}
}

func TestEncodeDictionary(t *testing.T) {
	tests := []struct {
		input    map[string]interface{}
		expected string
	}{
		{
			input:    map[string]interface{}{},
			expected: "de",
		},
		{
			input:    map[string]interface{}{"first": 1, "second": 2, "": "empty"},
			expected: "d0:5:empty5:firsti1e6:secondi2ee",
		},
		{
			input:    map[string]interface{}{"test": 1, "list": []interface{}{1, map[string]interface{}{"chat": 8}}, "last": "3"},
			expected: "d4:last1:34:listli1ed4:chati8eee4:testi1ee",
		},
	}

	for index, test := range tests {
		output, err := encodeDictionary(test.input)

		if err != nil {
			t.Errorf("failed to encode dictionary %d: %v", index, err)
			continue
		}

		if test.expected != output {
			t.Errorf("test %d: expected [%s] | [%s] output", index, test.expected, output)
			continue
		}
	}
}

func TestEncodePieces(t *testing.T) {
	tests := []struct {
		input    []Piece
		expected string
	}{
		{
			input:    []Piece{{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}},
			expected: "20:00000000000000000000",
		},
		{
			input:    []Piece{{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', '1', '2', '3', '4', '5', '6', '7', '8', '9', '!', '@', '3'}},
			expected: "20:abcdefgh123456789!@3",
		},
	}

	for index, test := range tests {
		output := encodePieces(test.input)

		if test.expected != output {
			t.Errorf("test %d: expected [%s] | [%s] output", index, test.expected, output)
			continue
		}
	}
}

func TestEncodeFiles(t *testing.T) {
	tests := []struct {
		input    []File
		expected string
	}{
		{
			input:    []File{},
			expected: "5:filesle",
		}, {
			input: []File{
				{
					Length:         12,
					Path:           "ouiii.txt",
					DecomposedPath: []string{"ouiii.txt"},
					CompletePath:   "ouiii.txt",
				},
			},
			expected: "5:filesld6:lengthi12e4:pathl9:ouiii.txteee",
		}, {
			input: []File{
				{
					Length:         3400,
					Path:           "nooon.txt",
					DecomposedPath: []string{"chat", "nooon.txt"},
					CompletePath:   "chat/nooon.txt",
				},
				{
					Length:         12,
					Path:           "ouiii.txt",
					DecomposedPath: []string{"ouiii.txt"},
					CompletePath:   "ouiii.txt",
				},
			},
			expected: "5:filesld6:lengthi3400e4:pathl4:chat9:nooon.txteed6:lengthi12e4:pathl9:ouiii.txteee",
		},
	}

	for index, test := range tests {
		output, err := encodeFiles(test.input)

		if err != nil {
			t.Errorf("failed to encode files %d: %v", index, err)
			continue
		}

		if test.expected != output {
			t.Errorf("test %d: expected [%s] | [%s] output", index, test.expected, output)
			continue
		}
	}
}

func TestEncodeInfo(t *testing.T) {
	tests := []struct {
		Bc       Bencode
		expected string
	}{
		{
			Bc: Bencode{
				Info: Info{
					Files: []File{
						{
							Length:         12,
							Path:           "ouiii.txt",
							DecomposedPath: []string{"ouiii.txt"},
							CompletePath:   "ouiii.txt",
						},
					},
					PieceLength: 233,
					Pieces: []Piece{
						{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'},
					},
					DirectoryName: "ouiii.txt",
				},
			},
			expected: "d6:lengthi12e4:name9:ouiii.txt12:piece lengthi233e6:pieces20:0123456789abcdefghije",
		},
		{
			Bc: Bencode{
				Info: Info{
					Files: []File{
						{
							Length:         19710976,
							Path:           "resource1.bin",
							DecomposedPath: []string{"resource1.bin"},
							CompletePath:   "Minecraft 1.15.2/resource1.bin",
						},
						{
							Length:         9050674,
							Path:           "resource.bin",
							DecomposedPath: []string{"resource.bin"},
							CompletePath:   "Minecraft 1.15.2/resource.bin",
						},
						{
							Length:         1056768,
							Path:           "setup.exe",
							DecomposedPath: []string{"setup.exe"},
							CompletePath:   "Minecraft 1.15.2/setup.exe",
						},
					},
					DirectoryName: "Minecraft 1.15.2",
					PieceLength:   16384,
					Pieces: []Piece{
						{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'},
					},
				},
			},
			expected: "d5:filesld6:lengthi19710976e4:pathl13:resource1.bineed6:lengthi9050674e4:pathl12:resource.bineed6:lengthi1056768e4:pathl9:setup.exeeee4:name16:Minecraft 1.15.212:piece lengthi16384e6:pieces20:0123456789abcdefghije",
		},
	}

	for _, test := range tests {
		encoded_info, err := encodeInfo(test.Bc.Info)

		if err != nil {
			t.Errorf("failed to encode the bencode: %v", err)
			continue
		}

		if test.expected != encoded_info {
			t.Errorf("expected [%s...] | [%s...] output", test.expected[:60], encoded_info[:60])
			continue
		}
	}
}

func TestEncodeElement(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{
			input:    "oui",
			expected: "3:oui",
		},
		{
			input:    1,
			expected: "i1e",
		},
		{
			input:    []interface{}{},
			expected: "le",
		},
		{
			input:    map[string]interface{}{},
			expected: "de",
		},
		{
			input:    []Piece{},
			expected: "0:",
		},
		{
			input:    Info{},
			expected: "d12:piece lengthi0e6:pieces0:e",
		},
	}

	for index, test := range tests {
		output, err := encodeElement(test.input)

		if err != nil {
			t.Errorf("failed to encode files %d: %v", index, err)
			continue
		}

		if test.expected != output {
			t.Errorf("test %d: expected [%s] | [%s] output", index, test.expected, output)
			continue
		}
	}
}
