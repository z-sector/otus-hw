package configs

import "fmt"

const (
	ServerHTTP = "http"
	ServerGRPC = "grpc"
)

type ServerConf struct {
	Type string `validate:"required,oneof=http grpc"`
	Host string `validate:"required,ip"`
	Port int    `validate:"required,number"`
}

func (h ServerConf) String() string {
	return fmt.Sprintf("{Type:%s Host:%s Port:%d}", h.Type, h.Host, h.Port)
}
