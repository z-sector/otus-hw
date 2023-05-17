package configs

import (
	"fmt"
)

const (
	StorageSQL    = "sql"
	StorageMemory = "memory"
)

type StorageConf struct {
	Type string `validate:"required,oneof=sql memory"`
	DB   *DBConf
}

func (s StorageConf) String() string {
	dbStr := ""
	if s.Type == StorageSQL {
		dbStr = fmt.Sprintf(" DB:%s", s.DB)
	}
	return fmt.Sprintf("{Type:%s%s}", s.Type, dbStr)
}
