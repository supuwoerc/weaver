package api

import (
	v1 "gin-web/api/v1"
	"github.com/google/wire"
)

var ApiProvider = wire.NewSet(
	v1.NewAttachmentApi,
)
