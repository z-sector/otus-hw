package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/z-sector/otus-hw/hw10_program_optimization/internal"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domainWithDot := "." + domain
	reader := bufio.NewReader(r)
	result := make(DomainStat)
	user := new(User)

	for shouldExit := false; !shouldExit; {
		line, err := readLine(reader)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, err
			}

			if len(line) == 0 {
				break
			}
			shouldExit = true
		}

		if err := internal.Unmarshal(line, user); err != nil {
			return nil, err
		}

		if strings.HasSuffix(user.Email, domainWithDot) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}

func readLine(r *bufio.Reader) ([]byte, error) {
	var (
		isPrefix = true
		err      error
		res, ln  []byte
	)

	for isPrefix && err == nil {
		ln, isPrefix, err = r.ReadLine()
		res = append(res, ln...)
	}

	return res, err
}
