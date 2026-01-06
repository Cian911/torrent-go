package main

import (
	"fmt"
	"regexp"
	"strings"
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

	t.Run("Decode list", func(t *testing.T) {
		str := "l5:helloi52ei34ee"
		str2 := "i52ei34ei-90ee"
		want := "[\"hello\",52,34]"

		json, err := decodeBencode(str)
		if err != nil {
			t.Error(err)
		}

		if json != want {
			r, _ := regexp.Compile(`-?\d+`)
			for _, v := range strings.Split(str2, "i") {
				fmt.Println(r.FindString(v))
			}
			t.Errorf("Expected %s but got value: %v", want, json)
		}
	})
}
