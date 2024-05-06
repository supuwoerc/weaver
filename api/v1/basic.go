package v1

import (
	"gin-web/pkg/global"
	"go.uber.org/zap"
)

type BasicApi struct {
	Logger *zap.SugaredLogger
}

func NewBasicApi() BasicApi {
	return BasicApi{Logger: global.Logger}
}
