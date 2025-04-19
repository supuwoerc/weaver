package constant

//go:generate stringer -type=ResourceType -linecomment -output resource_type_string.go
type ResourceType int8

const (
	ViewRoute    ResourceType = iota + 1 // viewRoute
	ViewResource                         // viewResource
	ApiRoute                             // apiRoute
)
