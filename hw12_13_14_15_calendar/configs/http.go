package configs

import "fmt"

type HTTPConf struct {
	Host string `validate:"required,ip"`
	Port int    `validate:"required,number"`
}

func (h HTTPConf) String() string {
	return fmt.Sprintf("{Host:%s Port:%d}", h.Host, h.Port)
}
