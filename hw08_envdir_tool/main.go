package main

import (
	"fmt"
	"os"
)

const (
	argsCount      int = 2
	defaultErrCode int = 1
)

func main() {
	args := os.Args

	if len(args) < argsCount {
		fmt.Println("Invalid count of arguments, needed minimum ", argsCount)
		os.Exit(defaultErrCode)
	}

	dirPath, cmdWithArgs := args[1], args[2:]

	env, err := ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(defaultErrCode)
	}

	os.Exit(RunCmd(cmdWithArgs, env))
}
