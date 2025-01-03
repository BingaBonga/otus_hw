package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		log.Panicf("Usage: go-envdir %s %s %s", "/path/to/env/dir", "command", "separate arguments")
	}

	envPath := os.Args[1]
	env, err := ReadDir(envPath)
	if err != nil {
		log.Panicln(err)
	}

	RunCmd(os.Args[2:], env)
}
