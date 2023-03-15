package hw09structvalidator

import "strings"

type Condition struct {
	operator string
	operand  string
}

func parseConditions(tag tagT) ([]Condition, error) {
	conditions := strings.Split(string(tag), andSeparator)
	cond := make([]Condition, 0)

	for _, c := range conditions {
		splitedCond := strings.Split(c, ruleSeparator)
		if len(splitedCond) > 2 {
			return nil, ErrValidationFormat
		}

		operator, operand := splitedCond[0], ""
		if len(splitedCond) == 2 {
			operand = splitedCond[1]
		}

		cond = append(cond, Condition{operator: operator, operand: operand})
	}

	return cond, nil
}
