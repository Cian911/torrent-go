package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var _ = json.Marshal

func decodeBencode(bencodedString string) (any, error) {
	v, next, err := decodeAt(bencodedString, 0)
	if err != nil {
		return nil, err
	}
	if next != len(bencodedString) {
		return nil, fmt.Errorf("trailing data after index %d", next)
	}
	return v, nil
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

