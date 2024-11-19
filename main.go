package main

import (
	"os"

	"github.com/eddycharly/generic-auth-server/pkg/commands/root"
)

func main() {
	root := root.Command()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
