package api

import (
	v1 "gin-web/api/v1"
	"github.com/google/wire"
)

var ApiProvider = wire.NewSet(
	AttachmentApiProvider,
)

// AttachmentApiProvider 附件模块controller
var AttachmentApiProvider = wire.NewSet(
	v1.NewUserApi,
)
