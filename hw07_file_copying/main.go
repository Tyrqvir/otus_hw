package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if len(from) == 0 || len(to) == 0 {
		fmt.Println("Usage: main.go -from -to -limit -offset")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println("Start copy process...")
	fmt.Printf("From %s to %s.\n", from, to)
	fmt.Printf("Params: limit: %d, offset %d.\n", limit, offset)

	err := Copy(from, to, offset, limit)
	if err != nil {
		fmt.Println(err)
	}
}
