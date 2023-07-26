package configs

import "fmt"

type MQConf struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
	Host     string `validate:"required,hostname"`
	Port     int    `validate:"required,number"`
	Queue    string `validate:"required"`
}

func (m MQConf) String() string {
	return fmt.Sprintf("{Host:%s Port:%d Queue:%s}", m.Host, m.Port, m.Queue)
}

func (m MQConf) DSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", m.Username, m.Password, m.Host, m.Port)
}
