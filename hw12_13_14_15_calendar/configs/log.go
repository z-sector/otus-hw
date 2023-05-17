package configs

import "fmt"

type LoggerConf struct {
	Level string `validate:"required"`
	JSON  bool
}

func (l LoggerConf) String() string {
	return fmt.Sprintf("{Level:%s Json:%t}", l.Level, l.JSON)
}
