package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

// Ensures gofmt doesn't remove the "os" encoding/json import (feel free to remove this!)
var _ = json.Marshal

type EncodingType string

var encodingTypes = [4]EncodingType{
	"string",
	"int",
	"list",
	"directionary",
}

type Encoding struct {
	typ    EncodingType
	value  string
	length int
	sign   bool
}

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
// - i42e -> 42
// - i-42e -> -42
func decodeBencode(bencodedString string) (interface{}, error) {
	if unicode.IsDigit(rune(bencodedString[0])) {
		var firstColonIndex int

		for i := 0; i < len(bencodedString); i++ {
			if bencodedString[i] == ':' {
				firstColonIndex = i
				break
			}
		}

		lengthStr := bencodedString[:firstColonIndex]

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return "", err
		}

		return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil
	} else if bencodedString[0] == 'i' && bencodedString[len(bencodedString)-1] == 'e' {
		startIndex := 1
		sign := false

		if bencodedString[1] == '-' {
			sign = true
		}

		// Check for leading zero value
		if bencodedString[startIndex+1:startIndex+2] == "0" {
			return "", fmt.Errorf("Leading zero value found.")
		}

		encoder := Encoding{
			typ:    "int",
			value:  bencodedString[startIndex : len(bencodedString)-1],
			length: len(bencodedString) - 1,
			sign:   sign,
		}

		// Convert to int
		iVal, err := strconv.Atoi(encoder.value)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse string value: %v", err)
		}

		return iVal, nil
	} else {
		return "", fmt.Errorf("Only strings are supported at the moment")
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]

		decoded, err := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
