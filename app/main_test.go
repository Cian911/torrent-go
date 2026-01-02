package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	t.Run("Decode string", func(t *testing.T) {
		str := "5:hello"

		json, err := decodeBencode(str)
		if err != nil {
			t.Error(err)
		}

		if json != "hello" {
			t.Errorf("Expected hello, got %s", json)
		}
	})

	t.Run("Decode integer", func(t *testing.T) {
		str := "i42e"

		json, err := decodeBencode(str)
		if err != nil {
			t.Error(err)
		}

		if json != 42 {
			t.Errorf("Expected 42, got %s", json)
		}
	})

	t.Run("Decode negative integer", func(t *testing.T) {
		str := "i-42e"

		json, err := decodeBencode(str)
		if err != nil {
			t.Error(err)
		}

		if json != -42 {
			t.Errorf("Expected -42, got %s", json)
		}
	})

	t.Run("Decode leading zero integer", func(t *testing.T) {
		str := "i-0e"

		json, err := decodeBencode(str)
		if err == nil {
			t.Errorf("Expected error but got value: %v", json)
		}
	})
}
