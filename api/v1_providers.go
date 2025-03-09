package api

import (
	v1 "gin-web/api/v1"
	"gin-web/service"
	"github.com/google/wire"
)

// V1Provider api-provider集合
var V1Provider = wire.NewSet(
	AttachmentApiProvider,
)

var AttachmentApiProvider = wire.NewSet(
	v1.NewAttachmentApi,
	v1.NewBasicApi,
	wire.Bind(new(v1.AttachmentService), new(service.AttachmentService)),
	service.NewAttachmentService,
)
