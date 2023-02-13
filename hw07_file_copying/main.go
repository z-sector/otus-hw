package main

import (
	"flag"
	"fmt"
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
	if err := validateFlags(); err != nil {
		log.Fatal(err)
	}

	if err := Copy(from, to, offset, limit); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func validateFlags() error {
	isFrom := checkFlag(flagNameFrom)
	isTo := checkFlag(flagNameTo)

	if isFrom && isTo {
		return nil
	}

	msg := "please provide a filename to "
	if isFrom && !isTo {
		return fmt.Errorf(msg + "copy into")
	}
	if isTo && !isFrom {
		return fmt.Errorf(msg + "copy from")
	}
	return fmt.Errorf(msg + "copy from and copy into")
}

func checkFlag(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
