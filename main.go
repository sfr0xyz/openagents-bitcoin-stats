package main

import (
	"github.com/extism/go-pdk"
)

func run() int32 {
	input := pdk.InputString()
	pdk.OutputString(input)
	return 0
}
