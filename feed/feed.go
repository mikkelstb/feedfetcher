package feed

type Feed interface {
	Parse() error
	HasNext() bool
	GetNext() (*NewsItem, error)
	Read() error
	Init(config map[string]string)
}
