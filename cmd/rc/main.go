package main

import (
	"os"

	"github.com/shaneoxm/recall/cmd/rc/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
