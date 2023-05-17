package configs

import "fmt"

type DBConf struct {
	Host     string `validate:"required,hostname"`
	Port     int    `validate:"required,number"`
	Database string `validate:"required"`
	Username string `validate:"required"`
	Password string `validate:"required"`
}

func (d DBConf) String() string {
	return fmt.Sprintf("{Host:%s Port:%d, Database:%s}", d.Host, d.Port, d.Database)
}

func (d DBConf) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", d.Username, d.Password, d.Host, d.Port, d.Database)
}
