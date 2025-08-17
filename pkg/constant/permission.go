package constant

//go:generate stringer -type=PermissionType -linecomment -output permission_type_string.go
type PermissionType int8

const (
	ViewRoute    PermissionType = iota + 1 // viewRoute
	ViewMenu                               // viewMenu
	ViewResource                           // viewResource
	ApiRoute                               // apiRoute
)
