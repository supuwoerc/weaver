package job

import (
	"fmt"
	"gin-web/pkg/constant"
)

type ServerStatus struct{}

func NewServerStatus() *ServerStatus {
	return &ServerStatus{}
}

func (s *ServerStatus) Name() string {
	return constant.ServerStatus.String()
}

func (s *ServerStatus) Handle() {
	fmt.Println(s.Name())
}
