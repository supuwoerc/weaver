package job

type SystemJob interface {
	Name() string
	Handle()
}
