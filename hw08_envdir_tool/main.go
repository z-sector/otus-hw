package main

import (
	"fmt"
	"log"
	"os"
)

const (
	dirNamePos = 1
	cmdPos     = 2
)

func main() {
	args := os.Args
	if len(args) < cmdPos {
		fmt.Printf("Usage: %s <path_to_env_dir> <command>\n", args[0])
		log.Fatalln("Not enough arguments to run.")
	}

	env, err := ReadDir(args[dirNamePos])
	if err != nil {
		log.Fatalf("Error loading env: %s\n", err)
	}

	returnCode := RunCmd(args[cmdPos:], env)
	os.Exit(returnCode)
}
