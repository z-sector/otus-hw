package main

import (
	"flag"
	"log"
)

const (
	flagNameFrom   = "from"
	flagNameTo     = "to"
	flagNameLimit  = "limit"
	flagNameOffset = "offset"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, flagNameFrom, "", "file to read from")
	flag.StringVar(&to, flagNameTo, "", "file to write to")
	flag.Int64Var(&limit, flagNameLimit, 0, "limit of bytes to copy")
	flag.Int64Var(&offset, flagNameOffset, 0, "offset in input file")
}

func main() {
	flag.Parse()
	if from == "" || to == "" {
		flag.Usage()
		log.Fatal("arguments -from and -to are required")
	}

	if err := Copy(from, to, offset, limit); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
