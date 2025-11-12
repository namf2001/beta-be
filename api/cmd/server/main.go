package main

import (
	"os"

	"beta-be/cmd/server/commands"

	_ "entgo.io/ent"
)

func main() {
	if err := commands.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
