package main

import (
	"os"

	"github.com/AlexRojer31/sandbox/internal/sandbox"
)

func main() {
	os.Exit(sandbox.Run(os.Args[1:]))
}
